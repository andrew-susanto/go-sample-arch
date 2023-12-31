package main

import (
	// golang package
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/cacheservice"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/db"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/docdb"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/goenv"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/httpclient"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/opensearchclient"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/paramstore"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/rpcclient"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/secretmanager"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
	cacheRepo "github.com/andrew-susanto/go-sample-arch/repository/cache"
	httpRepo "github.com/andrew-susanto/go-sample-arch/repository/http"
	jsonRpcRepo "github.com/andrew-susanto/go-sample-arch/repository/jsonrpc"
	mongoRepo "github.com/andrew-susanto/go-sample-arch/repository/mongodb"
	opensearchRepo "github.com/andrew-susanto/go-sample-arch/repository/opensearch"
	psqlRepo "github.com/andrew-susanto/go-sample-arch/repository/psql"
	s3Repo "github.com/andrew-susanto/go-sample-arch/repository/s3"
	snsRepo "github.com/andrew-susanto/go-sample-arch/repository/sns"
	sqsRepo "github.com/andrew-susanto/go-sample-arch/repository/sqs"
	"github.com/andrew-susanto/go-sample-arch/server/httphandler"
	userHttpHandler "github.com/andrew-susanto/go-sample-arch/server/httphandler/user"
	"github.com/andrew-susanto/go-sample-arch/server/jsonrpchandler"
	tripJsonRpcHandler "github.com/andrew-susanto/go-sample-arch/server/jsonrpchandler/trip"
	userJsonRpcHandler "github.com/andrew-susanto/go-sample-arch/server/jsonrpchandler/user"
	"github.com/andrew-susanto/go-sample-arch/service/cxpigw"
	"github.com/andrew-susanto/go-sample-arch/service/user"
	"github.com/andrew-susanto/go-sample-arch/usecase/account"
	"github.com/andrew-susanto/go-sample-arch/usecase/trip"

	// external package
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

func initHttpApp(httpServerPort string, rpcServerPort string) {
	// initialize base components
	log.InitLogger()
	goenv.InitEnvironmentVariable()

	environment := os.Getenv(environmentVariable)
	if environment == "" {
		environment = environmentDefault
	}

	tracer := tracer.InitTracer(serviceName, environment)
	monitor := monitor.InitMonitor(environment, serviceName)
	appConfig := config.ParseConfig(environment)
	httpClient := httpclient.InitHttpClient(appConfig)

	// initialize aws
	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(appConfig.AWS.Region),
	)
	if err != nil {
		log.Error(err, nil, "awsConfig.LoadDefaultConfig() got error - main")
	}

	// initialize aws service
	sqsService := awssqs.NewFromConfig(awsConfig)
	s3Service := awss3.NewFromConfig(awsConfig)
	snsService := awssns.NewFromConfig(awsConfig)
	secrets := secretmanager.InitSecretManager(awsConfig)
	paramStore := paramstore.InitParamstore(awsConfig)
	_ = infrastructure.InitInfrastructure(appConfig, paramStore)

	// initialize databases
	cacheClient := cacheservice.InitCache(secrets.Redis)
	dbClient := db.InitDatabaseClient(secrets.Postgresql)
	docDbDatabaseCrmClient := docdb.InitDocDB(secrets.DocDB).Database(secrets.DocDB.DBName)
	openSearchClient := opensearchclient.InitOpenSearchClient(secrets.OpenSearch)

	// initialize rpc client
	cxpIgwRpcClient := rpcclient.NewClient(appConfig.RpcClientConfig.CxpIgwTrip)

	// initialize application stack
	cacheRepository := cacheRepo.NewRepository(cacheClient)
	_ = httpRepo.NewRepository(httpClient)
	jsonRpcRepository := jsonRpcRepo.NewRepository(cxpIgwRpcClient)
	psqlReposiotry := psqlRepo.NewRepository(dbClient)
	_ = sqsRepo.NewRepository(sqsService)
	_ = s3Repo.NewRepository(s3Service)
	_ = mongoRepo.NewRepository(docDbDatabaseCrmClient)
	_ = opensearchRepo.NewRepository(openSearchClient)
	_ = snsRepo.NewRepository(snsService)

	// resource
	userRsc := user.NewResource(&cacheRepository, &psqlReposiotry)
	cxpigwRsc := cxpigw.NewResource(&jsonRpcRepository)

	// service
	userSvc := user.NewService(&userRsc)
	cxpigwSvc := cxpigw.NewService(&cxpigwRsc)

	// usecase
	accountUc := account.NewUsecase(&userSvc)
	tripUc := trip.NewUsecase(&cxpigwSvc)

	// http handler
	userHttpHandler := userHttpHandler.NewHandler(&accountUc)
	httpHandler := httphandler.NewHandler(&userHttpHandler, monitor)
	httpMuxRouter := httpHandler.Register()

	// json rpc handler
	userJsonRpcHandler := userJsonRpcHandler.NewHandler(&accountUc)
	tripJsonRpcHandler := tripJsonRpcHandler.NewHandler(&tripUc)
	jsonRpcHandler := jsonrpchandler.NewHandler(&userJsonRpcHandler, &tripJsonRpcHandler, monitor)
	jsonRpcMuxRouter := jsonRpcHandler.Register()

	var httpSrv *http.Server
	var jsonRpcSrv *http.Server

	// Start JSON Rpc server
	go func() {
		jsonRpcSrv = &http.Server{
			Handler:      jsonRpcMuxRouter,
			Addr:         fmt.Sprintf(":%v", rpcServerPort),
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  30 * time.Second,
		}

		log.Info(nil, fmt.Sprintf("Starting rpc server on :%v", rpcServerPort))

		err := jsonRpcSrv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(err, "JSON RPC server error: %v", err.Error())
		}
	}()

	// Start http server
	go func() {
		httpSrv = &http.Server{
			Handler:      httpMuxRouter,
			Addr:         fmt.Sprintf(":%v", httpServerPort),
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  30 * time.Second,
		}

		log.Info(nil, fmt.Sprintf("Starting http server on %v", httpServerPort))

		err := httpSrv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(err, "HTTP server error: %v", err.Error())
		}
	}()

	var shutdownWaitGroup sync.WaitGroup
	shutdownWaitGroup.Add(1)

	go func() {
		// Wait shutdown signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownRelease()
		log.Info(nil, "Graceful shutdown started.")
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf(err, "HTTP shutdown error: %v", err.Error())
		}

		if err := jsonRpcSrv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf(err, "JSONRpc shutdown error: %v", err.Error())
		}

		// Graceful shutdown db
		db.CloseDBConnection(dbClient)

		// Graceful shutdown tracer
		tracer.Shutdown(shutdownCtx)

		log.Info(nil, "Graceful shutdown completed.")
		shutdownWaitGroup.Done()
	}()

	// Wait graceful shutdown procedure complete
	shutdownWaitGroup.Wait()
}
