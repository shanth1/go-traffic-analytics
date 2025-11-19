package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type DataSeeder struct {
	LinkRepo  ports.LinkRepository
	ClickRepo ports.ClickRepository
}

func (s *DataSeeder) Seed(campaignCount, linksPerCampaign, clicksPerLink int) {
	ctx := context.Background()

	// Словари для генерации
	referers := []string{"Google", "Facebook", "Twitter", "Direct", "Email Newsletter"}
	osList := []string{"iOS", "Android", "Windows", "MacOS", "Linux"}
	countries := []string{"US", "DE", "GB", "FR", "IN", "BR", "JP"}

	// 1. Создаем кампании (пока просто логически, если нет репо кампаний)
	// ...

	// 2. Создаем ссылки
	for i := 0; i < linksPerCampaign; i++ {
		linkID := uuid.New().String()
		link := &domain.Link{
			ID:        linkID,
			Slug:      fmt.Sprintf("promo-%d", rand.Intn(10000)),
			TargetURL: "https://example.com",
			CreatedAt: time.Now().AddDate(0, -1, 0),
		}
		_ = s.LinkRepo.Save(ctx, link)

		// 3. Генерируем клики (Волна трафика)
		// Симулируем последние 7 дней
		now := time.Now()
		for j := 0; j < clicksPerLink; j++ {
			// Случайное время за последние 7 дней
			// Добавляем "вес" времени, чтобы днем было больше кликов (синусоида)
			daysAgo := rand.Intn(7)
			hour := rand.Intn(24)

			// Простой хак: если час ночь (0-6), пропускаем с вероятностью 80%
			if hour < 6 && rand.Float32() > 0.2 {
				continue
			}

			ts := now.AddDate(0, 0, -daysAgo).Add(time.Duration(hour)*time.Hour + time.Duration(rand.Intn(60))*time.Minute)

			// Корреляция данных (Если iOS, то скорее всего Mobile)
			os := osList[rand.Intn(len(osList))]
			device := "Desktop"
			if os == "iOS" || os == "Android" {
				device = "Mobile"
			}

			click := &domain.ClickEvent{
				ID:        uuid.New().String(),
				LinkID:    linkID,
				Timestamp: ts,
				IP:        "192.168.1.1", // Fake
				Country:   countries[rand.Intn(len(countries))],
				OS:        os,
				Device:    device,
				Referer:   referers[rand.Intn(len(referers))],
			}

			_ = s.ClickRepo.Save(ctx, click)
		}
	}
}
