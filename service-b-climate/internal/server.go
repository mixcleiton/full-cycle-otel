package internal

import (
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"time"

	"br.com.cleiton/service-b-climate/internal/application/controllers"
	"br.com.cleiton/service-b-climate/internal/application/usecases"
	weatherapi "br.com.cleiton/service-b-climate/internal/infrastructure/services"
	"github.com/labstack/echo/v4"
)

type Config struct {
	urlClimate      string
	keyClimate      string
	port            string
	otelServiceName string
}

func NewServer(urlClimate, keyClimate, port, otelServiceName string) *Config {
	return &Config{
		urlClimate:      urlClimate,
		keyClimate:      keyClimate,
		port:            port,
		otelServiceName: otelServiceName,
	}
}

func (c *Config) StartServer() {

	climateApi := weatherapi.NewWeatherApi(c.urlClimate, c.keyClimate)
	currentClimateUsecase := usecases.NewCurrentClimateUsecase(climateApi)
	tracer := otel.Tracer(c.otelServiceName)
	currentClimateHandler := controllers.NewCurrentClimateHandler(currentClimateUsecase, tracer, c.otelServiceName)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))
	e.Use(echoprometheus.NewMiddleware(c.otelServiceName))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/climate/:locality", currentClimateHandler.CurrentClimate)

	port := fmt.Sprintf(":%s", c.port)
	e.Logger.Fatal(e.Start(port))
}
