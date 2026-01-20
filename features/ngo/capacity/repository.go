package capacity

import (
	"database/sql"
	"encoding/json"
	"foodlink_backend/database"
	"foodlink_backend/errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{db: database.GetDB()}
}

func (r *Repository) GetByUserID(userID uuid.UUID) (*NGOCapacitySettings, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	settings := &NGOCapacitySettings{}
	var geoPointJSON, pickupWindowJSON, autoAcceptanceJSON []byte
	query := `SELECT id, user_id, org_name, location, geo_point, manager_name, contact_phone, contact_email, preferred_food_types, restricted_items, storage_types, safety_rules, policy_notes, pickup_window, daily_capacity_kg, refrigerated_capacity_kg, dry_capacity_kg, current_utilization_kg, xp_points, level, level_progress_pct, auto_acceptance, preferred_pickup_radius_km, updated_at FROM ngo_capacity_settings WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&settings.ID, &settings.UserID, &settings.OrgName, &settings.Location, &geoPointJSON, &settings.ManagerName, &settings.ContactPhone, &settings.ContactEmail, pq.Array(&settings.PreferredFoodTypes), pq.Array(&settings.RestrictedItems), pq.Array(&settings.StorageTypes), pq.Array(&settings.SafetyRules), &settings.PolicyNotes, &pickupWindowJSON, &settings.DailyCapacityKg, &settings.RefrigeratedCapacityKg, &settings.DryCapacityKg, &settings.CurrentUtilizationKg, &settings.XPPoints, &settings.Level, &settings.LevelProgressPct, &autoAcceptanceJSON, &settings.PreferredPickupRadiusKm, &settings.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(geoPointJSON) > 0 {
		json.Unmarshal(geoPointJSON, &settings.GeoPoint)
	}
	if len(pickupWindowJSON) > 0 {
		json.Unmarshal(pickupWindowJSON, &settings.PickupWindow)
	}
	if len(autoAcceptanceJSON) > 0 {
		json.Unmarshal(autoAcceptanceJSON, &settings.AutoAcceptance)
	}
	return settings, nil
}

func (r *Repository) CreateOrUpdate(settings *NGOCapacitySettings) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	geoPointJSON, _ := json.Marshal(settings.GeoPoint)
	pickupWindowJSON, _ := json.Marshal(settings.PickupWindow)
	autoAcceptanceJSON, _ := json.Marshal(settings.AutoAcceptance)
	query := `INSERT INTO ngo_capacity_settings (id, user_id, org_name, location, geo_point, manager_name, contact_phone, contact_email, preferred_food_types, restricted_items, storage_types, safety_rules, policy_notes, pickup_window, daily_capacity_kg, refrigerated_capacity_kg, dry_capacity_kg, current_utilization_kg, xp_points, level, level_progress_pct, auto_acceptance, preferred_pickup_radius_km, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24) ON CONFLICT (user_id) DO UPDATE SET org_name=EXCLUDED.org_name, location=EXCLUDED.location, geo_point=EXCLUDED.geo_point, manager_name=EXCLUDED.manager_name, contact_phone=EXCLUDED.contact_phone, contact_email=EXCLUDED.contact_email, preferred_food_types=EXCLUDED.preferred_food_types, restricted_items=EXCLUDED.restricted_items, storage_types=EXCLUDED.storage_types, safety_rules=EXCLUDED.safety_rules, policy_notes=EXCLUDED.policy_notes, pickup_window=EXCLUDED.pickup_window, daily_capacity_kg=EXCLUDED.daily_capacity_kg, refrigerated_capacity_kg=EXCLUDED.refrigerated_capacity_kg, dry_capacity_kg=EXCLUDED.dry_capacity_kg, auto_acceptance=EXCLUDED.auto_acceptance, preferred_pickup_radius_km=EXCLUDED.preferred_pickup_radius_km, updated_at=EXCLUDED.updated_at RETURNING id, user_id, org_name, location, geo_point, manager_name, contact_phone, contact_email, preferred_food_types, restricted_items, storage_types, safety_rules, policy_notes, pickup_window, daily_capacity_kg, refrigerated_capacity_kg, dry_capacity_kg, current_utilization_kg, xp_points, level, level_progress_pct, auto_acceptance, preferred_pickup_radius_km, updated_at`
	var geoPointJSONOut, pickupWindowJSONOut, autoAcceptanceJSONOut []byte
	err := r.db.QueryRow(query, settings.ID, settings.UserID, settings.OrgName, settings.Location, geoPointJSON, settings.ManagerName, settings.ContactPhone, settings.ContactEmail, pq.Array(settings.PreferredFoodTypes), pq.Array(settings.RestrictedItems), pq.Array(settings.StorageTypes), pq.Array(settings.SafetyRules), settings.PolicyNotes, pickupWindowJSON, settings.DailyCapacityKg, settings.RefrigeratedCapacityKg, settings.DryCapacityKg, settings.CurrentUtilizationKg, settings.XPPoints, settings.Level, settings.LevelProgressPct, autoAcceptanceJSON, settings.PreferredPickupRadiusKm, time.Now()).Scan(&settings.ID, &settings.UserID, &settings.OrgName, &settings.Location, &geoPointJSONOut, &settings.ManagerName, &settings.ContactPhone, &settings.ContactEmail, pq.Array(&settings.PreferredFoodTypes), pq.Array(&settings.RestrictedItems), pq.Array(&settings.StorageTypes), pq.Array(&settings.SafetyRules), &settings.PolicyNotes, &pickupWindowJSONOut, &settings.DailyCapacityKg, &settings.RefrigeratedCapacityKg, &settings.DryCapacityKg, &settings.CurrentUtilizationKg, &settings.XPPoints, &settings.Level, &settings.LevelProgressPct, &autoAcceptanceJSONOut, &settings.PreferredPickupRadiusKm, &settings.UpdatedAt)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(geoPointJSONOut) > 0 {
		json.Unmarshal(geoPointJSONOut, &settings.GeoPoint)
	}
	if len(pickupWindowJSONOut) > 0 {
		json.Unmarshal(pickupWindowJSONOut, &settings.PickupWindow)
	}
	if len(autoAcceptanceJSONOut) > 0 {
		json.Unmarshal(autoAcceptanceJSONOut, &settings.AutoAcceptance)
	}
	return nil
}
