package menu

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

func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]*RestaurantMenuItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, user_id, name, category, ingredients, predicted_waste_score, price, margin, suggestions, created_at, updated_at FROM restaurant_menu_items WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var items []*RestaurantMenuItem
	for rows.Next() {
		item := &RestaurantMenuItem{}
		var ingredientsJSON []byte
		if err := rows.Scan(&item.ID, &item.UserID, &item.Name, &item.Category, &ingredientsJSON, &item.PredictedWasteScore, &item.Price, &item.Margin, pq.Array(&item.Suggestions), &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		if len(ingredientsJSON) > 0 {
			json.Unmarshal(ingredientsJSON, &item.Ingredients)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) GetByID(id uuid.UUID) (*RestaurantMenuItem, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	item := &RestaurantMenuItem{}
	var ingredientsJSON []byte
	query := `SELECT id, user_id, name, category, ingredients, predicted_waste_score, price, margin, suggestions, created_at, updated_at FROM restaurant_menu_items WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.UserID, &item.Name, &item.Category, &ingredientsJSON, &item.PredictedWasteScore, &item.Price, &item.Margin, pq.Array(&item.Suggestions), &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	if len(ingredientsJSON) > 0 {
		json.Unmarshal(ingredientsJSON, &item.Ingredients)
	}
	return item, nil
}

func (r *Repository) Create(item *RestaurantMenuItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	ingredientsJSON, _ := json.Marshal(item.Ingredients)
	query := `INSERT INTO restaurant_menu_items (id, user_id, name, category, ingredients, predicted_waste_score, price, margin, suggestions, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, user_id, name, category, ingredients, predicted_waste_score, price, margin, suggestions, created_at, updated_at`
	now := time.Now()
	var ingredientsJSONOut []byte
	err := r.db.QueryRow(query, item.ID, item.UserID, item.Name, item.Category, ingredientsJSON, item.PredictedWasteScore, item.Price, item.Margin, pq.Array(item.Suggestions), now, now).Scan(&item.ID, &item.UserID, &item.Name, &item.Category, &ingredientsJSONOut, &item.PredictedWasteScore, &item.Price, &item.Margin, pq.Array(&item.Suggestions), &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(ingredientsJSONOut) > 0 {
		json.Unmarshal(ingredientsJSONOut, &item.Ingredients)
	}
	return nil
}

func (r *Repository) Update(item *RestaurantMenuItem) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	ingredientsJSON, _ := json.Marshal(item.Ingredients)
	query := `UPDATE restaurant_menu_items SET name=$1, category=$2, ingredients=$3, predicted_waste_score=$4, price=$5, margin=$6, suggestions=$7, updated_at=$8 WHERE id=$9 RETURNING id, user_id, name, category, ingredients, predicted_waste_score, price, margin, suggestions, created_at, updated_at`
	var ingredientsJSONOut []byte
	err := r.db.QueryRow(query, item.Name, item.Category, ingredientsJSON, item.PredictedWasteScore, item.Price, item.Margin, pq.Array(item.Suggestions), time.Now(), item.ID).Scan(&item.ID, &item.UserID, &item.Name, &item.Category, &ingredientsJSONOut, &item.PredictedWasteScore, &item.Price, &item.Margin, pq.Array(&item.Suggestions), &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	if len(ingredientsJSONOut) > 0 {
		json.Unmarshal(ingredientsJSONOut, &item.Ingredients)
	}
	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	result, err := r.db.Exec(`DELETE FROM restaurant_menu_items WHERE id = $1`, id)
	if err != nil {
		return errors.WrapError(err, errors.ErrDatabase)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}
