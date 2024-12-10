package viacepapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"br.com.cleiton/service-a-cep/internal/interfaces/services"
)

type viacepapi struct {
	urlViaCep string
}

func NewViaCepApi(urlViaCep string) services.CepApiInterface {
	return &viacepapi{
		urlViaCep: urlViaCep,
	}
}

func (v *viacepapi) GetLocation(cep int) (*entities.CEP, error) {
	service := fmt.Sprintf("/ws/%d/json/", cep)
	resp, err := http.Get(v.urlViaCep + service)
	if err != nil {
		return nil, fmt.Errorf("error to get request, %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error to read the response, %w", err)
	}

	var viaCepResponse ViaCepResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		return nil, fmt.Errorf("error to convert json in response")
	}

	if viaCepResponse.CEP == "" {
		return nil, fmt.Errorf("error to find address")
	}

	return &entities.CEP{
		Locality:       viaCepResponse.Localidade,
		Identification: viaCepResponse.CEP,
	}, nil
}
