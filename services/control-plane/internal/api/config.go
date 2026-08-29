package api

import "os"

type Config struct {
	Addr           string
	RequestIDHeader string
}

func LoadConfig() Config {
	addr := os.Getenv("FTN_CONTROL_PLANE_ADDR")
	if addr == "" { addr = ":8080" }
	header := os.Getenv("FTN_REQUEST_ID_HEADER")
	if header == "" { header = "X-Request-ID" }
	return Config{Addr: addr, RequestIDHeader: header}
}
