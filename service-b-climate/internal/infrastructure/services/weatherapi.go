package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

func (w *weatherapi) GetCurrentClimate(ctx context.Context, locaty string) (*entities.CurrentClimate, error) {
	service := fmt.Sprintf("/v1/current.json?q=%s&key=%s", locaty, w.key)
	url := w.url + service

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error to get request, %w", err)
	}

	var headers = propagation.HeaderCarrier{}
	for _, value := range req.Header {
		headers.Set(value[0], value[1])
	}
	otel.GetTextMapPropagator().Inject(ctx, headers)

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			log.Println("error closing response body")
		}
	}(resp.Body)

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
