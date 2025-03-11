package services

import (
	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"context"
	"net/http/httptrace"
)

type CepApiInterface interface {
	GetLocation(ctx context.Context, cep int, traceClient *httptrace.ClientTrace) (*entities.CEP, error)
}
