package v1

import (
	"avito-internship/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userUC     usecase.UserUC
	teamUC     usecase.TeamUC
	prUC       usecase.PullRequestUC
	middleware *Middleware
}

func NewHandler(userUc usecase.UserUC, teamUc usecase.TeamUC, prUC usecase.PullRequestUC, middleware *Middleware) *Handler {
	return &Handler{
		userUC:     userUc,
		teamUC:     teamUc,
		prUC:       prUC,
		middleware: middleware,
	}
}

func (h *Handler) Init(r *gin.Engine) {
	r.Use(h.middleware.ErrorMiddleware())

	team := r.Group("/team")
	{
		team.POST("/add", h.addTeam)
		team.GET("/get", h.getTeam)
		team.POST("/deactivate", h.deactivateMembers)
	}

	users := r.Group("/users")
	{
		users.POST("/setIsActive", h.middleware.AdminMiddleware(), h.setIsActive)
		users.GET("/getReview", h.getReview)
	}

	pullRequest := r.Group("/pullRequest")
	{
		pullRequest.POST("/create", h.pullRequestCreate)
		pullRequest.POST("/merge", h.pullRequestMerge)
		pullRequest.POST("/reassign", h.reviewerReassign)
	}

}
