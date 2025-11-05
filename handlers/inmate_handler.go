package handlers

import (
	"net/http"
	"sample/go-gcp/inmate"

	"github.com/gin-gonic/gin"
)

type InmateHandler struct {
	inmateService inmate.InmateService
}

func NewInmate(inmateService inmate.InmateService) InmateHandler {
	return InmateHandler{inmateService}
}

func (handler *InmateHandler) GetInmates(c *gin.Context) {
	c.JSON(http.StatusOK, handler.inmateService.GetInmates())
}
