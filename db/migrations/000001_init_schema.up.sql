CREATE TABLE IF NOT EXISTS teams(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    team_id INT REFERENCES teams(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS statuses(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id VARCHAR(50) PRIMARY KEY,
    name TEXT NOT NULL,
    author_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status_id INT NOT NULL REFERENCES statuses(id) DEFAULT 1 ON DELETE RESTRICT,
    need_more_reviewers BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pr_reviewers(
    reviewer_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pr_id VARCHAR(50) NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    PRIMARY KEY (reviewer_id, pr_id)
);

CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);
CREATE INDEX idx_pr_reviewers_pr_id ON pr_reviewers(pr_id);
CREATE INDEX idx_users_team_active ON users(team_id, is_active);