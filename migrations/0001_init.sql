CREATE TABLE IF NOT EXISTS teams (
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) REFERENCES teams(name),
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'OPEN',
    assigned_reviewers JSONB,
    merged_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);