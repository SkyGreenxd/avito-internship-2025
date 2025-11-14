package v1

import (
	"avito-internship/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addTeam(c *gin.Context) {
	var team TeamAddReq
	if err := c.ShouldBindJSON(&team); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.teamUC.AddTeam(c.Request.Context(), toUseCaseTeamAddReq(team))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, toDeliveryTeamAddRes(res))
}

func (h *Handler) getTeam(c *gin.Context) {
	var req GetTeamQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(e.ErrInvalidRequestBody)
		return
	}

	res, err := h.teamUC.GetTeam(c.Request.Context(), req.TeamName)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toDeliveryGetTeamRes(res))
}
