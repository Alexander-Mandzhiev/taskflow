-- +goose Up
CREATE TABLE teams (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    CONSTRAINT fk_teams_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_teams_deleted_at ON teams(deleted_at);

-- +goose Down
DROP INDEX idx_teams_deleted_at ON teams;
DROP TABLE IF EXISTS teams;
