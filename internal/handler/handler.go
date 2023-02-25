package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/slim-fairy/internal/service"
	"github.com/l-orlov/slim-fairy/internal/store"
)

type Handler struct {
	storage *store.Storage
	svc     *service.Service
}

func New(
	storage *store.Storage,
	svc *service.Service,
) http.Handler {
	h := &Handler{
		storage: storage,
		svc:     svc,
	}

	router := gin.Default()

	auth := router.Group("/auth")
	{
		authClients := auth.Group("/clients")
		{
			authClients.POST("/sign-up", h.registerClient)
			authClients.POST("/sign-in", h.signInClient)
		}
		//authNutritionists := auth.Group("/nutritionists")
		{
			//auth.POST("/sign-up", h.registerClient)
			//auth.POST("/sign-in", h.SignIn)
		}
	}

	api := router.Group("/api")
	{
		clients := api.Group("/clients")
		{
			clients.GET("/:id", h.getClientByID)
		}
		nutritionists := api.Group("/nutritionists")
		{
			nutritionists.GET("/:id", h.getNutritionistByID)
		}
	}

	return CORS(router)
}
