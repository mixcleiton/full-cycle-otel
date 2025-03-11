package main

import (
	"br.com.cleiton/service-a-cep/internal"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	urlViaCep := viper.GetString("URL_VIA_CEP")
	port := viper.GetString("PORT")
	urlClimate := viper.GetString("URL_SERVICE_B_CLIMATE")
	otelServiceName := viper.GetString("OTEL_SERVICE_NAME")
	otelExporterOtlpEndpoint := viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT")

	server := internal.NewServer(urlViaCep, port, urlClimate, otelServiceName, otelExporterOtlpEndpoint)
	server.StartServer()
}
