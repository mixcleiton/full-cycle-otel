package internal

import (
	"fmt"

	"br.com.cleiton/service-a-cep/internal/application/controllers"
	"br.com.cleiton/service-a-cep/internal/application/usecases"
	"br.com.cleiton/service-a-cep/internal/infrastructure/services/climateapi"
	"br.com.cleiton/service-a-cep/internal/infrastructure/services/viacepapi"
	"github.com/labstack/echo"
)

type Config struct {
	urlViaCep  string
	port       string
	urlClimate string
}

func NewServer(urlViaCep, port, urlClimate string) *Config {
	return &Config{
		urlViaCep:  urlViaCep,
		port:       port,
		urlClimate: urlClimate,
	}
}

func (c *Config) StartServer() {

	viaCepApi := viacepapi.NewViaCepApi(c.urlViaCep)
	climateApi := climateapi.NewClimateApi(c.urlClimate)
	currentClimateUsecase := usecases.NewCurrentClimateUsecase(viaCepApi, climateApi)
	currentClimateHandler := controllers.NewCurrentClimateHandler(currentClimateUsecase)

	e := echo.New()

	e.GET("/cep/:cep", currentClimateHandler.CurrentClimate)

	port := fmt.Sprintf(":%s", c.port)
	e.Logger.Fatal(e.Start(port))
}
