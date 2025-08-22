-- Projects table (main project info)
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Project Before (1-to-1 with projects)
CREATE TABLE project_before (
    project_id INT PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    estimated_target INTEGER NOT NULL,
    video_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Project After (1-to-1 with projects, nullable until completed)
CREATE TABLE project_after (
    project_id INT PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    body TEXT,
    project_cost INTEGER,
    video_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Images table (can store multiple images for before/after)
CREATE TABLE project_images (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    phase VARCHAR(10) NOT NULL CHECK (phase IN ('before', 'after')), -- to distinguish which part
    image_link VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
