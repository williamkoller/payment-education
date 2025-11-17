package dtos


type UpdateUserDto struct {
	Name     *string `json:"name" validate:"min=2,max=100" example:"John"`
	Surname  *string `json:"surname" validate:"min=2,max=100" example:"Doe"`
	Nickname *string `json:"nickname" validate:"min=2,max=50" example:"johnd"`
	Age      *int32  `json:"age" validate:"gte=1,lte=130" example:"30"`
	Email    *string `json:"email" validate:"email" example:"john.doe@example.com"`
	Password *string `json:"password" validate:"min=8" example:"strongPassword123"`
}
