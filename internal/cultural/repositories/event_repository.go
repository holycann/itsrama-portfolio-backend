package repositories

import (
	"context"
	"fmt"
	"sort"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
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

func (r *eventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(event, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create event")
	}
	return event, nil
}

func (r *eventRepository) FindByID(ctx context.Context, id string) (*models.EventDTO, error) {
	var eventDTO models.EventDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&eventDTO)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch event")
	}
	return &eventDTO, nil
}

func (r *eventRepository) Update(ctx context.Context, event *models.Event) (*models.Event, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(event, "minimal", "").
		Eq("id", event.ID.String()).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update event")
	}
	return event, nil
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete event")
	}
	return nil
}

func (r *eventRepository) List(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, error) {
	var events []models.EventDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		case base.OperatorGreaterThan:
			query = query.Gt(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLessThan:
			query = query.Lt(filter.Field, fmt.Sprintf("%v", filter.Value))
		}
	}

	// Apply search if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("description.ilike.%%%s%%", opts.Search),
		)
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == base.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err := query.ExecuteTo(&events)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list events")
	}

	return events, nil
}

func (r *eventRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	query := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true)

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	_, count, err := query.Execute()
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count events")
	}

	return int(count), nil
}

func (r *eventRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check event existence")
	}

	return count > 0, nil
}

func (r *eventRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.EventDTO, error) {
	var events []models.EventDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&events)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find events by field")
	}
	return events, nil
}

func (r *eventRepository) FindEventsByLocation(ctx context.Context, locationID string) ([]models.EventDTO, error) {
	var events []models.EventDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false).
		Eq("location_id", locationID).
		ExecuteTo(&events)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find events by location")
	}
	return events, nil
}

func (r *eventRepository) FindRelatedEvents(ctx context.Context, eventID, locationID string, limit int) ([]models.EventDTO, error) {
	var events []models.EventDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false).
		Eq("location_id", locationID).
		Neq("id", eventID).
		Limit(limit, "").
		ExecuteTo(&events)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find related events")
	}
	return events, nil
}

func (r *eventRepository) GetEventViews(ctx context.Context, eventID string) (int, error) {
	var views struct {
		Views int `json:"views"`
	}
	_, err := r.supabaseClient.
		From("event_with_views").
		Select("views", "", false).
		Eq("id", eventID).
		Single().
		ExecuteTo(&views)

	if err != nil {
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to get event views")
	}

	return views.Views, nil
}

func (r *eventRepository) UpdateViews(ctx context.Context, userID, eventID string) string {
	return r.supabaseClient.
		Rpc("insert_event_view_once", "", map[string]interface{}{
			"_event_id": eventID,
			"_user_id":  userID,
		})
}

func (r *eventRepository) FindPopularEvents(ctx context.Context, limit int) ([]models.EventDTO, error) {
	var events []models.EventDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false).
		Limit(limit, "").
		ExecuteTo(&events)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find popular events")
	}

	// Sort events manually by views in descending order
	sort.Slice(events, func(i, j int) bool {
		viewsI := 0
		viewsJ := 0

		if events[i].Views != nil {
			viewsI = events[i].Views["views"]
		}

		if events[j].Views != nil {
			viewsJ = events[j].Views["views"]
		}

		return viewsI > viewsJ
	})

	return events, nil
}

func (r *eventRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, int, error) {
	var events []models.EventDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, location:locations(*, city:cities(*, province:provinces(*))), creator:users_view!events_user_id_fkey(*), views:event_with_views(views)", "", false)

	// Apply search if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("description.ilike.%%%s%%", opts.Search),
		)
	}

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Count total results
	_, totalCount, err := query.Execute()
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to count search results")
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err = query.ExecuteTo(&events)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to execute search")
	}

	return events, int(totalCount), nil
}

func (r *eventRepository) BulkCreate(ctx context.Context, events []*models.Event) ([]models.Event, error) {
	var results []models.Event
	for _, event := range events {
		createdEvent, err := r.Create(ctx, event)
		if err != nil {
			return nil, err
		}
		results = append(results, *createdEvent)
	}
	return results, nil
}

func (r *eventRepository) BulkUpdate(ctx context.Context, events []*models.Event) ([]models.Event, error) {
	var results []models.Event
	for _, event := range events {
		updatedEvent, err := r.Update(ctx, event)
		if err != nil {
			return nil, err
		}
		results = append(results, *updatedEvent)
	}
	return results, nil
}

func (r *eventRepository) BulkDelete(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.Delete(context.Background(), id); err != nil {
			return err
		}
	}
	return nil
}

func (r *eventRepository) BulkUpsert(ctx context.Context, events []*models.Event) ([]models.EventDTO, error) {
	var results []models.EventDTO
	for _, event := range events {
		// Check if event exists
		exists, err := r.Exists(ctx, event.ID.String())
		if err != nil {
			return nil, err
		}

		var upsertedEvent *models.EventDTO
		if exists {
			// Update existing event
			updatedEvent, err := r.Update(ctx, event)
			if err != nil {
				return nil, err
			}
			upsertedEvent, err = r.FindByID(ctx, updatedEvent.ID.String())
			if err != nil {
				return nil, err
			}
		} else {
			// Create new event
			createdEvent, err := r.Create(ctx, event)
			if err != nil {
				return nil, err
			}
			upsertedEvent, err = r.FindByID(ctx, createdEvent.ID.String())
			if err != nil {
				return nil, err
			}
		}

		results = append(results, *upsertedEvent)
	}
	return results, nil
}
