package services

import "br.com.cleiton/service-a-cep/internal/domain/entities"

type CepApiInterface interface {
	GetLocation(cep int) (*entities.CEP, error)
}
