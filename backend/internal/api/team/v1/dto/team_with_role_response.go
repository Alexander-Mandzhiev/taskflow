package dto

// TeamWithRoleResponse — команда с ролью текущего пользователя (для списка «мои команды»).
type TeamWithRoleResponse struct {
	TeamResponse
	Role string `json:"role"`
}
