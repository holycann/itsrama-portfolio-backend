package gemini

import (
	"context"
	"fmt"
	"strings"

	achievementModels "github.com/holycann/cultour-backend/internal/achievement/models"
	achievementServices "github.com/holycann/cultour-backend/internal/achievement/services"
	culturalModels "github.com/holycann/cultour-backend/internal/cultural/models"
	culturalServices "github.com/holycann/cultour-backend/internal/cultural/services"
	placeModels "github.com/holycann/cultour-backend/internal/place/models"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// Error types for domain adapters
const (
	ErrInvalidInput errors.ErrorType = "INVALID_INPUT"
)

// DomainAdapterInterface defines the contract for domain-specific adapters
type DomainAdapterInterface interface {
	// Initialize the domain adapter with the knowledge base
	Initialize(kb KnowledgeBaseInterface) error

	// Load domain-specific data into the knowledge base
	LoadData(ctx context.Context) error

	// Build domain-specific context for prompts
	BuildContext(ctx context.Context, params map[string]interface{}) (string, error)
}

// CulturalDomainAdapter provides integration with cultural module
type CulturalDomainAdapter struct {
	eventService culturalServices.EventService
	kb           KnowledgeBaseInterface
	logger       *logger.Logger
}

// NewCulturalDomainAdapter creates a new adapter for the cultural domain
func NewCulturalDomainAdapter(eventService culturalServices.EventService, logger *logger.Logger) *CulturalDomainAdapter {
	return &CulturalDomainAdapter{
		eventService: eventService,
		logger:       logger,
	}
}

// Initialize sets up the domain adapter with the knowledge base
func (a *CulturalDomainAdapter) Initialize(kb KnowledgeBaseInterface) error {
	if kb == nil {
		return errors.New(ErrInvalidInput, "knowledge base cannot be nil", nil)
	}
	a.kb = kb
	return nil
}

// LoadData loads cultural domain data into the knowledge base
func (a *CulturalDomainAdapter) LoadData(ctx context.Context) error {
	if a.kb == nil {
		return errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	// Load events data (can be limited to popular/featured events to avoid overloading)
	events, err := a.eventService.GetTrendingEvents(ctx, 10)
	if err != nil {
		a.logger.Error("Failed to load trending events", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Add each event to the knowledge base
	for _, event := range events {
		eventModel := culturalModels.Event{
			ID:            event.ID,
			Name:          event.Name,
			Description:   event.Description,
			StartDate:     event.StartDate,
			EndDate:       event.EndDate,
			IsKidFriendly: event.IsKidFriendly,
		}
		a.kb.AddEvent(&eventModel)
	}

	return nil
}

// BuildEventContext builds context for event-related AI interactions
func (a *CulturalDomainAdapter) BuildContext(ctx context.Context, params map[string]interface{}) (string, error) {
	if a.kb == nil {
		return "", errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	var contextParts []string

	// Check for event ID in params
	eventIDInterface, hasEventID := params["event_id"]
	if hasEventID {
		eventID, ok := eventIDInterface.(string)
		if ok && eventID != "" {
			// Fetch full event details
			event, err := a.eventService.GetEventByID(ctx, eventID)
			if err == nil && event != nil {
				contextParts = append(contextParts, fmt.Sprintf("Event Context: %s", event.Name))
				contextParts = append(contextParts, fmt.Sprintf("Description: %s", event.Description))
				contextParts = append(contextParts, fmt.Sprintf("Dates: %s to %s",
					event.StartDate.Format("2006-01-02"),
					event.EndDate.Format("2006-01-02")))

				if event.Location != nil {
					contextParts = append(contextParts, fmt.Sprintf("Location: %s", event.Location.Name))
				}

				if event.IsKidFriendly {
					contextParts = append(contextParts, "This event is kid-friendly.")
				}
			}
		}
	}

	return strings.Join(contextParts, "\n"), nil
}

// PlaceDomainAdapter provides integration with place module
type PlaceDomainAdapter struct {
	locationService placeServices.LocationService
	cityService     placeServices.CityService
	provinceService placeServices.ProvinceService
	kb              KnowledgeBaseInterface
	logger          *logger.Logger
}

// NewPlaceDomainAdapter creates a new adapter for the place domain
func NewPlaceDomainAdapter(
	locationService placeServices.LocationService,
	cityService placeServices.CityService,
	provinceService placeServices.ProvinceService,
	logger *logger.Logger,
) *PlaceDomainAdapter {
	return &PlaceDomainAdapter{
		locationService: locationService,
		cityService:     cityService,
		provinceService: provinceService,
		logger:          logger,
	}
}

// Initialize sets up the domain adapter with the knowledge base
func (a *PlaceDomainAdapter) Initialize(kb KnowledgeBaseInterface) error {
	if kb == nil {
		return errors.New(ErrInvalidInput, "knowledge base cannot be nil", nil)
	}
	a.kb = kb
	return nil
}

// LoadData loads place domain data into the knowledge base
func (a *PlaceDomainAdapter) LoadData(ctx context.Context) error {
	if a.kb == nil {
		return errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	// Load provinces (limited number to avoid overloading)
	provinces, err := a.provinceService.ListProvinces(ctx, base.ListOptions{})
	if err != nil {
		a.logger.Error("Failed to load provinces", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Add each province to the knowledge base
	for _, province := range provinces {
		provinceModel := placeModels.Province{
			ID:          province.ID,
			Name:        province.Name,
			Description: province.Description,
		}
		a.kb.AddProvince(&provinceModel)

		// Load cities for each province (limited number)
		cities, err := a.cityService.GetCitiesByProvince(ctx, province.ID.String())
		if err != nil {
			continue // Skip to next province if city loading fails
		}

		// Add each city to the knowledge base
		for _, city := range cities {
			cityModel := placeModels.City{
				ID:          city.ID,
				Name:        city.Name,
				Description: city.Description,
				ProvinceID:  city.ProvinceID,
			}
			a.kb.AddCity(&cityModel)

			// Load popular locations for each city (limited number)
			locations, err := a.locationService.GetLocationsByCity(ctx, city.ID.String())
			if err != nil {
				continue // Skip to next city if location loading fails
			}

			// Add each location to the knowledge base
			for _, location := range locations {
				locationModel := placeModels.Location{
					ID:        location.ID,
					Name:      location.Name,
					CityID:    location.CityID,
					Latitude:  location.Latitude,
					Longitude: location.Longitude,
				}
				a.kb.AddLocation(&locationModel)
			}
		}
	}

	return nil
}

// BuildPlaceContext builds context for place-related AI interactions
func (a *PlaceDomainAdapter) BuildContext(ctx context.Context, params map[string]interface{}) (string, error) {
	if a.kb == nil {
		return "", errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	var contextParts []string

	// Check for location ID in params
	locationIDInterface, hasLocationID := params["location_id"]
	if hasLocationID {
		locationID, ok := locationIDInterface.(string)
		if ok && locationID != "" {
			// Fetch full location details
			location, err := a.locationService.GetLocationByID(ctx, locationID)
			if err == nil && location != nil {
				contextParts = append(contextParts, fmt.Sprintf("Location Context: %s", location.Name))

				if location.City != nil {
					contextParts = append(contextParts, fmt.Sprintf("City: %s", location.City.Name))

					// Fetch province information if city is available
					if location.City.ProvinceID.String() != "" {
						province, err := a.provinceService.GetProvinceByID(ctx, location.City.ProvinceID.String())
						if err == nil && province != nil {
							contextParts = append(contextParts, fmt.Sprintf("Province: %s", province.Name))
						}
					}
				}
			}
		}
	}

	// Check for city ID in params
	cityIDInterface, hasCityID := params["city_id"]
	if hasCityID && len(contextParts) == 0 { // Only process if location wasn't processed
		cityID, ok := cityIDInterface.(string)
		if ok && cityID != "" {
			// Fetch city details
			city, err := a.cityService.GetCityByID(ctx, cityID)
			if err == nil && city != nil {
				contextParts = append(contextParts, fmt.Sprintf("City Context: %s", city.Name))
				contextParts = append(contextParts, fmt.Sprintf("Description: %s", city.Description))

				// Fetch province information
				province, err := a.provinceService.GetProvinceByID(ctx, city.ProvinceID.String())
				if err == nil && province != nil {
					contextParts = append(contextParts, fmt.Sprintf("Province: %s", province.Name))
				}
			}
		}
	}

	return strings.Join(contextParts, "\n"), nil
}

// UserDomainAdapter provides integration with users module
type UserDomainAdapter struct {
	userService        userServices.UserService
	userProfileService userServices.UserProfileService
	userBadgeService   userServices.UserBadgeService
	badgeService       achievementServices.BadgeService
	kb                 KnowledgeBaseInterface
	logger             *logger.Logger
}

// NewUserDomainAdapter creates a new adapter for the user domain
func NewUserDomainAdapter(
	userService userServices.UserService,
	userProfileService userServices.UserProfileService,
	userBadgeService userServices.UserBadgeService,
	badgeService achievementServices.BadgeService,
	logger *logger.Logger,
) *UserDomainAdapter {
	return &UserDomainAdapter{
		userService:        userService,
		userProfileService: userProfileService,
		userBadgeService:   userBadgeService,
		badgeService:       badgeService,
		logger:             logger,
	}
}

// Initialize sets up the domain adapter with the knowledge base
func (a *UserDomainAdapter) Initialize(kb KnowledgeBaseInterface) error {
	if kb == nil {
		return errors.New(ErrInvalidInput, "knowledge base cannot be nil", nil)
	}
	a.kb = kb
	return nil
}

// LoadData loads user domain data into the knowledge base
func (a *UserDomainAdapter) LoadData(ctx context.Context) error {
	if a.kb == nil {
		return errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	// Load badges data
	badges, count, err := a.badgeService.ListBadges(ctx, base.ListOptions{})
	if err != nil {
		a.logger.Error("Failed to load badges", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	a.logger.Info("Loaded badges for knowledge base", map[string]interface{}{
		"count": count,
	})

	// Add each badge to the knowledge base
	for _, badge := range badges {
		badgeModel := achievementModels.Badge{
			ID:          badge.ID,
			Name:        badge.Name,
			Description: badge.Description,
			// Skip ImageURL as it appears to not be in the Badge model
		}
		a.kb.AddBadge(&badgeModel)
	}

	return nil
}

// LoadUserData loads specific user data into the knowledge base
func (a *UserDomainAdapter) LoadUserData(ctx context.Context, userID string) error {
	if a.kb == nil {
		return errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	if userID == "" {
		return errors.New(ErrInvalidInput, "user ID cannot be empty", nil)
	}

	// Load user data
	a.logger.Info("Loading user data", map[string]interface{}{
		"user_id": userID,
	})

	user, err := a.userService.GetUserByID(ctx, userID)
	if err != nil {
		a.logger.Error("Failed to load user data", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return err
	}

	// Add user to knowledge base
	userModel := userModels.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	a.kb.AddUser(&userModel)

	// Load user profile
	profile, err := a.userProfileService.GetProfileByUserID(ctx, userID)
	if err == nil && profile != nil {
		profileModel := userModels.UserProfile{
			ID:               profile.ID,
			UserID:           user.ID, // Use the user ID from the user object
			Fullname:         profile.Fullname,
			Bio:              profile.Bio,
			AvatarUrl:        profile.AvatarUrl,
			IdentityImageUrl: profile.IdentityImageUrl,
		}
		a.kb.AddUserProfile(&profileModel)
	}

	// Load user badges
	badges, err := a.userBadgeService.GetUserBadgesByUser(ctx, userID)
	if err == nil {
		for _, badge := range badges {
			userBadgeModel := userModels.UserBadge{
				UserID:  user.ID,
				BadgeID: badge.BadgeID,
			}
			a.kb.AddUserBadge(userID, &userBadgeModel)
		}
	}

	return nil
}

// BuildUserContext builds context for user-related AI interactions
func (a *UserDomainAdapter) BuildContext(ctx context.Context, params map[string]interface{}) (string, error) {
	if a.kb == nil {
		return "", errors.New(errors.ErrInternal, "knowledge base not initialized", nil)
	}

	var contextParts []string

	// Check for user ID in params
	userIDInterface, hasUserID := params["user_id"]
	if hasUserID {
		userID, ok := userIDInterface.(string)
		if ok && userID != "" {
			// Load user data if not already loaded
			_ = a.LoadUserData(ctx, userID)

			// Fetch user profile
			profile, err := a.userProfileService.GetProfileByUserID(ctx, userID)
			if err == nil && profile != nil {
				contextParts = append(contextParts, fmt.Sprintf("User Context: %s", profile.Fullname))
				if profile.Bio != nil {
					contextParts = append(contextParts, fmt.Sprintf("Bio: %s", profile.Bio))
				}
			}

			// Fetch user badges
			badges, err := a.userBadgeService.GetUserBadgesByUser(ctx, userID)
			if err == nil && len(badges) > 0 {
				badgeNames := make([]string, 0, len(badges))
				for _, badge := range badges {
					if badge.Badge != nil {
						badgeNames = append(badgeNames, badge.Badge.Name)
					}
				}

				if len(badgeNames) > 0 {
					contextParts = append(contextParts, fmt.Sprintf("User has earned badges: %s", strings.Join(badgeNames, ", ")))
				}
			}
		}
	}

	return strings.Join(contextParts, "\n"), nil
}

// DomainAdapterManager orchestrates multiple domain adapters
type DomainAdapterManager struct {
	adapters map[string]DomainAdapterInterface
	kb       KnowledgeBaseInterface
	logger   *logger.Logger
}

// NewDomainAdapterManager creates a new domain adapter manager
func NewDomainAdapterManager(kb KnowledgeBaseInterface, logger *logger.Logger) *DomainAdapterManager {
	return &DomainAdapterManager{
		adapters: make(map[string]DomainAdapterInterface),
		kb:       kb,
		logger:   logger,
	}
}

// RegisterAdapter adds a domain adapter to the manager
func (m *DomainAdapterManager) RegisterAdapter(name string, adapter DomainAdapterInterface) error {
	if adapter == nil {
		return errors.New(ErrInvalidInput, "adapter cannot be nil", nil)
	}

	if err := adapter.Initialize(m.kb); err != nil {
		return err
	}

	m.adapters[name] = adapter
	return nil
}

// LoadAllData loads data from all registered adapters
func (m *DomainAdapterManager) LoadAllData(ctx context.Context) error {
	for name, adapter := range m.adapters {
		if err := adapter.LoadData(ctx); err != nil {
			m.logger.Error("Failed to load data from adapter", map[string]interface{}{
				"adapter": name,
				"error":   err.Error(),
			})
			// Continue with other adapters even if one fails
		}
	}
	return nil
}

// BuildComprehensiveContext builds a complete context using all available adapters
func (m *DomainAdapterManager) BuildComprehensiveContext(ctx context.Context, params map[string]interface{}) (string, error) {
	var allContexts []string

	for name, adapter := range m.adapters {
		adapterContext, err := adapter.BuildContext(ctx, params)
		if err != nil {
			m.logger.Warn("Error building context from adapter", map[string]interface{}{
				"adapter": name,
				"error":   err.Error(),
			})
			// Continue with other adapters even if one fails
		}

		if adapterContext != "" {
			allContexts = append(allContexts, adapterContext)
		}
	}

	return strings.Join(allContexts, "\n\n"), nil
}
