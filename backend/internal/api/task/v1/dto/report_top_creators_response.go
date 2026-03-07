package dto

// TopCreatorResponse — топ создатель задач в команде за период.
type TopCreatorResponse struct {
	TeamID       string `json:"team_id"`
	UserID       string `json:"user_id"`
	Rank         int    `json:"rank"`
	CreatedCount int64  `json:"created_count"`
}

// TopCreatorsListResponse — список топ создателей по командам.
type TopCreatorsListResponse struct {
	Items []TopCreatorResponse `json:"items"`
}
