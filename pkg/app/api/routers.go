package api

import (
	"fmt"
	"net/http"
	"net/url"
	"practice/pkg/logger"
	"strings"

	"practice/pkg/conf"

	docs "practice/docs/app_swagger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	g            = gin.Default()
	HdrRequestID = "X-Request-Id"
)

func New() *gin.Engine {
	InitApi()
	return g
}

func InitApi() {
	g.Use(AddRequestID)
	g.Use(Cors())

	if conf.GetConfig().EnableSwagger {
		// load swagger api
		registerSwaggerRouter(g)
	}

	// add healthCheck Api
	registerHealthCheckApi(g)

	// add tool api
	registerToolV1Api(g)

}

func registerSwaggerRouter(engine *gin.Engine) {
	// programatically set swagger infos
	c := conf.GetConfig()
	log := logger.Logger("api")
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", c.Host, c.Port))
	if err != nil {
		log.Error(err, fmt.Sprintf("Parse app issuer error: %v", err.Error()))
	}

	docs.SwaggerInfo.Schemes = []string{u.Scheme}
	// docs.SwaggerInfo.Host = u.Host
	docs.SwaggerInfo.BasePath = u.Path
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func AddRequestID(c *gin.Context) {
	if c.GetHeader(HdrRequestID) == "" {
		c.Request.Header.Set(HdrRequestID, uuid.Version(4).String())
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			// This is to allow access to all domains
			c.Header("Access-Control-Allow-Origin", "*")
			// All cross-domain request methods supported by the server, in order to avoid multiple 'preflight' requests for browsing requests
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE,UPDATE")
			//  header's type
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//Allow cross domain
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			//Cache request information in seconds
			c.Header("Access-Control-Max-Age", "3600")
			//Whether cross-domain requests require cookie information. The default setting is true
			c.Header("Access-Control-Allow-Credentials", "false")
			//Set the return format to json
			c.Set("content-type", "application/json")
		}

		//Release all OPTIONS methods
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// handle the request
		c.Next()
	}
}
