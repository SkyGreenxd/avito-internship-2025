DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS pr_reviewers;

DROP INDEX IF EXISTS idx_pr_reviewers_reviewer_id;
DROP INDEX IF EXISTS idx_pr_reviewers_pr_id;