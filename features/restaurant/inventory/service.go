package inventory

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

func (s *Service) GetAllByUserID(userID uuid.UUID) ([]*RestaurantInventoryItem, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *Service) GetByID(id uuid.UUID) (*RestaurantInventoryItem, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, req *CreateRestaurantInventoryItemRequest) (*RestaurantInventoryItem, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	item := &RestaurantInventoryItem{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Category:    req.Category,
		ExpiryDate:  req.ExpiryDate,
		StorageType: req.StorageType,
		BatchCode:   req.BatchCode,
		AlertTags:   req.AlertTags,
		Status:      "normal",
		InvoiceImage: req.InvoiceImage,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateRestaurantInventoryItemRequest) (*RestaurantInventoryItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item.UserID != userID {
		return nil, errors.ErrForbidden
	}
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.ExpiryDate != nil {
		item.ExpiryDate = *req.ExpiryDate
	}
	if req.StorageType != "" {
		item.StorageType = req.StorageType
	}
	if req.BatchCode != "" {
		item.BatchCode = req.BatchCode
	}
	if req.AlertTags != nil {
		item.AlertTags = req.AlertTags
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.InvoiceImage != "" {
		item.InvoiceImage = req.InvoiceImage
	}
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Delete(id uuid.UUID, userID uuid.UUID) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return errors.ErrForbidden
	}
	return s.repo.Delete(id)
}

func (s *Service) GetExpiring(userID uuid.UUID, days int) ([]*RestaurantInventoryItem, error) {
	if days <= 0 {
		days = 7
	}
	return s.repo.GetExpiring(userID, days)
}
