package menu

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

func (s *Service) GetAllByUserID(userID uuid.UUID) ([]*RestaurantMenuItem, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *Service) GetByID(id uuid.UUID) (*RestaurantMenuItem, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uuid.UUID, req *CreateRestaurantMenuItemRequest) (*RestaurantMenuItem, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	ingredientsJSON := JSONB{"ingredients": req.Ingredients}
	item := &RestaurantMenuItem{
		ID:                 uuid.New(),
		UserID:             userID,
		Name:               req.Name,
		Category:           req.Category,
		Ingredients:        ingredientsJSON,
		PredictedWasteScore: req.PredictedWasteScore,
		Price:              req.Price,
		Margin:             req.Margin,
		Suggestions:        req.Suggestions,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(id uuid.UUID, userID uuid.UUID, req *UpdateRestaurantMenuItemRequest) (*RestaurantMenuItem, error) {
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
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Ingredients != nil {
		item.Ingredients = JSONB{"ingredients": req.Ingredients}
	}
	if req.PredictedWasteScore != "" {
		item.PredictedWasteScore = req.PredictedWasteScore
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Margin != nil {
		item.Margin = *req.Margin
	}
	if req.Suggestions != nil {
		item.Suggestions = req.Suggestions
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
