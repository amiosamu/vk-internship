package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gin-contrib/cors"

	_ "github.com/amiosamu/vk-internship/docs"
	"github.com/amiosamu/vk-internship/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(router *gin.Engine, services *service.Services) *gin.Engine {
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf(`{"time":"%s", "method":"%s","uri":"%s", "status":%d,"error":"%s"}`,
				param.TimeStamp.Format(time.RFC3339Nano),
				param.Method,
				param.Path,
				param.StatusCode,
				param.ErrorMessage,
			)
		},
		Output: setLogsFile(),
	}))
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}
	authMiddleware := AuthMiddleware{authService: services.Auth}
	v1 := router.Group("/api/v1", authMiddleware.UserIdentity())
	{
		newAdvertisementRoutes(v1.Group("/advertisements"), services.Advertisement, services.Auth)
	}

	return router
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
