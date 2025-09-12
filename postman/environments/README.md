# Itsrama Portfolio Backend Postman Environments

## Overview
These Postman environment files provide configuration for different deployment stages.

## Environments

### 1. Local Development (`local.json`)
- Base URL: `http://localhost:8080/api/v1`
- Ideal for local development and testing
- Uses local database connection
- Placeholder for development credentials

### 2. Staging Environment (`staging.json`)
- Base URL: `https://staging-api.itsrama.com/api/v1`
- Used for pre-production testing
- Connects to staging database
- Separate credentials from production

### 3. Production Environment (`production.json`)
- Base URL: `https://api.itsrama.com/api/v1`
- Used for final production testing
- Disabled database URL for security
- Requires specific production credentials

## Common Variables

Each environment includes:
- `BASE_URL`: API endpoint
- `BEARER_TOKEN`: Authentication token
- Resource-specific IDs:
  - `TECH_STACK_ID`
  - `PROJECT_ID`
  - `EXPERIENCE_ID`
- `USER_EMAIL`: User identifier
- `DATABASE_URL`: Database connection (disabled in production)
- `SUPABASE_PROJECT_ID`: Supabase project identifier

## Usage

1. Import the environment in Postman
2. Select the appropriate environment
3. Fill in secret values like `BEARER_TOKEN`
4. Use environment-specific variables in requests

## Security Notes

- Never commit real credentials to version control
- Use Postman's secret variable type for sensitive information
- Rotate tokens and credentials regularly
- Limit access to production environment

## Contributing

- Keep environments minimal and purpose-specific
- Use placeholder values for shared repositories
- Document any environment-specific configurations

## License
[Your Project License] 