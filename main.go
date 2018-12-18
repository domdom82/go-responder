package main

func main() {

	config, err := NewServerConfigFromFile("config.yml")

	if err != nil {
		panic(err)
	}

	for _, serverConfig := range config.ServerConfigs {

		if serverConfig.Http != nil {
			go serverConfig.Http.NewServer().Run()
		}
	}

	c := make(chan bool, 1)

	<-c

}
