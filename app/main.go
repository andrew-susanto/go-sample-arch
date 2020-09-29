package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joez-tkpd/go-sample-arch/infrastructure/infrahttp"
	"github.com/joez-tkpd/go-sample-arch/repository/pgsqlx"
	"github.com/joez-tkpd/go-sample-arch/repository/redispool"
	"github.com/joez-tkpd/go-sample-arch/server/httphandler"
	"github.com/joez-tkpd/go-sample-arch/service/user"
	"github.com/joez-tkpd/go-sample-arch/usecase/account"
)

func main() {
	var port string

	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.StringVar(&port, "port", "8000", "serve port")
	flag.Parse()

	router := infrahttp.NewRouter()

	httpHandler := initHandler()
	httpHandler.Register(router)

	log.Print("serving on:", port)
	log.Fatal(http.ListenAndServe(":"+port, router.GetHandler()))
}

func initHandler() httphandler.Handler {
	sqlxDB := pgsqlx.NewSqlxDB("user=user password=user host=localhost port=5432 sslmode=disable")
	redis := redispool.NewRedisPool("localhost:6379")

	pgsqlxRepo := pgsqlx.NewRepository(sqlxDB)
	redisRepo := redispool.NewResource(redis)

	userRsc := user.NewResource(pgsqlxRepo, redisRepo)
	userSvc := user.NewService(userRsc)
	accountUc := account.NewUsecase(userSvc)

	httpHandler := httphandler.NewHandler(accountUc)
	return httpHandler
}
