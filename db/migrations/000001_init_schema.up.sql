CREATE TABLE IF NOT EXISTS teams(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    team_id INT NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
);

CREATE TABLE IF NOT EXISTS pull_requests(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    author_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(10) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
    need_more_reviewers BOOLEAN NOT NULL DEFAULT TRUE,
);

CREATE TABLE IF NOT EXISTS pr_reviewers(
    reviewer_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pr_id INT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    PRIMARY KEY (reviewer_id, pr_id)
);

CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);
CREATE INDEX idx_pr_reviewers_pr_id ON pr_reviewers(pr_id);
