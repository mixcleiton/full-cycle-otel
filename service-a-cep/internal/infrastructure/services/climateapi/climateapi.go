package climateapi

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
	"net/http/httptrace"

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

func (v *climateApi) GetClimate(ctx context.Context, locality string, trace *httptrace.ClientTrace) (*entities.CurrentClimate, error) {
	service := fmt.Sprintf("/climate/%s", locality)
	req, err := http.NewRequestWithContext(ctx, "GET", v.urlClimate+service, nil)
	if err != nil {
		return nil, fmt.Errorf("error to get request, %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultClient.Transport,
			otelhttp.WithSpanNameFormatter(func(_ string, req *http.Request) string {
				return "get-service-b"
			}),
		),
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error to get response, %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error("error closing response body")
		}
	}(resp.Body)

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
