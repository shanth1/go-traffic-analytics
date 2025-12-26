package services

import (
	"context"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/sync/errgroup"
)

type UserService struct {
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
	campRepo ports.CampaignRepository
	linkRepo ports.LinkRepository
}

func NewUserService(ur ports.UserRepository, pr ports.PlanRepository, cr ports.CampaignRepository, lr ports.LinkRepository) *UserService {
	return &UserService{
		userRepo: ur,
		planRepo: pr,
		campRepo: cr,
		linkRepo: lr,
	}
}

func (s *UserService) GetAll(ctx context.Context, page, limit int) ([]*domain.User, error) {
	offset := (page - 1) * limit
	return s.userRepo.FindAll(ctx, limit, offset)
}

func (s *UserService) SetStatus(ctx context.Context, userID domain.UserID, isActive bool) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.IsActive = isActive
	return s.userRepo.Save(ctx, user)
}

func (s *UserService) ChangePlan(ctx context.Context, userID domain.UserID, planID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if _, err := s.planRepo.FindByID(ctx, planID); err != nil {
		return err
	}
	user.PlanID = planID
	return s.userRepo.Save(ctx, user)
}

func (s *UserService) GetHierarchy(ctx context.Context, id domain.UserID) (*domain.HierarchyNode, error) {
	g, ctx := errgroup.WithContext(ctx)

	var (
		user      *domain.User
		campaigns []*domain.Campaign
		links     []*domain.Link
	)

	g.Go(func() error {
		var err error
		user, err = s.userRepo.FindByID(ctx, id)
		return err
	})

	g.Go(func() error {
		var err error
		campaigns, err = s.campRepo.FindAllByUserID(ctx, id)
		return err
	})

	g.Go(func() error {
		var err error
		filter := domain.LinkFilter{
			UserID: id,
			Limit:  10000,
		}
		links, err = s.linkRepo.FindAll(ctx, filter)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	linksByCampaign := make(map[string][]*domain.Link, len(links))
	for _, l := range links {
		if l.CampaignID != "" {
			linksByCampaign[l.CampaignID] = append(linksByCampaign[l.CampaignID], l)
		}
	}

	root := &domain.HierarchyNode{
		Name:     user.Email,
		Type:     "root",
		Children: make([]*domain.HierarchyNode, 0, len(campaigns)),
	}

	for _, camp := range campaigns {
		campNode := &domain.HierarchyNode{
			Name:     camp.Name,
			Type:     "campaign",
			Children: make([]*domain.HierarchyNode, 0),
		}

		if campLinks, ok := linksByCampaign[camp.ID]; ok {
			for _, l := range campLinks {
				linkNode := &domain.HierarchyNode{
					Name:  l.TargetURL,
					Type:  "link",
					Value: 1,
				}
				campNode.Children = append(campNode.Children, linkNode)
			}
		}

		root.Children = append(root.Children, campNode)
	}

	return root, nil
}
