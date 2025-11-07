package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sample/go-gcp/inmate"
	"strings"

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
	var body map[string]any

	if err := c.BindJSON(&body); err != nil {
		jsonBody := parseRawJson(c.Request.Body)
		slog.Error("failed to parse request body to map", "req", jsonBody, "err", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON to expected map object"})
		return
	}

	attempts, ok := body["attempts"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field `attempts` passed in not a number"})

		return
	}

	if err := handler.inmateService.Attempt(fmt.Sprint(body["id"]),
		fmt.Sprint(body["reason"]), int8(attempts)); err != nil {
		slog.Error("failed to save attempt", "err", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save inmate attempt"})
	}

}

func (handler *InmateHandler) GetInmateAttempts(c *gin.Context) {
	inmateId := c.Param("inmateId")

	if strings.TrimSpace(inmateId) == "" {
		slog.Debug("attempt to get inmate attemtes: received an empty inmate id")

		c.JSON(http.StatusBadRequest, gin.H{"error": "inmate id request parameter cannot be empty"})

		return
	}

	inmates, err := handler.inmateService.GetAttempts(inmateId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, inmates)
}

func parseRawJson(body io.ReadCloser) string {
	bodyAsByteArray, err := io.ReadAll(body)

	if err != nil {
		slog.Debug("parsing request body to raw json failed", "err", err.Error())

		return "empty json, parsing failed"
	}

	return string(bodyAsByteArray)
}
