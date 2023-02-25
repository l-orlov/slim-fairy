package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	detailBusinessError = "check request parameters"
	detailServerError   = "something went wrong"
)

func (h *Handler) writeErrorResponse(c *gin.Context, statusCode int, msg string) {
	detail := detailServerError
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		detail = detailBusinessError
	}

	c.AbortWithStatusJSON(statusCode, gin.H{
		"message": msg,
		"detail":  detail,
	})
}
