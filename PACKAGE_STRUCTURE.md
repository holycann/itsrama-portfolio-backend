# Package Structure Standardization Guide

## Overview
This document outlines the standard structure for domain packages in the Cultour backend.

## Standard Package Layout

### Directory Structure
```
domain/
├── handlers/
│   └── domain_handler.go
├── models/
│   ├── domain.go
│   └── domain_request_response.go
├── repositories/
│   ├── domain_repository.go
│   └── interface.go
└── services/
    ├── domain_service.go
    └── interface.go
```

## Naming Conventions

### Handlers
- Suffix with `Handler`
- Implement HTTP request/response logic
- Minimal business logic
- Use dependency injection for services

### Models
- Primary model struct
- Request/Response specific structs
- Validation tags
- Use `validate` tags for field validation

### Repositories
- Implement data access layer
- Use `BaseRepository` interface
- Handle database/storage interactions
- Implement common CRUD methods

### Services
- Implement business logic
- Use `BaseService` interface
- Coordinate between repositories
- Apply validation and business rules

## Interface Design

### Repository Interface
```go
type DomainRepository interface {
    repository.BaseRepository[DomainModel, DomainResponse]
    
    // Domain-specific methods
    FindBySpecificField(ctx context.Context, value string) ([]DomainResponse, error)
}
```

### Service Interface
```go
type DomainService interface {
    repository.BaseService[DomainModel, DomainResponse]
    
    // Domain-specific methods
    GetBySpecificField(ctx context.Context, value string) ([]DomainResponse, error)
}
```

## Validation
- Use centralized `validator` package
- Apply struct-level and field-level validations
- Use `validate` tags

## Error Handling
- Use centralized `errors` package
- Create domain-specific error types
- Wrap and contextualize errors

## Best Practices
1. Keep handlers thin
2. Push logic to services
3. Use interfaces for dependency injection
4. Apply consistent error handling
5. Use generic base interfaces
6. Validate input at service layer

## Example Implementation

### models/user.go
```go
type User struct {
    ID        uuid.UUID `json:"id" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    Name      string    `json:"name" validate:"required,min=2,max=50"`
}

type UserCreate struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### repositories/user_repository.go
```go
type userRepository struct {
    client *supabase.Client
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    // Implementation
}
```

### services/user_service.go
```go
func (s *userService) CreateUser(ctx context.Context, user *models.UserCreate) error {
    // Validate input
    if err := validator.ValidateStruct(user); err != nil {
        return errors.Wrap(err, errors.ErrValidation, "Invalid user data")
    }
    
    // Business logic
    // ...
}
``` 