package inventory

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantInventoryItem struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	Quantity     float64    `json:"quantity" db:"quantity"`
	Unit         string     `json:"unit" db:"unit"`
	Category     string     `json:"category" db:"category"`
	ExpiryDate   time.Time  `json:"expiry_date" db:"expiry_date"`
	StorageType  string     `json:"storage_type" db:"storage_type"`
	BatchCode    string     `json:"batch_code,omitempty" db:"batch_code"`
	AlertTags    []string   `json:"alert_tags,omitempty" db:"alert_tags"`
	Status       string     `json:"status" db:"status"`
	InvoiceImage string     `json:"invoice_image,omitempty" db:"invoice_image"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateRestaurantInventoryItemRequest struct {
	Name         string    `json:"name" validate:"required,min=1,max=255"`
	Quantity     float64   `json:"quantity" validate:"required,gt=0"`
	Unit         string    `json:"unit" validate:"required,min=1,max=50"`
	Category     string    `json:"category" validate:"required,min=1,max=100"`
	ExpiryDate   time.Time `json:"expiry_date" validate:"required"`
	StorageType  string    `json:"storage_type" validate:"required,oneof=fresh chilled frozen dry"`
	BatchCode    string    `json:"batch_code,omitempty" validate:"omitempty,max=100"`
	AlertTags    []string  `json:"alert_tags,omitempty"`
	InvoiceImage string    `json:"invoice_image,omitempty"`
}

type UpdateRestaurantInventoryItemRequest struct {
	Name         string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Quantity     *float64  `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit         string    `json:"unit,omitempty" validate:"omitempty,min=1,max=50"`
	Category     string    `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	StorageType  string    `json:"storage_type,omitempty" validate:"omitempty,oneof=fresh chilled frozen dry"`
	BatchCode    string    `json:"batch_code,omitempty" validate:"omitempty,max=100"`
	AlertTags    []string  `json:"alert_tags,omitempty"`
	Status       string    `json:"status,omitempty"`
	InvoiceImage string    `json:"invoice_image,omitempty"`
}
