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

func (h *Handler) registerUser(c *gin.Context) {
	userToReg := &model.UserToRegister{}
	err := c.BindJSON(userToReg)
	if err != nil {
		log.Printf("c.BindJSON: %v", err)
		h.writeErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.svc.RegisterUser(c, userToReg)
	if err != nil {
		log.Printf("h.svc.RegisterUser: %v", err)
		h.writeErrorResponse(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("uuid.Parse: %v", err)
		h.writeErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	user, err := h.storage.GetUserByID(c, id)
	if err != nil {
		log.Printf("h.storage.GetUserByID: %v", err)

		if errors.Is(err, store.ErrNotFound) {
			h.writeErrorResponse(c, http.StatusNotFound, "user not found")
			return
		}

		h.writeErrorResponse(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
