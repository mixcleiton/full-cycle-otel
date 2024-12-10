package services

import "br.com.cleiton/service-b-climate/internal/domain/entities"

type ClimaApiInterface interface {
	GetCurrentClimate(locaty string) (*entities.CurrentClimate, error)
}
