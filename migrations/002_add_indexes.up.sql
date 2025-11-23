CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_prs_author ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_prs_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pr_reviewers(user_id);