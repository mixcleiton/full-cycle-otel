package usecases

import (
	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"br.com.cleiton/service-a-cep/internal/interfaces/services"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/trace"
)

type CurrentClimateInterface interface {
	GetCurrentClimate(ctx context.Context, cep int) (*entities.CurrentClimate, error)
}

var (
	ErrCep     = errors.New("error to get location by cep")
	ErrClimate = errors.New("error to get current climate")
)

type CurrentClimate struct {
	cepApi     services.CepApiInterface
	climateApi services.ClimateInterface
	tracer     trace.Tracer
}

func NewCurrentClimateUsecase(cepApi services.CepApiInterface, climateApi services.ClimateInterface, tracer trace.Tracer) *CurrentClimate {
	return &CurrentClimate{
		cepApi:     cepApi,
		climateApi: climateApi,
		tracer:     tracer,
	}
}

func (c *CurrentClimate) GetCurrentClimate(ctx context.Context, cep int) (*entities.CurrentClimate, error) {
	if cep <= 0 {
		return nil, fmt.Errorf("error to read cep value")
	}

	traceClient := otelhttptrace.NewClientTrace(ctx)
	cepResponse, err := c.cepApi.GetLocation(ctx, cep, traceClient)
	if err != nil || cepResponse == nil {
		log.Printf("err %s", err)
		return nil, ErrCep
	}

	climateResponse, err := c.climateApi.GetClimate(ctx, cepResponse.Locality, traceClient)
	if err != nil {
		log.Errorf("error in get climate, err: %s", err)
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
