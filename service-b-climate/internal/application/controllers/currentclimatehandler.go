package controllers

import (
	"errors"
	"log"
	"net/http"

	"br.com.cleiton/service-b-climate/internal/application/usecases"
	"github.com/labstack/echo"
)

type CurrentClimateResponse struct {
	City       string  `json:"city"`
	Fahrenheit float64 `json:"temp_F"`
	Celsius    float64 `json:"temp_C"`
	Kelvin     float64 `json:"temp_K"`
}

type CurrentClimateHandler struct {
	currentClimateUsecase usecases.CurrentClimateInterface
}

func NewCurrentClimateHandler(currentClimateUsecase usecases.CurrentClimateInterface) *CurrentClimateHandler {
	return &CurrentClimateHandler{
		currentClimateUsecase: currentClimateUsecase,
	}
}

func (h *CurrentClimateHandler) CurrentClimate(c echo.Context) error {
	locality := c.Param("locality")

	if len(locality) <= 0 {
		c.JSON(http.StatusNotFound, "can not find locality")
	}

	currentClimateApiResponse, err := h.currentClimateUsecase.GetCurrentClimate(locality)
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
