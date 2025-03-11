package controllers

import (
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"

	"br.com.cleiton/service-b-climate/internal/application/usecases"
	"github.com/labstack/echo/v4"
)

type CurrentClimateResponse struct {
	City       string  `json:"city"`
	Fahrenheit float64 `json:"temp_F"`
	Celsius    float64 `json:"temp_C"`
	Kelvin     float64 `json:"temp_K"`
}

type CurrentClimateHandler struct {
	currentClimateUsecase usecases.CurrentClimateInterface
	tracer                trace.Tracer
	requestNameOtel       string
}

func NewCurrentClimateHandler(currentClimateUsecase usecases.CurrentClimateInterface,
	tracer trace.Tracer, requestNameOtel string) *CurrentClimateHandler {
	return &CurrentClimateHandler{
		currentClimateUsecase: currentClimateUsecase,
		tracer:                tracer,
		requestNameOtel:       requestNameOtel,
	}
}

func (h *CurrentClimateHandler) CurrentClimate(c echo.Context) error {
	locality := c.Param("locality")

	ctx := c.Request().Context()
	var headers = propagation.HeaderCarrier{}
	for key, value := range c.Request().Header {
		headers.Set(key, value[0])
	}

	ctx = otel.GetTextMapPropagator().Extract(ctx, headers)
	ctx, span := h.tracer.Start(ctx, h.requestNameOtel)
	defer span.End()

	if len(locality) <= 0 {
		c.JSON(http.StatusNotFound, "can not find locality")
	}

	currentClimateApiResponse, err := h.currentClimateUsecase.GetCurrentClimate(ctx, locality)
	if err != nil {
		log.Printf("error to consume api rest, err: %s", err)
		if errors.Is(err, usecases.ErrClimate) {
			c.JSON(http.StatusNotFound, "can not find locality")
			return nil
		}

		c.JSON(http.StatusInternalServerError, "internal server error")
		return nil
	}

	c.JSON(http.StatusOK, CurrentClimateResponse{
		City:       currentClimateApiResponse.Location,
		Fahrenheit: currentClimateApiResponse.TempF,
		Celsius:    currentClimateApiResponse.TempC,
		Kelvin:     currentClimateApiResponse.TempK,
	})

	return nil
}
