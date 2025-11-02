package user_handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/system-education/internal/user/dtos"
	user_mapper "github.com/williamkoller/system-education/internal/user/mapper"
	port_user_handler "github.com/williamkoller/system-education/internal/user/port/handler"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
	port_user_usecase "github.com/williamkoller/system-education/internal/user/port/usecase"
)

type UserHandler struct {
	usecase port_user_usecase.UserUsecase
}

func NewUserHandler(usecase port_user_usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

var _ port_user_handler.UserHandler = &UserHandler{}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input dtos.AddUserDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	user, err := h.usecase.Create(input)

	if err != nil {
		if errors.Is(err, port_user_repository.ErrUserAlreadyExists) {
			c.Status(http.StatusConflict)
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	res := user_mapper.ToUser(user)
	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) FindAllUsers(c *gin.Context) {
	users, err := h.usecase.FindAll()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	resp := user_mapper.ToUsers(users)
	c.JSON(http.StatusOK, resp)
}
