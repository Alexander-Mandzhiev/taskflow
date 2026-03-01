-- +goose Up
CREATE TABLE task_comments (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    task_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL DEFAULT NULL,
    CONSTRAINT fk_task_comments_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_task_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_task_comments_task_deleted ON task_comments(task_id, deleted_at);

-- +goose Down
DROP INDEX idx_task_comments_task_deleted ON task_comments;
DROP TABLE IF EXISTS task_comments;
