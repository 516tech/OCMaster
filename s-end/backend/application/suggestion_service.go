package application

import (
	"context"

	"github.com/ocmaster/backend/domain/suggestion"
	"github.com/ocmaster/backend/infrastructure/pdf"
	"github.com/ocmaster/backend/pkg/apperrors"
)

type SuggestionService struct {
	repo suggestion.Repository
}

func NewSuggestionService(r suggestion.Repository) *SuggestionService {
	return &SuggestionService{repo: r}
}

func (s *SuggestionService) Create(ctx context.Context, sug *suggestion.Suggestion) error {
	return s.repo.Save(ctx, sug)
}

func (s *SuggestionService) GetByID(ctx context.Context, id uint64) (*suggestion.Suggestion, error) {
	sug, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return sug, nil
}

func (s *SuggestionService) ListHistory(ctx context.Context, merchantID uint64, offset, limit int) ([]*suggestion.Suggestion, int64, error) {
	return s.repo.FindByMerchant(ctx, merchantID, offset, limit)
}

func (s *SuggestionService) GeneratePDF(ctx context.Context, id uint64) ([]byte, error) {
	sug, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	data := &pdf.ReportData{
		RiskWarning: sug.RiskWarning,
	}
	if sug.CpuSuggestion != nil {
		data.CpuFreq = sug.CpuSuggestion.Frequency
		data.CpuVoltage = sug.CpuSuggestion.Voltage
		data.CpuNotes = sug.CpuSuggestion.Notes
	}
	if sug.RamSuggestion != nil {
		data.RamFreq = sug.RamSuggestion.Frequency
		data.RamTimings = sug.RamSuggestion.Timings
		data.RamVoltage = sug.RamSuggestion.Voltage
		data.RamNotes = sug.RamSuggestion.Notes
	}
	if sug.StabilityTest != nil {
		data.TestTool = sug.StabilityTest.Tool
		data.TestDuration = sug.StabilityTest.Duration
	}
	return pdf.GenerateReport(data)
}
