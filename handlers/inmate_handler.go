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
	inmates, err := handler.inmateService.GetInmates()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, inmates)
}
