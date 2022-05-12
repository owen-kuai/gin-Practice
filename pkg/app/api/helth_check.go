package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerHealthCheckApi(g *gin.Engine) {
	group := g.Group("/health")
	group.GET("/liveness", liveness)
	group.GET("/readiness", readiness)
}

// @Summary  get liveness message
// @Description liveness api
// @Tags health
// @Accept  json
// @Product json
// @Success 200
// @Router /health/liveness [GET]
func liveness(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// @Summary  get readiness message
// @Description readiness api
// @Tags health
// @Accept  json
// @Product json
// @Success 200
// @Router /health/readiness [GET]
func readiness(c *gin.Context) {
	// todo complete self check logic
	c.JSON(http.StatusOK, "status ok")
	return
}
