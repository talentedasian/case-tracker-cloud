package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"sample/go-gcp/inmate"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	putInmateRequestsCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "inmate_requests_total",
			Help: "Total number of inmate API requests",
		},
	)
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

func (handler *InmateHandler) PutInmate(c *gin.Context) {
	putInmateRequestsCounter.Inc()

	var inmate inmate.Inmate
	if err := c.BindJSON(&inmate); err != nil {
		jsonBody := parseRawJson(c.Request.Body)
		slog.Error("failed to parse request body to inmate", "req", jsonBody, "err", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON to expected inmate"})
		return
	}

	if putErr := handler.inmateService.PutInmate(inmate); putErr != nil {
		slog.Error("failed to put item to dynamodb", "err", putErr.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "operation to put item to dynamodb failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"inmate": inmate})
}

func (handler *InmateHandler) Attempt(c *gin.Context) {
	var body map[string]string

	if err := c.BindJSON(&body); err != nil {
		jsonBody := parseRawJson(c.Request.Body)
		slog.Error("failed to parse request body to map", "req", jsonBody, "err", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON to expected map object"})
		return
	}

	if err := handler.inmateService.Attempt(body["id"], body["reason"]); err != nil {
		slog.Error("failed to save attempt", "err", err.Error())
	}

}

func parseRawJson(body io.ReadCloser) string {
	bodyAsByteArray, err := io.ReadAll(body)

	if err != nil {
		slog.Debug("parsing request body to raw json failed", "err", err.Error())

		return "empty json, parsing failed"
	}

	return string(bodyAsByteArray)
}
