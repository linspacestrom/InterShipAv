CREATE TABLE IF NOT EXISTS "user" (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    FOREIGN KEY (team_name) REFERENCES team(team_name)
);