# Data Transfer Objects (DTOs) Usage Guide

## Overview

Data Transfer Objects (DTOs) are used to control data exposure, validate input, and provide a clean interface between different layers of the application.

## User DTO Usage Examples

### 1. User Creation

```go
// Handler Layer
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Validate and bind input
    var userCreate models.UserCreate
    if err := c.ShouldBindJSON(&userCreate); err != nil {
        // Handle validation error
        return
    }

    // Convert UserCreate to User model
    user := models.User{
        Email:    userCreate.Email,
        Password: hashPassword(userCreate.Password),
        Phone:    userCreate.Phone,
        Role:     userCreate.Role,
    }

    // Service layer creates the user
    if err := userService.CreateUser(ctx, &user); err != nil {
        // Handle service error
        return
    }

    // Convert to DTO for response
    userDTO := user.ToDTO()
    response.SuccessCreated(c, userDTO, "User created successfully")
}
```

### 2. User Authentication

```go
// Service Layer
func (s *userService) Authenticate(ctx context.Context, login models.UserLogin) (*models.UserDTO, error) {
    // Find user by email
    user, err := s.repo.FindByEmail(ctx, login.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Verify password
    if !verifyPassword(user.Password, login.Password) {
        return nil, errors.New("invalid credentials")
    }

    // Convert to DTO
    userDTO := user.ToDTO()
    return &userDTO, nil
}
```

## User Profile DTO Usage

### 1. Profile Creation

```go
// Handler Layer
func (h *UserProfileHandler) CreateProfile(c *gin.Context) {
    // Validate input
    var profileCreate models.UserProfileCreate
    if err := c.ShouldBindJSON(&profileCreate); err != nil {
        response.BadRequest(c, "Invalid profile data", err.Error(), "")
        return
    }

    // Convert to UserProfile model
    profile := models.UserProfile{
        UserID:           profileCreate.UserID,
        Fullname:         profileCreate.Fullname,
        Bio:              profileCreate.Bio,
        AvatarUrl:        profileCreate.AvatarURL,
        IdentityImageUrl: profileCreate.IdentityImageURL,
    }

    // Service creates profile
    if err := profileService.CreateProfile(ctx, &profile); err != nil {
        response.InternalServerError(c, "Profile creation failed", err.Error(), "")
        return
    }

    // Convert to DTO for response
    profileDTO := profile.ToDTO()
    response.SuccessCreated(c, profileDTO, "Profile created successfully")
}
```

### 2. Profile Update

```go
// Service Layer
func (s *userProfileService) UpdateProfile(
    ctx context.Context,
    userID uuid.UUID,
    updateDTO models.UserProfileUpdate,
) (*models.UserProfileDTO, error) {
    // Find existing profile
    profile, err := s.repo.FindByUserID(ctx, userID.String())
    if err != nil {
        return nil, errors.New("profile not found")
    }

    // Update fields if provided
    if updateDTO.Fullname != "" {
        profile.Fullname = updateDTO.Fullname
    }
    if updateDTO.Bio != "" {
        profile.Bio = updateDTO.Bio
    }
    if updateDTO.AvatarURL != "" {
        profile.AvatarUrl = updateDTO.AvatarURL
    }
    if updateDTO.IdentityImageURL != "" {
        profile.IdentityImageUrl = updateDTO.IdentityImageURL
    }

    // Save updated profile
    if err := s.repo.Update(ctx, profile); err != nil {
        return nil, err
    }

    // Convert to DTO
    return profile.ToDTO(), nil
}
```

## User Badge DTO Usage

### 1. Assigning Badges

```go
// Service Layer
func (s *userBadgeService) AssignBadge(
    ctx context.Context,
    badgeCreate models.UserBadgeCreate,
) (*models.UserBadgeDTO, error) {
    // Validate user and badge exist
    if _, err := s.userRepo.FindByID(ctx, badgeCreate.UserID.String()); err != nil {
        return nil, errors.New("user not found")
    }

    // Create user badge
    userBadge := models.UserBadge{
        ID:        uuid.New(),
        UserID:    badgeCreate.UserID,
        BadgeID:   badgeCreate.BadgeID,
        CreatedAt: time.Now(),
    }

    // Save to repository
    if err := s.repo.Create(ctx, &userBadge); err != nil {
        return nil, err
    }

    // Convert to DTO
    return userBadge.ToDTO(), nil
}
```

### 2. Searching User Badges

```go
// Handler Layer
func (h *UserBadgeHandler) SearchUserBadges(c *gin.Context) {
    // Parse search parameters
    var search models.UserBadgeSearch
    if err := c.ShouldBindQuery(&search); err != nil {
        response.BadRequest(c, "Invalid search parameters", err.Error(), "")
        return
    }

    // Prepare list options
    listOptions := repository.ListOptions{
        Limit:  search.Limit,
        Offset: search.Offset,
        Filters: []repository.FilterOption{},
    }

    // Add optional filters
    if search.UserID != uuid.Nil {
        listOptions.Filters = append(listOptions.Filters,
            repository.FilterOption{
                Field:    "user_id",
                Operator: "=",
                Value:    search.UserID,
            },
        )
    }

    // Retrieve badges
    badges, err := s.service.ListUserBadges(ctx, listOptions)
    if err != nil {
        response.InternalServerError(c, "Failed to retrieve badges", err.Error(), "")
        return
    }

    // Convert to DTOs
    badgeDTOs := make([]models.UserBadgeDTO, len(badges))
    for i, badge := range badges {
        badgeDTOs[i] = badge.ToDTO()
    }

    response.SuccessOK(c, badgeDTOs, "Badges retrieved successfully")
}
```

## Best Practices

1. **Validation**: Always validate input using DTO validation tags
2. **Conversion**: Use `ToDTO()` and `FromDTO()` methods for conversions
3. **Exposure Control**: Use DTOs to limit data exposure between layers
4. **Separation of Concerns**: Keep models focused on data storage, DTOs on data transfer

## Common Pitfalls to Avoid

- Don't expose sensitive fields like passwords in DTOs
- Always validate input before converting to models
- Use appropriate conversion methods to maintain data integrity
- Be consistent in using DTOs across different layers

## Performance Considerations

- DTOs add a small overhead due to conversion
- For high-performance scenarios, consider using more lightweight mapping techniques
- Profile and optimize if conversion becomes a bottleneck
