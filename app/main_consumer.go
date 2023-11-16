package main

import (
	// golang package
	"context"
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
	"github.com/andrew-susanto/go-sample-arch/server/cronhandler"
	userCronHandler "github.com/andrew-susanto/go-sample-arch/server/cronhandler/user"
	"github.com/andrew-susanto/go-sample-arch/server/sqshandler"
	userSqsHandler "github.com/andrew-susanto/go-sample-arch/server/sqshandler/user"
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

func initConsumerApp() {
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
	infra := infrastructure.InitInfrastructure(appConfig, paramStore)

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
	_ = trip.NewUsecase(&cxpigwSvc)

	// sqs handler
	userSqsHandler := userSqsHandler.NewHandler(&accountUc)
	sqsHandler := sqshandler.NewHandler(&userSqsHandler)

	// cron handler
	userCronHandler := userCronHandler.NewHandler(&accountUc)
	cronHandler := cronhandler.NewHandler(&userCronHandler)

	var handlerWaitGroup sync.WaitGroup

	// Start SQS client handler
	ctxSqsHandler, shutdownSqsHandler := context.WithCancel(context.Background())
	sqsHandler.Register(ctxSqsHandler, &handlerWaitGroup, infra, monitor, sqsService)

	// Start Cron Handler
	ctxCronHandler, shutdownCronHandler := context.WithCancel(context.Background())
	cronHandler.Register(ctxCronHandler, &handlerWaitGroup, monitor)

	var shutdownWaitGroup sync.WaitGroup
	shutdownWaitGroup.Add(1)

	go func() {
		// Wait shutdown signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCronHandler()
		shutdownSqsHandler()

		handlerWaitGroup.Wait()

		// Graceful shutdown db
		db.CloseDBConnection(dbClient)

		// Graceful shutdown tracer
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownRelease()

		tracer.Shutdown(shutdownCtx)

		log.Info(nil, "Graceful shutdown completed.")
		shutdownWaitGroup.Done()
	}()

	// Wait graceful shutdown procedure complete
	shutdownWaitGroup.Wait()
}
