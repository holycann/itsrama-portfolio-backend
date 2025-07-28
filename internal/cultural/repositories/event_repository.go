// Package repositories provides an implementation of repository for event data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type eventRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewEventRepository(supabaseClient *supabase.Client) EventRepository {
	return &eventRepository{
		supabaseClient: supabaseClient,
		table:          "events",
	}
}

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(event, false, "", "minimal", "").
		ExecuteTo(&event)
	return err
}

func (r *eventRepository) FindByID(ctx context.Context, id string) (*models.Event, error) {
	var event *models.Event
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&event)
	return event, err
}

// Update methods to handle type conversions and UUID correctly
func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(event, "minimal", "").
		Eq("id", event.ID.String()).
		Execute()
	return err
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *eventRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.Event, error) {
	var events []models.Event
	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == repository.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	query = query.Range(opts.Offset, opts.Offset+opts.Limit-1, "")

	_, err := query.ExecuteTo(&events)
	return events, err
}

func (r *eventRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
	query := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false)

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	_, count, err := query.Execute()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *eventRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *eventRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Event, error) {
	var events []models.Event
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&events)
	return events, err
}

// Specialized methods
func (r *eventRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.Event, int, error) {
	var events []models.Event
	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.SearchQuery != "" {
		escapedQuery := strings.ReplaceAll(strings.ReplaceAll(opts.SearchQuery, "%", "\\%"), "_", "\\_")
		likeQuery := "%" + escapedQuery + "%"
		query = query.Or(fmt.Sprintf("name.ilike.%s,description.ilike.%s", likeQuery, likeQuery), "")
	}

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == repository.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	query = query.Range(opts.Offset, opts.Offset+opts.Limit-1, "")

	_, err := query.ExecuteTo(&events)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

// Add a new method to get event views
func (r *eventRepository) GetEventViews(ctx context.Context, eventID string) (int, error) {
	var views int
	_, err := r.supabaseClient.
		From("event_views").
		Select("views", "", false).
		Eq("event_id", eventID).
		Single().
		ExecuteTo(&views)

	// If no views found, return 0 instead of an error
	if err != nil {
		return 0, nil
	}

	return views, nil
}

// Modify FindPopularEvents to use event_with_views
func (r *eventRepository) FindPopularEvents(ctx context.Context, limit int) ([]models.Event, error) {
	var events []models.Event
	_, err := r.supabaseClient.
		From("event_with_views").
		Select("*", "", false).
		Order("views", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&events)
	return events, err
}

// Update the ResponseEvent model to include views
func (r *eventRepository) GetEventWithViews(ctx context.Context, id string) (*models.ResponseEvent, error) {
	var responseEvent models.ResponseEvent

	// First, get the event details
	event, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get the views
	views, _ := r.GetEventViews(ctx, id)

	// Construct the response event
	responseEvent.Event = *event
	responseEvent.Views = views // Assuming you'll add a Views field to the ResponseEvent struct

	return &responseEvent, nil
}

func (r *eventRepository) FindRecentEvents(ctx context.Context, limit int) ([]models.Event, error) {
	var events []models.Event
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&events)
	return events, err
}

func (r *eventRepository) FindEventsByLocation(ctx context.Context, locationID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("location_id", locationID.String()).
		ExecuteTo(&events)
	return events, err
}

func (r *eventRepository) UpdateViews(ctx context.Context, id string) string {
	return r.supabaseClient.
		Rpc("increment_or_create_event_views", "", map[string]interface{}{
			"event_id": id,
		})
}

func (r *eventRepository) FindRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.Event, error) {
	var events []models.Event

	// First, get the original event's details to find related events
	originalEvent, err := r.FindByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find original event: %w", err)
	}

	// Find related events based on similar location, city, or province
	_, err = r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Or(
			fmt.Sprintf(
				"location_id.eq.%s,city_id.eq.%s,province_id.eq.%s",
				originalEvent.LocationID.String(),
				originalEvent.CityID.String(),
				originalEvent.ProvinceID.String(),
			),
			"",
		).
		Neq("id", eventID).
		Limit(limit, "").
		ExecuteTo(&events)

	return events, err
}
