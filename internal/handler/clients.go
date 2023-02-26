package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/pkg/errors"
)

func (h *Handler) registerClient(c *gin.Context) {
	clientToReg := &model.ClientToRegister{}
	err := c.BindJSON(clientToReg)
	if err != nil {
		log.Printf("c.BindJSON: %v", err)
		h.writeErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	client, err := h.svc.RegisterClient(c, clientToReg)
	if err != nil {
		log.Printf("h.svc.RegisterClient: %v", err)
		h.writeErrorResponse(c, http.StatusInternalServerError, "failed to create client")
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *Handler) getClientByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("uuid.Parse: %v", err)
		h.writeErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	client, err := h.storage.GetClientByID(c, id)
	if err != nil {
		log.Printf("h.storage.GetClientByID: %v", err)

		if errors.Is(err, store.ErrNotFound) {
			h.writeErrorResponse(c, http.StatusNotFound, "client not found")
			return
		}

		h.writeErrorResponse(c, http.StatusInternalServerError, "failed to get client")
		return
	}

	c.IndentedJSON(http.StatusOK, client)
}
