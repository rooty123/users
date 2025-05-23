package main

import (
	"users/actions"
	"users/db"
	"users/metrics"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rooty123/libs/logger"
)

func main() {
	dbh := db.UsersDBHandler{}
	dbh.RunMigrations()

	e := echo.New()

	// Add Prometheus middleware
	e.Use(metrics.PrometheusMiddleware())

	// Routes
	e.GET("/users", actions.GetUsers)
	e.GET("/users/:id", actions.GetUser)
	e.POST("/users", actions.CreateUser)
	e.PUT("/users/:id", actions.UpdateUser)
	e.DELETE("/users/:id", actions.DeleteUser)

	// Add metrics endpoint
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	logger.WithFields(map[string]any{
		"event": "program_started",
	}).Info("Program started")
	e.Logger.Fatal(e.Start(":8080"))
}
