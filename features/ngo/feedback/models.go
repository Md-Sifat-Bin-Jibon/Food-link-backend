package feedback

import (
	"time"

	"github.com/google/uuid"
)

type NGOFeedbackEntry struct {
	ID              uuid.UUID `json:"id" db:"id"`
	NGOUserID       uuid.UUID `json:"ngo_user_id" db:"ngo_user_id"`
	RecipientName   string    `json:"recipient_name" db:"recipient_name"`
	PartnerName     string    `json:"partner_name" db:"partner_name"`
	DeliveryDate    time.Time `json:"delivery_date" db:"delivery_date"`
	Rating          *int      `json:"rating,omitempty" db:"rating"`
	Comment         string    `json:"comment" db:"comment"`
	Tags            []string  `json:"tags,omitempty" db:"tags"`
	Photo           string    `json:"photo,omitempty" db:"photo"`
	Status          string    `json:"status" db:"status"`
	CorrectiveAction string   `json:"corrective_action,omitempty" db:"corrective_action"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type NGOImpactStory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	NGOUserID   uuid.UUID `json:"ngo_user_id" db:"ngo_user_id"`
	Title       string    `json:"title" db:"title"`
	Story       string    `json:"story" db:"story"`
	Beneficiaries int     `json:"beneficiaries" db:"beneficiaries"`
	MealsProvided int     `json:"meals_provided" db:"meals_provided"`
	Tags        []string  `json:"tags,omitempty" db:"tags"`
	Photo       string    `json:"photo,omitempty" db:"photo"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateNGOFeedbackRequest struct {
	RecipientName string    `json:"recipient_name" validate:"required,min=1"`
	PartnerName   string    `json:"partner_name" validate:"required,min=1"`
	DeliveryDate  time.Time `json:"delivery_date" validate:"required"`
	Rating        *int      `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment       string    `json:"comment" validate:"required,min=1"`
	Tags          []string  `json:"tags,omitempty"`
	Photo         string    `json:"photo,omitempty"`
}

type CreateNGOImpactStoryRequest struct {
	Title         string   `json:"title" validate:"required,min=1,max=255"`
	Story         string   `json:"story" validate:"required,min=1"`
	Beneficiaries int      `json:"beneficiaries,omitempty"`
	MealsProvided int      `json:"meals_provided,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Photo         string   `json:"photo,omitempty"`
}
