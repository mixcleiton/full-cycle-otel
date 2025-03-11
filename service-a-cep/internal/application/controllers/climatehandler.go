package controllers

import (
	"errors"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strconv"

	"br.com.cleiton/service-a-cep/internal/application/usecases"
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
	log.Println("iniciando a busca por cep")
	ctx := c.Request().Context()

	var headers = propagation.HeaderCarrier{}
	for key, value := range c.Request().Header {
		headers.Set(key, value[0])
	}

	traceId := uuid.New().String()
	if c.Request().Header.Get("X-Request-ID") == "" {
		headers.Set("X-Request-ID", traceId)
	}

	ctx, rootSpan := h.tracer.Start(ctx, "service-a-cep-span")

	defer rootSpan.End()
	otel.GetTextMapPropagator().Inject(ctx, headers)

	cep, err := strconv.Atoi(c.Param("cep"))
	if err != nil {
		log.Printf("erro ao converter cep, %s", err)
		return c.JSON(http.StatusUnprocessableEntity, "invalid zipcode")
	}

	log.Printf("cep %d", cep)
	currentClimateApiResponse, err := h.currentClimateUsecase.GetCurrentClimate(ctx, cep)
	if err != nil {
		log.Printf("error to consume api rest, err: %s", err)
		if errors.Is(err, usecases.ErrCep) {
			return c.JSON(http.StatusUnprocessableEntity, "invalid zipcode")
		}

		if errors.Is(err, usecases.ErrClimate) {
			return c.JSON(http.StatusNotFound, "can not find zipcode")
		}

		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, CurrentClimateResponse{
		City:       currentClimateApiResponse.Location,
		Fahrenheit: currentClimateApiResponse.TempF,
		Celsius:    currentClimateApiResponse.TempC,
		Kelvin:     currentClimateApiResponse.TempK,
	})
}
