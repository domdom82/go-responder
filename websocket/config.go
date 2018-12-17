package websocket

type Config struct {

}

func (cfg *Config) NewServer() *WsServer {

	server := &WsServer{cfg}

	return server
}