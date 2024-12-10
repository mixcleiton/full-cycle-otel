package main

import (
	"br.com.cleiton/service-b-climate/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	urlClimate := viper.GetString("URL_CLIMATE")
	keyClimate := viper.GetString("KEY_CLIMATE")
	port := viper.GetString("PORT")

	logrus.WithField("UrlClimate", urlClimate).WithField("key climate", keyClimate).Info("configs")

	server := internal.NewServer(urlClimate, keyClimate, port)
	server.StartServer()

}
