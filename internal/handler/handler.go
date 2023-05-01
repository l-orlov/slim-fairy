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
		authUsers := auth.Group("/users")
		{
			authUsers.POST("/sign-up", h.registerUser)
			authUsers.POST("/sign-in", h.signInUser)
		}
		//authNutritionists := auth.Group("/nutritionists")
		{
			//auth.POST("/sign-up", h.registerUser)
			//auth.POST("/sign-in", h.SignIn)
		}
	}

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/:id", h.getUserByID)
		}
		nutritionists := api.Group("/nutritionists")
		{
			nutritionists.GET("/:id", h.getNutritionistByID)
		}
	}

	return CORS(router)
}
