package services

import (
	"context"
	"fmt"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/sync/errgroup"
)

type UserService struct {
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
	campRepo ports.CampaignRepository
	linkRepo ports.LinkRepository
	tx       ports.Transactor
	logger   log.Logger
}

func NewUserService(
	ur ports.UserRepository,
	pr ports.PlanRepository,
	cr ports.CampaignRepository,
	lr ports.LinkRepository,
	tx ports.Transactor,
	logger log.Logger,
) *UserService {
	return &UserService{
		userRepo: ur,
		planRepo: pr,
		campRepo: cr,
		linkRepo: lr,
		tx:       tx,
		logger:   logger,
	}
}

func (s *UserService) GetAll(ctx context.Context, page, limit int) ([]*domain.User, int64, error) {
	const op = "UserService.GetAll"
	offset := (page - 1) * limit
	var (
		users []*domain.User
		total int64
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		users, err = s.userRepo.FindAll(gCtx, limit, offset)
		return err
	})

	g.Go(func() error {
		var err error
		total, err = s.userRepo.Count(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, ops.WrapMsg(op, ops.KindInternal, err, "failed to list users")
	}

	return users, total, nil
}

func (s *UserService) SetStatus(ctx context.Context, userID domain.UserID, isActive bool) error {
	const op = "UserService.SetStatus"
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ops.WrapMsg(op, ops.KindNotFound, err, "user not found")
	}
	user.IsActive = isActive
	if err := s.userRepo.Save(ctx, user); err != nil {
		return ops.WrapMsg(op, ops.KindInternal, err, "failed to update status")
	}
	return nil
}

func (s *UserService) ChangePlan(ctx context.Context, userID domain.UserID, planID string) error {
	const op = "UserService.ChangePlan"

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		user, err := s.userRepo.FindByID(txCtx, userID)
		if err != nil {
			return err
		}
		if _, err := s.planRepo.FindByID(txCtx, planID); err != nil {
			return fmt.Errorf("plan not found: %w", err)
		}
		user.PlanID = planID
		return s.userRepo.Save(txCtx, user)
	})

	if err != nil {
		return ops.WrapMsg(op, ops.KindInternal, err, "failed to change plan")
	}
	return nil
}

func (s *UserService) GetHierarchy(ctx context.Context, id domain.UserID) (*domain.HierarchyNode, error) {
	const op = "UserService.GetHierarchy"

	g, gCtx := errgroup.WithContext(ctx)

	var (
		user      *domain.User
		campaigns []*domain.Campaign
		links     []*domain.Link
	)

	g.Go(func() error {
		var err error
		user, err = s.userRepo.FindByID(gCtx, id)
		return err
	})

	// Get ALL campaigns
	g.Go(func() error {
		var err error
		campaigns, err = s.campRepo.FindAll(gCtx, domain.CampaignFilter{UserID: id, Limit: 1000})
		return err
	})

	// Get ALL links
	g.Go(func() error {
		var err error
		filter := domain.LinkFilter{UserID: id, Limit: 10000}
		links, err = s.linkRepo.FindAll(gCtx, filter)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, ops.WrapMsg(op, ops.KindInternal, err, "failed to build hierarchy")
	}

	// Group links by CampaignID
	linksByCampaign := make(map[string][]*domain.Link, len(links))
	var orphans []*domain.Link

	for _, l := range links {
		if l.CampaignID != "" {
			linksByCampaign[l.CampaignID] = append(linksByCampaign[l.CampaignID], l)
		} else {
			orphans = append(orphans, l)
		}
	}

	// Build Tree
	// Note: Assuming HierarchyNode children are []*HierarchyNode based on your code
	root := &domain.HierarchyNode{
		Name:     user.Email,
		Type:     "root",
		Value:    len(links), // Total links
		Children: make([]*domain.HierarchyNode, 0),
	}

	for _, camp := range campaigns {
		campNode := &domain.HierarchyNode{
			Name:     camp.Name,
			Type:     "campaign",
			Children: make([]*domain.HierarchyNode, 0),
		}

		if campLinks, ok := linksByCampaign[camp.ID]; ok {
			campNode.Value = len(campLinks)
			for _, l := range campLinks {
				linkNode := &domain.HierarchyNode{
					Name:  l.Slug + " (" + l.TargetURL + ")",
					Type:  "link",
					Value: 1,
				}
				campNode.Children = append(campNode.Children, linkNode)
			}
		}
		root.Children = append(root.Children, campNode)
	}

	// Handle Links without Campaign
	if len(orphans) > 0 {
		orphanNode := &domain.HierarchyNode{
			Name:     "Uncategorized",
			Type:     "collection",
			Value:    len(orphans),
			Children: make([]*domain.HierarchyNode, 0),
		}
		for _, l := range orphans {
			orphanNode.Children = append(orphanNode.Children, &domain.HierarchyNode{
				Name:  l.Slug,
				Type:  "link",
				Value: 1,
			})
		}
		root.Children = append(root.Children, orphanNode)
	}

	return root, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID domain.UserID) error {
	const op = "UserService.DeleteUser"

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Soft Delete User
		if err := s.userRepo.Delete(txCtx, userID); err != nil {
			return err
		}

		// 2. Soft Delete Campaigns
		if err := s.campRepo.DeleteByUserID(txCtx, userID); err != nil {
			return err
		}

		// 3. Soft Delete Links
		if err := s.linkRepo.DeleteByUserID(txCtx, userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return ops.WrapMsg(op, ops.KindInternal, err, "failed to delete user account")
	}

	return nil
}
