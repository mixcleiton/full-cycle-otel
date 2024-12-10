package internal

import (
	"fmt"

	"br.com.cleiton/service-b-climate/internal/application/controllers"
	"br.com.cleiton/service-b-climate/internal/application/usecases"
	weatherapi "br.com.cleiton/service-b-climate/internal/infrastructure/services"
	"github.com/labstack/echo"
)

type Config struct {
	urlClimate string
	keyClimate string
	port       string
}

func NewServer(urlClimate, keyClimate, port string) *Config {
	return &Config{
		urlClimate: urlClimate,
		keyClimate: keyClimate,
		port:       port,
	}
}

func (c *Config) StartServer() {

	climateApi := weatherapi.NewWeatherApi(c.urlClimate, c.keyClimate)
	currentClimateUsecase := usecases.NewCurrentClimateUsecase(climateApi)
	currentClimateHandler := controllers.NewCurrentClimateHandler(currentClimateUsecase)

	e := echo.New()

	e.GET("/climate/:locality", currentClimateHandler.CurrentClimate)

	port := fmt.Sprintf(":%s", c.port)
	e.Logger.Fatal(e.Start(port))
}
