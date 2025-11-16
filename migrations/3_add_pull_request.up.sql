CREATE TABLE IF NOT EXISTS pull_request (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL,
    status            TEXT NOT NULL,
    merged_at         TIMESTAMPTZ DEFAULT NULL,
    created_at        TIMESTAMPTZ DEFAULT NOW(),

    FOREIGN KEY (author_id) REFERENCES "user"(id)
);