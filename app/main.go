package main

import (
	// golang package
	"flag"
)

const (
	// Http server constants
	httpServerPortFlag = "port"
	httpDefaultPort    = "8000"

	// Rpc server constants
	rpcServerPortFlag = "rpcport"
	rpcDefaultPort    = "9000"

	// Service mode constants
	serviceModeFlag     = "mode"
	serviceModeHttp     = "http"
	serviceModeConsumer = "consumer"

	// Environment constants
	environmentVariable = "TRAVELOKA_ENV"
	environmentDefault  = "dev"

	serviceName = "cxp-crmapp"
)

func main() {
	var httpServerPort string
	var rpcServerPort string
	var serviceMode string

	// parse command flag
	flag.StringVar(&httpServerPort, httpServerPortFlag, httpDefaultPort, "http server port")
	flag.StringVar(&rpcServerPort, rpcServerPortFlag, rpcDefaultPort, "rpc server port")
	flag.StringVar(&serviceMode, serviceModeFlag, serviceModeHttp, "service mode")
	flag.Parse()

	if serviceMode == serviceModeHttp {
		initHttpApp(httpServerPort, rpcServerPort)
	} else if serviceMode == serviceModeConsumer {
		initConsumerApp()
	}
}
