package services

import (
	"fmt"
	"math"

	"mo-da-backend/internal/models"
	"mo-da-backend/internal/repositories"
)

type QuarryService struct {
	repo *repositories.QuarryRepo
}

func NewQuarryService() *QuarryService {
	return &QuarryService{repo: repositories.NewQuarryRepo()}
}

func (s *QuarryService) ListQuarries() ([]models.Quarry, error) {
	return s.repo.ListQuarries()
}

func (s *QuarryService) GetQuarryByID(id string) (*models.Quarry, error) {
	return s.repo.GetQuarryByID(id)
}

func (s *QuarryService) CreateQuarry(q *models.Quarry) error {
	if q.Code == "" || q.Name == "" {
		return fmt.Errorf("mã mỏ và tên mỏ không được để trống")
	}
	if q.Status == "" {
		q.Status = "active"
	}
	return s.repo.CreateQuarry(q)
}

func (s *QuarryService) ListQuarryAreas(quarryID string) ([]models.QuarryArea, error) {
	return s.repo.ListQuarryAreas(quarryID)
}

func (s *QuarryService) CreateQuarryArea(a *models.QuarryArea) error {
	if a.QuarryID == "" || a.Code == "" || a.Name == "" {
		return fmt.Errorf("quarryId, code và name không được để trống")
	}
	if a.Status == "" {
		a.Status = "active"
	}
	return s.repo.CreateQuarryArea(a)
}

func (s *QuarryService) ListSurveyCycles(quarryID string) ([]models.SurveyCycle, error) {
	return s.repo.ListSurveyCycles(quarryID)
}

func (s *QuarryService) GetSurveyCycleByID(id string) (*models.SurveyCycle, error) {
	return s.repo.GetSurveyCycleByID(id)
}

func (s *QuarryService) CreateSurveyCycle(cycle *models.SurveyCycle) error {
	if cycle.QuarryID == "" || cycle.CycleCode == "" {
		return fmt.Errorf("quarryId và cycleCode không được để trống")
	}
	if cycle.Status == "" {
		cycle.Status = "completed"
	}
	return s.repo.CreateSurveyCycle(cycle)
}

func (s *QuarryService) ListVolumeCalculations(cycleID string) ([]models.VolumeCalculation, error) {
	return s.repo.ListVolumeCalculations(cycleID)
}

func (s *QuarryService) ListSurfaceModels(cycleID string) ([]models.SurfaceModel, error) {
	return s.repo.ListSurfaceModels(cycleID)
}

func (s *QuarryService) ListMaterials(quarryID string) ([]models.QuarryMaterial, error) {
	return s.repo.ListMaterials(quarryID)
}

func (s *QuarryService) ListWeighingTransactions(quarryID, cycleID string, limit, offset int) ([]models.WeighingTransaction, int, error) {
	return s.repo.ListWeighingTransactions(quarryID, cycleID, limit, offset)
}

func (s *QuarryService) ListReconciliations(quarryID string) ([]models.ProductionReconciliation, error) {
	return s.repo.ListReconciliations(quarryID)
}

func (s *QuarryService) ListAlerts(quarryID string) ([]models.QuarryAlert, error) {
	return s.repo.ListQuarryAlerts(quarryID)
}

func (s *QuarryService) GetDashboardSummary(quarryID string) (*models.QuarryDashboardSummary, error) {
	return s.repo.GetDashboardSummary(quarryID)
}

// CalculateReconciliation executes the core formula comparing 3D Volume vs Actual Weighbridge weight
func (s *QuarryService) CalculateReconciliation(
	quarryID string,
	areaID *string,
	cycleID string,
	materialID *string,
	extractedVolumeM3 float64,
	density float64,
	actualWeightTon float64,
) (*models.ProductionReconciliation, error) {
	if density <= 0 {
		density = 1.55 // default mineral density
	}

	expectedWeightTon := extractedVolumeM3 * density
	differenceTon := expectedWeightTon - actualWeightTon
	var diffPct float64
	if expectedWeightTon > 0 {
		diffPct = (differenceTon / expectedWeightTon) * 100.0
	}

	status := "matched"
	absPct := math.Abs(diffPct)
	if absPct >= 5.0 {
		status = "discrepancy"
	} else if absPct >= 2.0 {
		status = "warning"
	}

	recon := &models.ProductionReconciliation{
		QuarryID:          quarryID,
		QuarryAreaID:      areaID,
		SurveyCycleID:     cycleID,
		MaterialID:        materialID,
		VolumeM3:          extractedVolumeM3,
		DensityTonPerM3:   density,
		ExpectedWeightTon: expectedWeightTon,
		ActualWeightTon:   actualWeightTon,
		DifferenceTon:     differenceTon,
		DifferencePercent: diffPct,
		Status:            status,
	}

	if err := s.repo.CreateReconciliation(recon); err != nil {
		return nil, err
	}

	// Trigger automated alert if difference exceeds threshold (> 3%)
	if absPct >= 3.0 {
		severity := "warning"
		if absPct >= 7.0 {
			severity = "critical"
		}
		alert := &models.QuarryAlert{
			QuarryID:          quarryID,
			QuarryAreaID:      areaID,
			SurveyCycleID:     &cycleID,
			AlertType:         "volume_weight_discrepancy",
			Severity:          severity,
			Title:             fmt.Sprintf("Sai lệch sản lượng đo 3D vs Cân xe: %.2f%%", diffPct),
			Message:           fmt.Sprintf("Thể tích bóc 3D: %.1f m³ (quy đổi %.1f tấn) so với tổng cân thực tế: %.1f tấn. Chênh lệch: %.1f tấn (%.2f%%)", extractedVolumeM3, expectedWeightTon, actualWeightTon, differenceTon, diffPct),
			ExpectedValue:     expectedWeightTon,
			ActualValue:       actualWeightTon,
			DifferencePercent: diffPct,
		}
		_ = s.repo.CreateQuarryAlert(alert)
	}

	return recon, nil
}

func (s *QuarryService) GetQuarryOverview(identifier string) (map[string]interface{}, error) {
	return s.repo.GetQuarryOverview(identifier)
}
