package climateapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"br.com.cleiton/service-a-cep/internal/interfaces/services"
)

type climateApi struct {
	urlClimate string
}

func NewClimateApi(urlClimate string) services.ClimateInterface {
	return &climateApi{
		urlClimate: urlClimate,
	}
}

func (v *climateApi) GetClimate(locality string) (*entities.CurrentClimate, error) {
	service := fmt.Sprintf("/climate/%s", locality)
	resp, err := http.Get(v.urlClimate + service)
	if err != nil {
		return nil, fmt.Errorf("error to get request, %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error to read the response, %w", err)
	}

	var climateResponse CurrentClimateResponse
	err = json.Unmarshal(body, &climateResponse)
	if err != nil {
		return nil, fmt.Errorf("error to convert json in response")
	}

	if climateResponse.City == "" {
		return nil, fmt.Errorf("error to find address")
	}

	return &entities.CurrentClimate{
		Location: climateResponse.City,
		TempC:    climateResponse.Celsius,
		TempF:    climateResponse.Fahrenheit,
		TempK:    climateResponse.Kelvin,
	}, nil
}
