package usecases

import (
	"context"
	"errors"
	"fmt"

	"br.com.cleiton/service-b-climate/internal/domain/entities"
	"br.com.cleiton/service-b-climate/internal/interfaces/services"
	log "github.com/sirupsen/logrus"
)

type CurrentClimateInterface interface {
	GetCurrentClimate(ctx context.Context, locality string) (*entities.CurrentClimate, error)
}

const valueConvertFahrenheit = 273

var (
	ErrClimate = errors.New("error to get current climate")
)

type CurrentClimate struct {
	climateApi services.ClimaApiInterface
}

func NewCurrentClimateUsecase(climateApi services.ClimaApiInterface) CurrentClimateInterface {
	return &CurrentClimate{
		climateApi: climateApi,
	}
}

func (c CurrentClimate) GetCurrentClimate(ctx context.Context, locality string) (*entities.CurrentClimate, error) {
	climateResponse, err := c.climateApi.GetCurrentClimate(ctx, locality)
	if err != nil {
		log.Printf("error in api current climate, err: %s", err)
		return nil, fmt.Errorf("error in get currnet climate %w", err)
	}

	climateResponse.TempK = climateResponse.TempC + valueConvertFahrenheit

	return climateResponse, nil
}
