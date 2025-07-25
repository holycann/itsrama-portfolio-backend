package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/supabase-community/supabase-go"
)

// SupabaseBadgeRepository implements BadgeRepository using Supabase
type SupabaseBadgeRepository struct {
	client *supabase.Client
	table  string
}

// NewSupabaseBadgeRepository creates a new Supabase-based badge repository
func NewSupabaseBadgeRepository(client *supabase.Client, table string) BadgeRepository {
	if table == "" {
		table = "badges"
	}
	return &SupabaseBadgeRepository{
		client: client,
		table:  table,
	}
}

// Create adds a new badge to the system
func (r *SupabaseBadgeRepository) Create(ctx context.Context, badge *models.Badge) (*models.Badge, error) {
	var result []models.Badge

	_, err := r.client.
		From(r.table).
		Insert(badge, false, "", "minimal", "").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create badge: %w", err)
	}

	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, fmt.Errorf("no badge created")
}

// FindByID retrieves a badge by its unique identifier
func (r *SupabaseBadgeRepository) FindByID(ctx context.Context, id string) (*models.Badge, error) {
	var badges []models.Badge

	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		ExecuteTo(&badges)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve badge: %w", err)
	}

	if len(badges) > 0 {
		return &badges[0], nil
	}
	return nil, fmt.Errorf("badge not found")
}

// FindAll retrieves all badges with pagination
func (r *SupabaseBadgeRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Badge, error) {
	var badges []models.Badge

	query := r.client.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "")

	_, err := query.ExecuteTo(&badges)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve badges: %w", err)
	}

	return badges, nil
}

// Update modifies an existing badge
func (r *SupabaseBadgeRepository) Update(ctx context.Context, badge *models.Badge) (*models.Badge, error) {
	var result []models.Badge

	_, err := r.client.
		From(r.table).
		Update(badge, "minimal", "").
		Eq("id", badge.ID).
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to update badge: %w", err)
	}

	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, fmt.Errorf("no badge updated")
}

// Delete removes a badge by its ID
func (r *SupabaseBadgeRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.client.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete badge: %w", err)
	}

	return nil
}

// Count returns the total number of badges
func (r *SupabaseBadgeRepository) Count(ctx context.Context) (int, error) {
	_, count, err := r.client.
		From(r.table).
		Select("id", "exact", true).
		Execute()

	if err != nil {
		return 0, fmt.Errorf("failed to count badges: %w", err)
	}

	return int(count), nil
}
