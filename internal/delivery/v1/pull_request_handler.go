package v1

import (
	"avito-internship/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) pullRequestCreate(c *gin.Context) {
	var req CreatePullRequestReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.PullRequestCreate(c.Request.Context(), toUseCasePullRequestCreate(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, toDeliveryPullRequestCreate(res))
}

func (h *Handler) pullRequestMerge(c *gin.Context) {
	var req PullRequestMergeReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.PullRequestMerge(c.Request.Context(), toUseCasePullRequestMerge(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliveryPullRequestMerge(res))
}

func (h *Handler) reviewerReassign(c *gin.Context) {
	var req PullRequestReassignReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.ReviewerReassign(c.Request.Context(), toUseCaseReviewerReassignReq(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliveryReviewerReassignRes(res))
}
