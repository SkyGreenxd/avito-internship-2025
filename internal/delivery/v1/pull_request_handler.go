package v1

import (
	"avito-internship/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) pullRequestCreate(c *gin.Context) {
	var req CreatePullRequestReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.Wrap(err.Error(), e.ErrInvalidRequestBody))
		return
	}

	res, err := h.prUC.PullRequestCreate(c.Request.Context(), toUseCaseCreatePullRequestReq(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, toDeliveryCreatePullRequestRes(res))
}

func (h *Handler) pullRequestMerge(c *gin.Context) {
	var req PullRequestMergeReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.Wrap(err.Error(), e.ErrInvalidRequestBody))
		return
	}

	res, err := h.prUC.PullRequestMerge(c.Request.Context(), toUseCasePullRequestMergeReq(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliveryPullRequestMergeRes(res))
}

func (h *Handler) reviewerReassign(c *gin.Context) {
	var req PullRequestReassignReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.Wrap(err.Error(), e.ErrInvalidRequestBody))
		return
	}

	res, err := h.prUC.ReviewerReassign(c.Request.Context(), toUseCasePullRequestReassignReq(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliveryPullRequestReassignRes(res))
}
