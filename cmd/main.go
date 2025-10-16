package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teamcutter/subscriptions-service-task/internal/handler"
	"github.com/teamcutter/subscriptions-service-task/internal/repo"
	"github.com/teamcutter/subscriptions-service-task/pkg/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/teamcutter/subscriptions-service-task/docs"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "endpoint", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint"})
)

func metricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		method := c.Request().Method
		endpoint := c.Path()

		// block until next handler (ex. /subscriptions) ends
		err := next(c)
		
		duration := time.Since(start).Seconds()
		status := c.Response().Status

		httpRequestsTotal.WithLabelValues(method, endpoint, http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)

		return err
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// @title Subscription Service API
// @version 1.0
// @description REST API for manage users' subscriptions
// @host localhost:8080
// @BasePath /
func main() {
	logger := setupLogger()

	_ = godotenv.Load()

	time.Sleep(5 * time.Second)

	db, err := database.Connect()
	if err != nil {
		logger.Error("database connection failed", "error", err)
		return
	}
	defer db.Close()

	logger.Info("database connected successfully")

	repo := repo.NewSubscriptionRepo(db)
	h := handler.NewHandler(repo, logger)

	e := echo.New()
	
	subscriptionsGroup := e.Group("/subscriptions")
	subscriptionsGroup.Use(metricsMiddleware)
	subscriptionsGroup.POST("", h.Create)
	subscriptionsGroup.GET("", h.GetAll)
	subscriptionsGroup.DELETE("/:id", h.Delete)
	subscriptionsGroup.GET("/total", h.TotalCost)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			logger.Error("metrics failed", "error", err)
			return
		}
		logger.Info("metrics run on port 2112")
	}()

	port := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
