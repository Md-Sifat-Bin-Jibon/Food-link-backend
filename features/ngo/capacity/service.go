package capacity

import (
	"foodlink_backend/errors"
	"foodlink_backend/utils"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{repo: NewRepository()}
}

func (s *Service) GetByUserID(userID uuid.UUID) (*NGOCapacitySettings, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) CreateOrUpdate(userID uuid.UUID, req *CreateNGOCapacitySettingsRequest) (*NGOCapacitySettings, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	settings := &NGOCapacitySettings{
		ID:                      uuid.New(),
		UserID:                  userID,
		OrgName:                 req.OrgName,
		Location:                req.Location,
		GeoPoint:                JSONB(req.GeoPoint),
		ManagerName:             req.ManagerName,
		ContactPhone:            req.ContactPhone,
		ContactEmail:            req.ContactEmail,
		PreferredFoodTypes:      req.PreferredFoodTypes,
		RestrictedItems:         req.RestrictedItems,
		StorageTypes:            req.StorageTypes,
		SafetyRules:             req.SafetyRules,
		PolicyNotes:             req.PolicyNotes,
		PickupWindow:            JSONB(req.PickupWindow),
		DailyCapacityKg:         req.DailyCapacityKg,
		RefrigeratedCapacityKg:  req.RefrigeratedCapacityKg,
		DryCapacityKg:           req.DryCapacityKg,
		CurrentUtilizationKg:    0,
		XPPoints:                0,
		Level:                   1,
		LevelProgressPct:        0,
		AutoAcceptance:           JSONB(req.AutoAcceptance),
		PreferredPickupRadiusKm: req.PreferredPickupRadiusKm,
	}
	if settings.PreferredPickupRadiusKm == 0 {
		settings.PreferredPickupRadiusKm = 5.0
	}
	if err := s.repo.CreateOrUpdate(settings); err != nil {
		return nil, err
	}
	return settings, nil
}
