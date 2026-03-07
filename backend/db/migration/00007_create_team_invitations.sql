-- +goose Up
-- invitation_status: pending, accepted, declined, expired (используется в колонке status ниже)
CREATE TABLE team_invitations (
    id CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    team_id CHAR(36) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    invited_by CHAR(36) NOT NULL,
    status ENUM('pending', 'accepted', 'declined', 'expired') NOT NULL DEFAULT 'pending',
    token VARCHAR(64) NOT NULL,
    expires_at DATETIME(3) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_team_invitations_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_team_invitations_invited_by FOREIGN KEY (invited_by) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_team_invitations_team_id ON team_invitations(team_id);
CREATE INDEX idx_team_invitations_email ON team_invitations(email);
CREATE INDEX idx_team_invitations_status ON team_invitations(status);
CREATE UNIQUE INDEX idx_team_invitations_token ON team_invitations(token);

-- +goose Down
DROP INDEX idx_team_invitations_token ON team_invitations;
DROP INDEX idx_team_invitations_status ON team_invitations;
DROP INDEX idx_team_invitations_email ON team_invitations;
DROP INDEX idx_team_invitations_team_id ON team_invitations;
DROP TABLE IF EXISTS team_invitations;
