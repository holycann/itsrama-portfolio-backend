# Itsrama Portfolio Backend Postman Collections

## Overview
These Postman collections provide comprehensive API documentation and testing for the Itsrama Portfolio Backend.

## Prerequisites
- [Postman](https://www.postman.com/downloads/)
- Bearer Token for Authentication

## Collections

### 1. Base Collection
- Contains global configurations
- Shared authentication settings
- Common pre-request and test scripts

### 2. Tech Stack Collection
- CRUD operations for tech stacks
- Bulk create, update, and delete
- Category-based filtering
- Search functionality

### 3. Experience Collection
- CRUD operations for professional experience
- File upload support (logo and images)
- Bulk create, update, and delete
- Company-based filtering
- Search functionality

### 4. Project Collection
- CRUD operations for projects
- Bulk create, update, and delete
- Category-based filtering
- Search functionality

## Usage

1. Import the collections in Postman
2. Set environment variables:
   - `BASE_URL`: Your backend server URL (default: `http://localhost:8080/api/v1`)
   - `BEARER_TOKEN`: Your authentication token

## Merging Collections
Use [postman-combine-collections](https://www.npmjs.com/package/postman-combine-collections) to merge these collections:

```bash
npx postman-combine-collections --name Itsrama -f 'postman/*_collection.json' -o postman/itsrama_portfolio_backend_collection.json
```

## Contributing
- Ensure collections follow SOLID principles
- Keep collections DRY (Don't Repeat Yourself)
- Maintain single responsibility for each collection

## License
[Your Project License] 