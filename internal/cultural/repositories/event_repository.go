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
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(event, false, "", "minimal", "").
		Execute()
	return err
}

func (r *eventRepository) FindByID(ctx context.Context, id string) (*models.ResponseEvent, error) {
	var event *models.ResponseEvent
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&event)

	if err == nil && event != nil {
		views, _ := r.GetEventViews(ctx, id)
		event.Views = views
	}

	return event, err
}

// Update methods to handle type conversions and UUID correctly
func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(*event, "minimal", "").
		Eq("id", (*event).ID.String()).
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

func (r *eventRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.ResponseEvent, error) {
	var events []models.ResponseEvent
	query := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false)

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

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, err
}

func (r *eventRepository) GetEventViews(ctx context.Context, eventID string) (int, error) {
	var viewsData struct {
		Views int `json:"views"`
	}

	_, err := r.supabaseClient.
		From("event_views").
		Select("views", "", false).
		Single().
		Eq("event_id", eventID).
		ExecuteTo(&viewsData)

	views := viewsData.Views
	if err != nil {
		return 0, err
	}
	return views, nil
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

func (r *eventRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ResponseEvent, error) {
	var events []models.ResponseEvent
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&events)
	if err != nil {
		return nil, err
	}

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, nil
}

// Specialized methods
// Modify Search method to match base repository interface
func (r *eventRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.ResponseEvent, int, error) {
	var events []models.ResponseEvent
	query := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false)

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

	// Execute query
	_, err := query.ExecuteTo(&events)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, err
	}

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, count, nil
}

// FindPopularEventsWithDetails retrieves popular events with full details
func (r *eventRepository) FindPopularEvents(ctx context.Context, limit int) ([]models.ResponseEvent, error) {
	var events []models.ResponseEvent
	_, err := r.supabaseClient.
		From("event_with_views").
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false).
		Order("views", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&events)

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, err
}

func (r *eventRepository) UpdateViews(ctx context.Context, id string) string {
	return r.supabaseClient.
		Rpc("increment_or_create_event_views", "", map[string]interface{}{
			"event_id": id,
		})
}

// Similar modifications for FindRelatedEvents, FindRecentEvents, FindEventsByLocation
func (r *eventRepository) FindRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.ResponseEvent, error) {
	var events []models.ResponseEvent

	// First, get the original event's details to find related events
	originalEvent, err := r.FindByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find original event: %w", err)
	}

	// Find related events based on similar location, city, or province
	_, err = r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false).
		Or(
			fmt.Sprintf(
				"location_id.eq.%s,city_id.eq.%s,province_id.eq.%s",
				(*originalEvent).Location.ID.String(),
				(*originalEvent).City.ID.String(),
				(*originalEvent).Province.ID.String(),
			),
			"",
		).
		Neq("id", eventID).
		Limit(limit, "").
		ExecuteTo(&events)

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, err
}

// Similar implementation for FindRecentEventsWithDetails and FindEventsByLocationWithDetails
func (r *eventRepository) FindRecentEvents(ctx context.Context, limit int) ([]models.ResponseEvent, error) {
	var events []models.ResponseEvent
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&events)

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, err
}

func (r *eventRepository) FindEventsByLocation(ctx context.Context, locationID uuid.UUID) ([]models.ResponseEvent, error) {
	var events []models.ResponseEvent
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*)))", "", false).
		Eq("location_id", locationID.String()).
		ExecuteTo(&events)

	// Fetch views for each event
	for i := range events {
		views, _ := r.GetEventViews(ctx, events[i].ID.String())
		events[i].Views = views
	}

	return events, err
}
