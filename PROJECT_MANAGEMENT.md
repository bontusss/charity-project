# Project Management System

This document describes the project management system for the Chinedu Onyeizu Foundation, including API endpoints, authentication, and usage examples.

## Overview

The project management system allows authenticated admins to create, update, and delete projects, as well as manage project details including "before" and "after" sections, and upload images.

## Authentication

All project management endpoints (create, update, delete) require admin authentication. Public endpoints (view, list) are accessible without authentication.

## API Endpoints

### Project Management (Admin Only - Requires Authentication)

#### Create Project
```
POST /api/projects
Content-Type: application/json

{
  "name": "Project Name"
}
```

#### Update Project
```
PUT /api/projects/:id
Content-Type: application/json

{
  "name": "Updated Project Name"
}
```

#### Delete Project
```
DELETE /api/projects/:id
```

#### Create/Update Project Before Section
```
POST /api/projects/:id/before
PUT /api/projects/:id/before
Content-Type: application/json

{
  "body": "Project description before implementation",
  "estimated_target": "50000",
  "video_link": "https://youtube.com/watch?v=..."
}
```

#### Create/Update Project After Section
```
POST /api/projects/:id/after
PUT /api/projects/:id/after
Content-Type: application/json

{
  "body": "Project description after implementation",
  "project_cost": "45000",
  "video_link": "https://youtube.com/watch?v=..."
}
```

#### Upload Project Image
```
POST /api/projects/:id/images
Content-Type: multipart/form-data

Form fields:
- image: File (jpg, jpeg, png, gif)
- phase: String (e.g., "before", "after", "during")
```

#### Delete Project Image
```
DELETE /api/projects/images/:image_id
```

### Public Endpoints (No Authentication Required)

#### List All Projects
```
GET /api/projects
```

#### Get Project Details
```
GET /api/projects/:id
```

#### Get Project Before Section
```
GET /api/projects/:id/before
```

#### Get Project After Section
```
GET /api/projects/:id/after
```

#### List Project Images
```
GET /api/projects/:id/images
```

#### List Project Images by Phase
```
GET /api/projects/:id/images/phase?phase=before
```

## Image Upload

### Supported Formats
- JPG/JPEG
- PNG
- GIF

### File Storage
Images are stored in the `/static/images/` directory with the following naming convention:
```
project_{project_id}_{phase}_{timestamp}.{extension}
```

### File Size Limits
The system uses Gin's default file size limits. For production, consider configuring appropriate limits.

## Database Schema

### Projects Table
```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Project Before Table
```sql
CREATE TABLE project_before (
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    estimated_target INTEGER NOT NULL,
    video_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (project_id)
);
```

### Project After Table
```sql
CREATE TABLE project_after (
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    body TEXT,
    project_cost INTEGER,
    video_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (project_id)
);
```

### Project Images Table
```sql
CREATE TABLE project_images (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    phase VARCHAR(50) NOT NULL,
    image_link VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Usage Examples

### Creating a Complete Project

1. **Create the project:**
```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=your_jwt_token" \
  -d '{"name": "School Renovation Project"}'
```

2. **Add project before details:**
```bash
curl -X POST http://localhost:8080/api/projects/1/before \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=your_jwt_token" \
  -d '{
    "body": "The school building is in poor condition with leaking roofs and broken windows.",
    "estimated_target": "75000",
    "video_link": "https://youtube.com/watch?v=before_video"
  }'
```

3. **Upload before images:**
```bash
curl -X POST http://localhost:8080/api/projects/1/images \
  -H "Cookie: auth_token=your_jwt_token" \
  -F "image=@before_photo1.jpg" \
  -F "phase=before"
```

4. **After project completion, add after details:**
```bash
curl -X POST http://localhost:8080/api/projects/1/after \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=your_jwt_token" \
  -d '{
    "body": "The school has been completely renovated with new roofs, windows, and fresh paint.",
    "project_cost": "72000",
    "video_link": "https://youtube.com/watch?v=after_video"
  }'
```

5. **Upload after images:**
```bash
curl -X POST http://localhost:8080/api/projects/1/images \
  -H "Cookie: auth_token=your_jwt_token" \
  -F "image=@after_photo1.jpg" \
  -F "phase=after"
```

### Viewing Projects

**List all projects:**
```bash
curl http://localhost:8080/api/projects
```

**Get specific project:**
```bash
curl http://localhost:8080/api/projects/1
```

**Get project images:**
```bash
curl http://localhost:8080/api/projects/1/images
```

## Admin Dashboard

The admin dashboard provides a web interface for managing projects:

1. **Access:** Login at `/admin/login` and navigate to `/dashboard`
2. **Create Projects:** Use the "Create New Project" button
3. **View Projects:** See all projects in a list format
4. **Manage Projects:** View, edit, or delete projects using the action buttons

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid authentication
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

Error responses include a JSON object with an `error` field:
```json
{
  "error": "Error description"
}
```

## Security Considerations

1. **Authentication:** All admin operations require valid JWT tokens
2. **File Upload:** Only image files are allowed, with extension validation
3. **Input Validation:** All inputs are validated before processing
4. **SQL Injection:** Uses parameterized queries via sqlc
5. **File System:** Images are stored in a controlled directory

## Future Enhancements

1. **Image Resizing:** Automatic thumbnail generation
2. **File Compression:** Optimize image file sizes
3. **Bulk Operations:** Upload multiple images at once
4. **Project Categories:** Add project categorization
5. **Progress Tracking:** Track project completion status
6. **Donation Integration:** Link projects to donation campaigns 