package services

import "br.com.cleiton/service-a-cep/internal/domain/entities"

type ClimateInterface interface {
	GetClimate(locality string) (*entities.CurrentClimate, error)
}
