package application

import (
	"context"
	"time"

	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/domain/suggestion"
	"github.com/rs/zerolog/log"
)

type CleanupService struct {
	hardwareRepo   hardware.Repository
	suggestionRepo suggestion.Repository
}

func NewCleanupService(hr hardware.Repository, sr suggestion.Repository) *CleanupService {
	return &CleanupService{hardwareRepo: hr, suggestionRepo: sr}
}

func (s *CleanupService) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				n, err := s.hardwareRepo.DeleteExpired(ctx)
				if err != nil {
					log.Error().Err(err).Msg("cleanup hardware failed")
				} else if n > 0 {
					log.Info().Int64("count", n).Msg("cleaned expired hardware")
				}

				m, err := s.suggestionRepo.DeleteOldHistory(ctx, 30)
				if err != nil {
					log.Error().Err(err).Msg("cleanup history failed")
				} else if m > 0 {
					log.Info().Int64("count", m).Msg("cleaned old query history")
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
