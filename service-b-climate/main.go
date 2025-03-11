package main

import (
	"br.com.cleiton/service-b-climate/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	urlClimate := viper.GetString("URL_CLIMATE")
	keyClimate := viper.GetString("KEY_CLIMATE")
	port := viper.GetString("PORT")
	otelServiceName := viper.GetString("OTEL_SERVICE_NAME")

	logrus.WithField("UrlClimate", urlClimate).WithField("key climate", keyClimate).Info("configs")

	server := internal.NewServer(urlClimate, keyClimate, port, otelServiceName)
	server.StartServer()

}
