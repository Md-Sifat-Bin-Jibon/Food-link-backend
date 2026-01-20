package inventory

import (
	"database/sql"
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

func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]*RestaurantInventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, name, quantity, unit, category, expiry_date, storage_type, batch_code, alert_tags, status, invoice_image, created_at, updated_at FROM restaurant_inventory_items WHERE user_id = $1 ORDER BY expiry_date ASC, created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var items []*RestaurantInventoryItem
	for rows.Next() {
		item := &RestaurantInventoryItem{}
		if err := rows.Scan(&item.ID, &item.UserID, &item.Name, &item.Quantity, &item.Unit, &item.Category, &item.ExpiryDate, &item.StorageType, &item.BatchCode, pq.Array(&item.AlertTags), &item.Status, &item.InvoiceImage, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*RestaurantInventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	item := &RestaurantInventoryItem{}
	query := `SELECT id, user_id, name, quantity, unit, category, expiry_date, storage_type, batch_code, alert_tags, status, invoice_image, created_at, updated_at FROM restaurant_inventory_items WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.UserID, &item.Name, &item.Quantity, &item.Unit, &item.Category, &item.ExpiryDate, &item.StorageType, &item.BatchCode, pq.Array(&item.AlertTags), &item.Status, &item.InvoiceImage, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return item, nil
}

func (r *Repository) Create(item *RestaurantInventoryItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO restaurant_inventory_items (id, user_id, name, quantity, unit, category, expiry_date, storage_type, batch_code, alert_tags, status, invoice_image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, user_id, name, quantity, unit, category, expiry_date, storage_type, batch_code, alert_tags, status, invoice_image, created_at, updated_at`
	now := time.Now()
	return r.db.QueryRow(query, item.ID, item.UserID, item.Name, item.Quantity, item.Unit, item.Category, item.ExpiryDate, item.StorageType, item.BatchCode, pq.Array(item.AlertTags), item.Status, item.InvoiceImage, now, now).Scan(&item.ID, &item.UserID, &item.Name, &item.Quantity, &item.Unit, &item.Category, &item.ExpiryDate, &item.StorageType, &item.BatchCode, pq.Array(&item.AlertTags), &item.Status, &item.InvoiceImage, &item.CreatedAt, &item.UpdatedAt)
}

func (r *Repository) Update(item *RestaurantInventoryItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `UPDATE restaurant_inventory_items SET name=$1, quantity=$2, unit=$3, category=$4, expiry_date=$5, storage_type=$6, batch_code=$7, alert_tags=$8, status=$9, invoice_image=$10, updated_at=$11 WHERE id=$12 RETURNING id, user_id, name, quantity, unit, category, expiry_date, storage_type, batch_code, alert_tags, status, invoice_image, created_at, updated_at`
	return r.db.QueryRow(query, item.Name, item.Quantity, item.Unit, item.Category, item.ExpiryDate, item.StorageType, item.BatchCode, pq.Array(item.AlertTags), item.Status, item.InvoiceImage, time.Now(), item.ID).Scan(&item.ID, &item.UserID, &item.Name, &item.Quantity, &item.Unit, &item.Category, &item.ExpiryDate, &item.StorageType, &item.BatchCode, pq.Array(&item.AlertTags), &item.Status, &item.InvoiceImage, &item.CreatedAt, &item.UpdatedAt)
}

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	result, err := r.db.Exec(`DELETE FROM restaurant_inventory_items WHERE id = $1`, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *Repository) GetExpiring(userID uuid.UUID, days int) ([]*RestaurantInventoryItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	cutoffDate := time.Now().AddDate(0, 0, days)
	query := `SELECT id, user_id, name, quantity, unit, category, expiry_date, storage_type, batch_code, alert_tags, status, invoice_image, created_at, updated_at FROM restaurant_inventory_items WHERE user_id = $1 AND expiry_date <= $2 AND expiry_date >= CURRENT_TIMESTAMP ORDER BY expiry_date ASC`
	rows, err := r.db.Query(query, userID, cutoffDate)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var items []*RestaurantInventoryItem
	for rows.Next() {
		item := &RestaurantInventoryItem{}
		if err := rows.Scan(&item.ID, &item.UserID, &item.Name, &item.Quantity, &item.Unit, &item.Category, &item.ExpiryDate, &item.StorageType, &item.BatchCode, pq.Array(&item.AlertTags), &item.Status, &item.InvoiceImage, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		items = append(items, item)
	}
	return items, nil
}
