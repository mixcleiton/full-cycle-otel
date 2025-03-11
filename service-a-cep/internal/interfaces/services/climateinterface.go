package services

import (
	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"context"
	"net/http/httptrace"
)

type ClimateInterface interface {
	GetClimate(ctx context.Context, locality string, trace *httptrace.ClientTrace) (*entities.CurrentClimate, error)
}
