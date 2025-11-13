package v1

import "github.com/gin-gonic/gin"

type Handler struct {
	// TODO: сервисы/юзкейсы подключить
	middleware *Middleware
}

// TODO: передать сервисы
func NewHandler(middleware *Middleware) *Handler {
	return &Handler{
		// TODO: usecases
		middleware: middleware,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	v1.Use(h.middleware.ErrorMiddleware())
	{
		team := v1.Group("/team")
		{
			team.POST("/add", h.addTeam)
			team.GET("/get", h.getTeam)
		}

		users := v1.Group("/users")
		{
			users.POST("/setIsActive", h.middleware.AdminMiddleware(), h.setIsActive)
			users.GET("/getReview", h.getReview)
		}

		pullRequest := v1.Group("/pullRequest")
		{
			pullRequest.POST("/create", h.pullRequestCreate)
			pullRequest.POST("/merge", h.pullRequestMerge)
			pullRequest.POST("/reassign", h.reviewerReassign)
		}
	}
}
