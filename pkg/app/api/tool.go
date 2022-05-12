package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerToolV1Api(g *gin.Engine) {
	group := g.Group("/api/v1/")
	group.GET("/test", TT)
}

type TResp struct {
	Message string `json:"message"`
}

// TT
// @Summary  get tool component info
// @Description tool component info api
// @Tags test
// @Accept  json
// @Product json
// @Param namespace query string true "namespace"
// @Param toolName path string true "tool name"
// @Success 200 {object} TResp
// @Failure 400 {object} errs.PraticeException
// @Failure 401 {object} errs.PraticeException
// @Failure 403 {object} errs.PraticeException
// @Failure 500 {object} errs.PraticeException
// @Router /api/v1/test [GET]
func TT(c *gin.Context) {
	c.JSON(http.StatusOK, &TResp{Message: "ok"})
	return
}
