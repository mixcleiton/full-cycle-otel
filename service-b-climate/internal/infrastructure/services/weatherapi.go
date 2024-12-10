package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"br.com.cleiton/service-b-climate/internal/domain/entities"
	"br.com.cleiton/service-b-climate/internal/interfaces/services"
)

type weatherapi struct {
	url string
	key string
}

func NewWeatherApi(url string, key string) services.ClimaApiInterface {
	return &weatherapi{
		key: key,
		url: url,
	}
}

func (w *weatherapi) GetCurrentClimate(locaty string) (*entities.CurrentClimate, error) {
	service := fmt.Sprintf("/v1/current.json?q=%s&key=%s", locaty, w.key)
	url := w.url + service
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error to execute request, %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error to read the response, %w", err)
	}

	var currentResponse WeatherCurrentResponse
	err = json.Unmarshal(body, &currentResponse)
	if err != nil {
		log.Printf("error to convert json, err %s", err)
		return nil, fmt.Errorf("error to convert json in response")
	}

	return &entities.CurrentClimate{
		Location: currentResponse.Location.Name,
		TempC:    currentResponse.Current.TempC,
		TempF:    currentResponse.Current.TempF,
	}, nil
}
