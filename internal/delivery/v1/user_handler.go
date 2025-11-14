package v1

import (
	"avito-internship/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) setIsActive(c *gin.Context) {
	var req SetIsActiveReq
	if err := c.ShouldBind(&req); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.userUC.SetIsActive(c.Request.Context(), toUseCaseSetIsActiveReq(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliverySetIsActiveRes(res))
}

func (h *Handler) getReview(c *gin.Context) {
	var req GetReviewQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.userUC.GetReview(c.Request.Context(), req.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliveryGetReviewRes(res))
}
