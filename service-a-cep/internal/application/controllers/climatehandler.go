package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"br.com.cleiton/service-a-cep/internal/application/usecases"
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
	log.Println("iniciando a busca por cep")
	cep, err := strconv.Atoi(c.Param("cep"))
	if err != nil {
		log.Printf("erro ao converter cep, %s", err)
		c.JSON(http.StatusUnprocessableEntity, "invalid zipcode")
	}

	log.Printf("cep %d", cep)
	currentClimateApiResponse, err := h.currentClimateUsecase.GetCurrentClimate(cep)
	if err != nil {
		log.Printf("error to consume api rest, err: %s", err)
		if errors.Is(err, usecases.ErrCep) {
			c.JSON(http.StatusUnprocessableEntity, "invalid zipcode")
			return nil
		}

		if errors.Is(err, usecases.ErrClimate) {
			c.JSON(http.StatusNotFound, "can not find zipcode")
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
