CREATE TABLE IF NOT EXISTS pr_reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_request(pull_request_id) ON DELETE CASCADE,
    reviewer_id     TEXT NOT NULL REFERENCES "user"(id),

    PRIMARY KEY (pull_request_id, reviewer_id)
);
