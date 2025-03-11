package services

import (
	"br.com.cleiton/service-b-climate/internal/domain/entities"
	"context"
)

type ClimaApiInterface interface {
	GetCurrentClimate(ctx context.Context, locaty string) (*entities.CurrentClimate, error)
}
