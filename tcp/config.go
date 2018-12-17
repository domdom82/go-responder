package tcp

type Config struct {

}

func (cfg *Config) NewServer() *TcpServer {

	server := &TcpServer{cfg}

	return server
}