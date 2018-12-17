package main

import "fmt"

func main() {

	config, err := NewServerConfigFromFile("config.yml")

	if err != nil {
		panic(err)
	}

	for _,serverConfig := range config.ServerConfigs {

		fmt.Println(serverConfig)
	}

	//c := make(chan bool, 1)
	//
	//<-c

}
