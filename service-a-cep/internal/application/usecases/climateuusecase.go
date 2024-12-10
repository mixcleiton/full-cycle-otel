package usecases

import (
	"errors"
	"fmt"

	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"br.com.cleiton/service-a-cep/internal/interfaces/services"
	log "github.com/sirupsen/logrus"
)

type CurrentClimateInterface interface {
	GetCurrentClimate(cep int) (*entities.CurrentClimate, error)
}

var (
	ErrCep     = errors.New("error to get location by cep")
	ErrClimate = errors.New("error to get current climate")
)

type CurrentClimate struct {
	cepApi     services.CepApiInterface
	climateApi services.ClimateInterface
}

func NewCurrentClimateUsecase(cepApi services.CepApiInterface, climateApi services.ClimateInterface) *CurrentClimate {
	return &CurrentClimate{
		cepApi:     cepApi,
		climateApi: climateApi,
	}
}

func (c CurrentClimate) GetCurrentClimate(cep int) (*entities.CurrentClimate, error) {
	if cep <= 0 {
		return nil, fmt.Errorf("error to read cep value")
	}

	cepResponse, err := c.cepApi.GetLocation(cep)
	if err != nil || cepResponse == nil {
		log.Printf("err %s", err)
		return nil, ErrCep
	}

	climateResponse, err := c.climateApi.GetClimate(cepResponse.Locality)
	if err != nil {
		return nil, ErrClimate
	}

	log.WithField("cepResponse", climateResponse).Info("cep api result")
	return &entities.CurrentClimate{
		Location: cepResponse.Locality,
		TempC:    climateResponse.TempC,
		TempF:    climateResponse.TempF,
		TempK:    climateResponse.TempK,
	}, nil
}
