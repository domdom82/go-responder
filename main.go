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
			go serverConfig.Http.NewServer().Run()
		}
		if serverConfig.Tcp != nil {
			go serverConfig.Tcp.NewServer().Run()
		}

	}

	c := make(chan bool, 1)

	<-c

}
