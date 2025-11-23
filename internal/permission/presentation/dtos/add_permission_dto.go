package permission_dtos

type AddPermissionDto struct {
	UserID      string   `json:"user_id"`
	Modules     []string `json:"modules"`
	Actions     []string `json:"actions"`
	Level       string   `json:"level"`
	Description string   `json:"description"`
}
