package api

type ApiConfig struct {
	Port       int
	ProxyCount int
}

func NewApiConfig(port, proxyCount int) ApiConfig {
	return ApiConfig{
		Port:       port,
		ProxyCount: proxyCount,
	}
}
