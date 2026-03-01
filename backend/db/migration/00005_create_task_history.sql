-- +goose Up
CREATE TABLE task_history (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    task_id CHAR(36) NOT NULL,
    changed_by CHAR(36) NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    changed_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_task_history_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_task_history_changed_by FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_task_history_task_id ON task_history(task_id);
CREATE INDEX idx_task_history_changed_at ON task_history(changed_at);

-- +goose Down
DROP INDEX idx_task_history_changed_at ON task_history;
DROP INDEX idx_task_history_task_id ON task_history;
DROP TABLE IF EXISTS task_history;
