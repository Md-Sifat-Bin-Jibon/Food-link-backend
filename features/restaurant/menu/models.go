package menu

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type RestaurantMenuItem struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Name               string    `json:"name" db:"name"`
	Category           string    `json:"category" db:"category"`
	Ingredients        JSONB     `json:"ingredients" db:"ingredients"`
	PredictedWasteScore string   `json:"predicted_waste_score,omitempty" db:"predicted_waste_score"`
	Price              float64   `json:"price" db:"price"`
	Margin             float64   `json:"margin" db:"margin"`
	Suggestions        []string  `json:"suggestions,omitempty" db:"suggestions"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRestaurantMenuItemRequest struct {
	Name               string                 `json:"name" validate:"required,min=1,max=255"`
	Category           string                 `json:"category" validate:"required,min=1,max=100"`
	Ingredients        []map[string]interface{} `json:"ingredients" validate:"required"`
	PredictedWasteScore string               `json:"predicted_waste_score,omitempty" validate:"omitempty,oneof=low medium high"`
	Price              float64                `json:"price" validate:"required,gt=0"`
	Margin             float64                `json:"margin" validate:"required"`
	Suggestions        []string               `json:"suggestions,omitempty"`
}

type UpdateRestaurantMenuItemRequest struct {
	Name               string                 `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Category           string                 `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	Ingredients        []map[string]interface{} `json:"ingredients,omitempty"`
	PredictedWasteScore string               `json:"predicted_waste_score,omitempty" validate:"omitempty,oneof=low medium high"`
	Price              *float64               `json:"price,omitempty" validate:"omitempty,gt=0"`
	Margin             *float64               `json:"margin,omitempty"`
	Suggestions        []string               `json:"suggestions,omitempty"`
}
