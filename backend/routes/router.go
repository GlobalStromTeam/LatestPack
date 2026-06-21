package routes

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"latestpack/config"
	"latestpack/handlers"
	"latestpack/middleware"
	"latestpack/repository"
	"latestpack/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, repos *repository.Repositories) *gin.Engine {
	authSvc := services.NewAuthService(repos.User, cfg.JWTSecret, cfg.JWTTTL)
	versionSvc := services.NewVersionService(repos.Version)
	statsSvc := services.NewStatsService(repos.Stats)
	fileSvc := services.NewFileService(filepath.Join(cfg.DataDir, "files"))

	authHandler := handlers.NewAuthHandler(authSvc)
	dashHandler := handlers.NewDashboardHandler(versionSvc, statsSvc)
	versionHandler := handlers.NewVersionHandler(versionSvc)
	fileHandler := handlers.NewFileHandler(fileSvc)

	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20
	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))

	api := r.Group("/api")
	{
		api.POST("/auth/login", middleware.RateLimitMiddleware(10, time.Minute), authHandler.Login)

		authed := api.Group("")
		authed.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authed.GET("/dashboard/stats", dashHandler.GetStats)
			authed.GET("/dashboard/latest-version", dashHandler.GetLatestVersion)

			authed.GET("/versions", versionHandler.List)
			authed.POST("/versions", versionHandler.Create)
			authed.DELETE("/versions/:version", versionHandler.Delete)

			authed.GET("/files", fileHandler.List)
			authed.POST("/files/folder", fileHandler.CreateFolder)
			authed.POST("/files/upload", fileHandler.Upload)
			authed.PUT("/files/rename", fileHandler.Rename)
			authed.DELETE("/files", fileHandler.Delete)
		}
	}

	// Serve frontend static files
	frontendDist := resolveFrontendDist()
	if frontendDist != "" {
		r.Static("/assets", filepath.Join(frontendDist, "assets"))
		r.StaticFile("/vite.svg", filepath.Join(frontendDist, "vite.svg"))

		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if !strings.HasPrefix(path, "/api") && !strings.Contains(path, ".") {
				c.File(filepath.Join(frontendDist, "index.html"))
				return
			}
			c.Status(http.StatusNotFound)
		})
	}

	return r
}

func resolveFrontendDist() string {
	candidates := []string{
		filepath.Join("..", "frontend", "dist"),
		filepath.Join("frontend", "dist"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	log.Println("Frontend dist not found, API-only mode")
	return ""
}
