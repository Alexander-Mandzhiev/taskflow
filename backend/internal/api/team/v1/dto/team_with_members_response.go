package dto

// TeamWithMembersResponse — команда с участниками в ответе API.
type TeamWithMembersResponse struct {
	Team    TeamResponse     `json:"team"`
	Members []MemberResponse `json:"members"`
}
