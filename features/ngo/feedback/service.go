package feedback

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

func (s *Service) GetAllFeedback(ngoUserID uuid.UUID) ([]*NGOFeedbackEntry, error) {
	return s.repo.GetAllFeedbackByNGOUserID(ngoUserID)
}

func (s *Service) CreateFeedback(ngoUserID uuid.UUID, req *CreateNGOFeedbackRequest) (*NGOFeedbackEntry, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	feedback := &NGOFeedbackEntry{
		ID:            uuid.New(),
		NGOUserID:     ngoUserID,
		RecipientName: req.RecipientName,
		PartnerName:   req.PartnerName,
		DeliveryDate:  req.DeliveryDate,
		Rating:        req.Rating,
		Comment:       req.Comment,
		Tags:          req.Tags,
		Photo:         req.Photo,
		Status:        "pending",
	}
	if err := s.repo.CreateFeedback(feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (s *Service) GetAllStories(ngoUserID uuid.UUID) ([]*NGOImpactStory, error) {
	return s.repo.GetAllStoriesByNGOUserID(ngoUserID)
}

func (s *Service) CreateStory(ngoUserID uuid.UUID, req *CreateNGOImpactStoryRequest) (*NGOImpactStory, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.NewAppErrorWithErr(errors.ErrValidationFailed.Code, "Validation failed: "+validationErrors[0], nil)
	}
	story := &NGOImpactStory{
		ID:            uuid.New(),
		NGOUserID:     ngoUserID,
		Title:         req.Title,
		Story:         req.Story,
		Beneficiaries: req.Beneficiaries,
		MealsProvided: req.MealsProvided,
		Tags:          req.Tags,
		Photo:         req.Photo,
	}
	if err := s.repo.CreateStory(story); err != nil {
		return nil, err
	}
	return story, nil
}
