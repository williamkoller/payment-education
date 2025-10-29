package dtos

type AddUserDto struct {
	Name        string   `json:"name" validate:"required,min=2,max=100" example:"John"`
	Surname     string   `json:"surname" validate:"required,min=2,max=100" example:"Doe"`
	Nickname    string   `json:"nickname" validate:"required,min=2,max=50" example:"johnd"`
	Age         int32    `json:"age" validate:"required,gte=1,lte=130" example:"30"`
	Email       string   `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password    string   `json:"password" validate:"required,min=8" example:"strongPassword123"`
	Roles       []string `json:"roles" validate:"required,min=1,dive,required" example:"[\"admin\"]"`
	Permissions []string `json:"permissions" validate:"required,min=1,dive,required" example:"[\"read\",\"write\"]"`
}
