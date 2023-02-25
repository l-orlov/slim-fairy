package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/pkg/errors"
)

func (h *Handler) registerNutritionist(c *gin.Context) {
}

func (h *Handler) getNutritionistByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("uuid.Parse: %v", err)
		h.writeErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	client, err := h.storage.GetNutritionistByID(c, id)
	if err != nil {
		log.Printf("h.storage.GetNutritionistByID: %v", err)

		if errors.Is(err, store.ErrNotFound) {
			h.writeErrorResponse(c, http.StatusNotFound, "nutritionist not found")
			return
		}

		h.writeErrorResponse(c, http.StatusInternalServerError, "failed to get nutritionist")
		return
	}

	c.IndentedJSON(http.StatusOK, client)
}
