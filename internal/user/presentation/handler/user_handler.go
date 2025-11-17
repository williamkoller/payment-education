package user_handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	user_mapper "github.com/williamkoller/system-education/internal/user/application/mapper"
	portUserHandler "github.com/williamkoller/system-education/internal/user/port/handler"
	portUserRepository "github.com/williamkoller/system-education/internal/user/port/repository"
	portUserUsecase "github.com/williamkoller/system-education/internal/user/port/usecase"
	"github.com/williamkoller/system-education/internal/user/presentation/dtos"
)

type UserHandler struct {
	usecase portUserUsecase.UserUsecase
}

func NewUserHandler(usecase portUserUsecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

var _ portUserHandler.UserHandler = &UserHandler{}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input dtos.AddUserDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	user, err := h.usecase.Create(input)

	if err != nil {
		if errors.Is(err, portUserRepository.ErrUserAlreadyExists) {
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

func (h *UserHandler) FindByID(c *gin.Context) {
	idParams := c.Param("id")
	user, err := h.usecase.FindByID(idParams)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	resp := user_mapper.ToUser(user)
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Update(c *gin.Context) {
	idParams := c.Param("id")

	var input dtos.UpdateUserDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	user, err := h.usecase.Update(idParams, input)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	resp := user_mapper.ToUser(user)
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.usecase.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete user",
		})
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "user deleted successfully",
	})
}
