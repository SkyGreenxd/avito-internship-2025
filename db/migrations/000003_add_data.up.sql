INSERT INTO teams(name)
VALUES ('alpha_team');

INSERT INTO users(id, name, is_active, team_id)
VALUES
    ('u1', 'Alice', true, 1),
    ('u2', 'Bob', true, 1),
    ('u3', 'Carol', true, 1),
    ('u4', 'Dave', true, 1),
    ('u5', 'Eve', true, 1),
    ('u6', 'Frank', true, 1),
    ('u7', 'Grace', true, 1),
    ('u8', 'Hank', true, 1);

INSERT INTO pull_requests(id, name, author_id, status_id, need_more_reviewers, created_at)
VALUES
    ('pr-1001', 'Add login feature', 'u1', 1, TRUE, NOW()),
    ('pr-1002', 'Fix bug #42', 'u2', 1, TRUE, NOW() - INTERVAL '1 day'),
    ('pr-1003', 'Improve performance', 'u3', 1, TRUE, NOW() - INTERVAL '2 days'),
    ('pr-1004', 'Refactor module', 'u4', 1, TRUE, NOW() - INTERVAL '3 days'),
    ('pr-1005', 'Update README', 'u5', 1, TRUE, NOW() - INTERVAL '4 days'),
    ('pr-1006', 'Add tests', 'u6', 1, TRUE, NOW() - INTERVAL '5 days');

INSERT INTO pr_reviewers(reviewer_id, pr_id) VALUES
    ('u2', 'pr-1001'),
    ('u3', 'pr-1001'),
    ('u1', 'pr-1002'),
    ('u3', 'pr-1002'),
    ('u1', 'pr-1003'),
    ('u4', 'pr-1003'),
    ('u1', 'pr-1004'),
    ('u2', 'pr-1004'),
    ('u2', 'pr-1005'),
    ('u3', 'pr-1005'),
    ('u1', 'pr-1006'),
    ('u2', 'pr-1006');
