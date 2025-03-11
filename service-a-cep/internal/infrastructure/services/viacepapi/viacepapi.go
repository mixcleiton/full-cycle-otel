package viacepapi

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
	"net/http/httptrace"

	"br.com.cleiton/service-a-cep/internal/domain/entities"
	"br.com.cleiton/service-a-cep/internal/interfaces/services"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type viacepapi struct {
	urlViaCep string
}

func NewViaCepApi(urlViaCep string) services.CepApiInterface {
	return &viacepapi{
		urlViaCep: urlViaCep,
	}
}

func (v *viacepapi) GetLocation(ctx context.Context, cep int, traceClient *httptrace.ClientTrace) (*entities.CEP, error) {
	service := fmt.Sprintf("/ws/%d/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", v.urlViaCep+service, nil)
	if err != nil {
		return nil, fmt.Errorf("error to get request, %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultClient.Transport,
			otelhttp.WithSpanNameFormatter(func(_ string, req *http.Request) string {
				return "get-via-cep"
			}),
		),
	}

	//span := trace.SpanFromContext(ctx)
	//defer span.End()
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
