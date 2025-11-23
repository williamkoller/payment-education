package permission_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	permission_mapper "github.com/williamkoller/system-education/internal/permission/application/mapper"
	port_permission_handler "github.com/williamkoller/system-education/internal/permission/port/handler"
	port_permission_usecase "github.com/williamkoller/system-education/internal/permission/port/usecase"
	permission_dtos "github.com/williamkoller/system-education/internal/permission/presentation/dtos"
)

type PermissionHandler struct {
	usecase port_permission_usecase.PermissionUsecase
}

func NewPermissionHandler(usecase port_permission_usecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{
		usecase: usecase,
	}
}

var _ port_permission_handler.PermissionHandler = &PermissionHandler{}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var input permission_dtos.AddPermissionDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	p, err := h.usecase.Create(input)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	resp := permission_mapper.ToPermission(p)

	c.JSON(http.StatusCreated, resp)
}

func (h *PermissionHandler) FindAllPermission(c *gin.Context) {
	permissions, err := h.usecase.FindAll()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	resp := permission_mapper.ToPermissions(permissions)

	c.JSON(http.StatusOK, resp)
}
