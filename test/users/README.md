# Users Test Suite

## Overview
This test suite provides comprehensive testing for the Users module in the Cultour backend, covering services, repositories, and handlers.

## Test Coverage

### Services
- User Service
  - User creation
  - User retrieval
  - User update
  - User deletion
  - User listing and searching

- User Profile Service
  - Profile creation
  - Profile retrieval
  - Profile update
  - Profile deletion
  - Profile avatar and identity updates

- User Badge Service
  - Badge assignment
  - Badge removal
  - Badge retrieval
  - Badge listing and counting

### Handlers
- User Handler
  - HTTP endpoints for user management
  - Request validation
  - Response formatting
  - Authentication context handling

- User Profile Handler
  - HTTP endpoints for user profile management
  - Multipart form handling
  - File upload validation
  - Profile-related operations

- User Badge Handler
  - HTTP endpoints for user badge management
  - Badge assignment and removal
  - Badge listing and counting

## Testing Strategies

### Mocking
- Use `testify/mock` for creating mock implementations
- Mock external dependencies like repositories and services
- Simulate various scenarios including success and error cases

### Test Cases
- Happy path scenarios
- Edge cases
- Error handling
- Input validation
- Pagination and filtering
- Authentication and authorization checks

## Running Tests

### Prerequisites
- Go 1.21 or higher
- `testify` library for assertions and mocking

### Execute Tests
```bash
# Run all user-related tests
go test ./test/users/...

# Run specific test suites
go test ./test/users/services
go test ./test/users/handlers
```

## Best Practices
- Keep tests focused and isolated
- Use meaningful test names
- Cover both positive and negative scenarios
- Validate input validation, business logic, and error handling
- Maintain high test coverage

## Continuous Integration
These tests are integrated into our CI/CD pipeline to ensure code quality and prevent regressions. 