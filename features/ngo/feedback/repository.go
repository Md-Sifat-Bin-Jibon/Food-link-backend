package feedback

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

func (r *Repository) GetAllFeedbackByNGOUserID(ngoUserID uuid.UUID) ([]*NGOFeedbackEntry, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, ngo_user_id, recipient_name, partner_name, delivery_date, rating, comment, tags, photo, status, corrective_action, created_at, updated_at FROM ngo_feedback_entries WHERE ngo_user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, ngoUserID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var feedbacks []*NGOFeedbackEntry
	for rows.Next() {
		feedback := &NGOFeedbackEntry{}
		if err := rows.Scan(&feedback.ID, &feedback.NGOUserID, &feedback.RecipientName, &feedback.PartnerName, &feedback.DeliveryDate, &feedback.Rating, &feedback.Comment, pq.Array(&feedback.Tags), &feedback.Photo, &feedback.Status, &feedback.CorrectiveAction, &feedback.CreatedAt, &feedback.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		feedbacks = append(feedbacks, feedback)
	}
	return feedbacks, nil
}

func (r *Repository) CreateFeedback(feedback *NGOFeedbackEntry) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO ngo_feedback_entries (id, ngo_user_id, recipient_name, partner_name, delivery_date, rating, comment, tags, photo, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, ngo_user_id, recipient_name, partner_name, delivery_date, rating, comment, tags, photo, status, corrective_action, created_at, updated_at`
	return r.db.QueryRow(query, feedback.ID, feedback.NGOUserID, feedback.RecipientName, feedback.PartnerName, feedback.DeliveryDate, feedback.Rating, feedback.Comment, pq.Array(feedback.Tags), feedback.Photo, feedback.Status, time.Now(), time.Now()).Scan(&feedback.ID, &feedback.NGOUserID, &feedback.RecipientName, &feedback.PartnerName, &feedback.DeliveryDate, &feedback.Rating, &feedback.Comment, pq.Array(&feedback.Tags), &feedback.Photo, &feedback.Status, &feedback.CorrectiveAction, &feedback.CreatedAt, &feedback.UpdatedAt)
}

func (r *Repository) GetAllStoriesByNGOUserID(ngoUserID uuid.UUID) ([]*NGOImpactStory, error) {
	if r.db == nil {
		return nil, errors.ErrDatabase
	}
	query := `SELECT id, ngo_user_id, title, story, beneficiaries, meals_provided, tags, photo, created_at, updated_at FROM ngo_impact_stories WHERE ngo_user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, ngoUserID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	defer rows.Close()
	var stories []*NGOImpactStory
	for rows.Next() {
		story := &NGOImpactStory{}
		if err := rows.Scan(&story.ID, &story.NGOUserID, &story.Title, &story.Story, &story.Beneficiaries, &story.MealsProvided, pq.Array(&story.Tags), &story.Photo, &story.CreatedAt, &story.UpdatedAt); err != nil {
			return nil, errors.WrapError(err, errors.ErrDatabase)
		}
		stories = append(stories, story)
	}
	return stories, nil
}

func (r *Repository) CreateStory(story *NGOImpactStory) error {
	if r.db == nil {
		return errors.ErrDatabase
	}
	query := `INSERT INTO ngo_impact_stories (id, ngo_user_id, title, story, beneficiaries, meals_provided, tags, photo, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, ngo_user_id, title, story, beneficiaries, meals_provided, tags, photo, created_at, updated_at`
	return r.db.QueryRow(query, story.ID, story.NGOUserID, story.Title, story.Story, story.Beneficiaries, story.MealsProvided, pq.Array(story.Tags), story.Photo, time.Now(), time.Now()).Scan(&story.ID, &story.NGOUserID, &story.Title, &story.Story, &story.Beneficiaries, &story.MealsProvided, pq.Array(&story.Tags), &story.Photo, &story.CreatedAt, &story.UpdatedAt)
}
