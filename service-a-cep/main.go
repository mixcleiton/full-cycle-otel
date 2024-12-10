package main

import (
	"br.com.cleiton/service-a-cep/internal"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	urlViaCep := viper.GetString("URL_VIA_CEP")
	port := viper.GetString("PORT")
	urlClimate := viper.GetString("URL_SERVICE_B_CLIMATE")

	server := internal.NewServer(urlViaCep, port, urlClimate)
	server.StartServer()
}
