package capacity

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

type NGOCapacitySettings struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	UserID                  uuid.UUID `json:"user_id" db:"user_id"`
	OrgName                 string    `json:"org_name" db:"org_name"`
	Location                string    `json:"location" db:"location"`
	GeoPoint                JSONB     `json:"geo_point,omitempty" db:"geo_point"`
	ManagerName             string    `json:"manager_name" db:"manager_name"`
	ContactPhone            string    `json:"contact_phone" db:"contact_phone"`
	ContactEmail            string    `json:"contact_email,omitempty" db:"contact_email"`
	PreferredFoodTypes      []string  `json:"preferred_food_types,omitempty" db:"preferred_food_types"`
	RestrictedItems         []string  `json:"restricted_items,omitempty" db:"restricted_items"`
	StorageTypes            []string  `json:"storage_types,omitempty" db:"storage_types"`
	SafetyRules             []string  `json:"safety_rules,omitempty" db:"safety_rules"`
	PolicyNotes             string    `json:"policy_notes,omitempty" db:"policy_notes"`
	PickupWindow            JSONB     `json:"pickup_window" db:"pickup_window"`
	DailyCapacityKg         float64   `json:"daily_capacity_kg" db:"daily_capacity_kg"`
	RefrigeratedCapacityKg  float64   `json:"refrigerated_capacity_kg" db:"refrigerated_capacity_kg"`
	DryCapacityKg           float64   `json:"dry_capacity_kg" db:"dry_capacity_kg"`
	CurrentUtilizationKg    float64   `json:"current_utilization_kg" db:"current_utilization_kg"`
	XPPoints                int       `json:"xp_points" db:"xp_points"`
	Level                   int       `json:"level" db:"level"`
	LevelProgressPct        float64   `json:"level_progress_pct" db:"level_progress_pct"`
	AutoAcceptance          JSONB     `json:"auto_acceptance,omitempty" db:"auto_acceptance"`
	PreferredPickupRadiusKm float64  `json:"preferred_pickup_radius_km" db:"preferred_pickup_radius_km"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

type CreateNGOCapacitySettingsRequest struct {
	OrgName                 string                 `json:"org_name" validate:"required,min=1,max=255"`
	Location                string                 `json:"location" validate:"required,min=1"`
	GeoPoint                map[string]interface{} `json:"geo_point,omitempty"`
	ManagerName             string                 `json:"manager_name" validate:"required,min=1,max=255"`
	ContactPhone            string                 `json:"contact_phone" validate:"required,min=1,max=50"`
	ContactEmail            string                 `json:"contact_email,omitempty" validate:"omitempty,email"`
	PreferredFoodTypes      []string               `json:"preferred_food_types,omitempty"`
	RestrictedItems         []string               `json:"restricted_items,omitempty"`
	StorageTypes            []string               `json:"storage_types,omitempty"`
	SafetyRules             []string               `json:"safety_rules,omitempty"`
	PolicyNotes             string                 `json:"policy_notes,omitempty"`
	PickupWindow            map[string]interface{} `json:"pickup_window" validate:"required"`
	DailyCapacityKg         float64                `json:"daily_capacity_kg" validate:"required,gt=0"`
	RefrigeratedCapacityKg  float64                `json:"refrigerated_capacity_kg,omitempty"`
	DryCapacityKg           float64                `json:"dry_capacity_kg,omitempty"`
	AutoAcceptance          map[string]interface{} `json:"auto_acceptance,omitempty"`
	PreferredPickupRadiusKm float64                `json:"preferred_pickup_radius_km,omitempty"`
}
