CREATE TABLE teams (
    name VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) REFERENCES teams(name),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'OPEN',
    assigned_reviewers JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP
);