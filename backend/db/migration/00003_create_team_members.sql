-- +goose Up
CREATE TABLE team_members (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    team_id CHAR(36) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_team_members_user_team (user_id, team_id),
    CONSTRAINT fk_team_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_team_members_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);

-- +goose Down
DROP INDEX idx_team_members_team_id ON team_members;
DROP INDEX idx_team_members_user_id ON team_members;
DROP TABLE IF EXISTS team_members;
