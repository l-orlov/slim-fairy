package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/service"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/pkg/errors"
)

func (h *Handler) signInUser(c *gin.Context) {
	user := &model.UserToSignIn{}
	err := c.BindJSON(user)
	if err != nil {
		log.Printf("c.BindJSON: %v", err)
		h.writeErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := h.svc.SignInUser(c, user)
	if err != nil {
		log.Printf("h.svc.SignInUser: %v", err)

		msg := "failed to sign-in user"
		if errors.Is(err, store.ErrNotFound) {
			h.writeErrorResponse(c, http.StatusNotFound, msg)
			return
		}
		if errors.Is(err, service.ErrWrongPassword) {
			h.writeErrorResponse(c, http.StatusUnauthorized, "wrong password")
			return
		}

		h.writeErrorResponse(c, http.StatusInternalServerError, msg)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": token,
	})
}
