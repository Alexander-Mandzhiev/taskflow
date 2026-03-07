-- +goose Up
CREATE TABLE tasks (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'todo',
    assignee_id CHAR(36),
    team_id CHAR(36) NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    completed_at DATETIME(3) NULL DEFAULT NULL,
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    CONSTRAINT fk_tasks_assignee FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_tasks_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_tasks_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_tasks_team_id ON tasks(team_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_team_status ON tasks(team_id, status);
CREATE INDEX idx_tasks_team_status_updated ON tasks(team_id, status, updated_at);
CREATE INDEX idx_tasks_team_status_completed ON tasks(team_id, status, completed_at);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- +goose Down
DROP INDEX idx_tasks_deleted_at ON tasks;
DROP INDEX idx_tasks_team_status_completed ON tasks;
DROP INDEX idx_tasks_team_status_updated ON tasks;
DROP INDEX idx_tasks_team_status ON tasks;
DROP INDEX idx_tasks_assignee_id ON tasks;
DROP INDEX idx_tasks_status ON tasks;
DROP INDEX idx_tasks_team_id ON tasks;
DROP TABLE IF EXISTS tasks;
