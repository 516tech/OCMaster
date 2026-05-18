package application

import (
	"context"

	"github.com/ocmaster/backend/domain/hardware"
	"github.com/ocmaster/backend/domain/sharecode"
	"github.com/ocmaster/backend/pkg/apperrors"
)

type HardwareService struct {
	hardwareRepo  hardware.Repository
	sharecodeRepo sharecode.Repository
}

func NewHardwareService(hr hardware.Repository, sr sharecode.Repository) *HardwareService {
	return &HardwareService{hardwareRepo: hr, sharecodeRepo: sr}
}

func (s *HardwareService) Upload(ctx context.Context, info *hardware.HardwareInfo) (string, error) {
	code := sharecode.Generate()
	upload := &hardware.HardwareUpload{
		ShareCode:    code,
		HardwareData: *info,
		ExpiresAt:    sharecode.ExpiryDate(),
	}
	if err := s.hardwareRepo.Save(ctx, upload); err != nil {
		return "", err
	}
	return code, nil
}

func (s *HardwareService) GetByCode(ctx context.Context, code string) (*hardware.HardwareUpload, error) {
	upload, err := s.hardwareRepo.FindByShareCode(ctx, code)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	if upload.IsExpired() {
		return nil, apperrors.ErrShareCodeExpired
	}
	return upload, nil
}

func (s *HardwareService) DeleteByCode(ctx context.Context, code string) error {
	return s.hardwareRepo.DeleteByShareCode(ctx, code)
}
