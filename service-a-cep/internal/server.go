package internal

import (
	"br.com.cleiton/service-a-cep/internal/infrastructure/services/otelprovider"
	"context"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"os"
	"os/signal"
	"time"

	"br.com.cleiton/service-a-cep/internal/application/controllers"
	"br.com.cleiton/service-a-cep/internal/application/usecases"
	"br.com.cleiton/service-a-cep/internal/infrastructure/services/climateapi"
	"br.com.cleiton/service-a-cep/internal/infrastructure/services/viacepapi"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Config struct {
	urlViaCep                string
	port                     string
	urlClimate               string
	otelServiceName          string
	otelExporterOtlpEndpoint string
}

func NewServer(urlViaCep, port, urlClimate, otelServiceName, otelExporterOtlpEndpoint string) *Config {
	return &Config{
		urlViaCep:                urlViaCep,
		port:                     port,
		urlClimate:               urlClimate,
		otelServiceName:          otelServiceName,
		otelExporterOtlpEndpoint: otelExporterOtlpEndpoint,
	}
}

func (c *Config) StartServer() {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	viaCepApi := viacepapi.NewViaCepApi(c.urlViaCep)
	climateApi := climateapi.NewClimateApi(c.urlClimate)
	tracer := otel.Tracer("service-a-cep")
	currentClimateUsecase := usecases.NewCurrentClimateUsecase(viaCepApi, climateApi, tracer)
	currentClimateHandler := controllers.NewCurrentClimateHandler(currentClimateUsecase, tracer, c.otelServiceName)
	shutdown, err := otelprovider.InitProvider(c.otelServiceName, c.otelExporterOtlpEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown tracer provider", err)
		}
	}()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))
	e.Use(echoprometheus.NewMiddleware(c.otelServiceName))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/cep/:cep", currentClimateHandler.CurrentClimate)
	port := fmt.Sprintf(":%s", c.port)
	go func() {
		e.Logger.Fatal(e.Start(port))
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, Ctrl+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to another reason...")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
