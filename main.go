package main

import "fmt"

func main() {

	config, err := NewServerConfigFromFile("config.yml")

	if err != nil {
		panic(err)
	}

	for _, serverConfig := range config.ServerConfigs {
		fmt.Println(serverConfig)
		if serverConfig.Http != nil {
			serverConfig.Http.NewServer().Run()
		}
		if serverConfig.Tcp != nil {
			serverConfig.Tcp.NewServer().Run()
		}
		if serverConfig.Websocket != nil {
			serverConfig.Websocket.NewServer().Run()
		}
	}

	c := make(chan bool, 1)

	<-c

}
