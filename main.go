package main

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/reddaemon/calendargrpcsql/app"
	"github.com/reddaemon/calendargrpcsql/config"
	database "github.com/reddaemon/calendargrpcsql/db"
	"github.com/reddaemon/calendargrpcsql/internal/domain/grpc/server"
	"github.com/reddaemon/calendargrpcsql/internal/domain/grpc/service"
	"github.com/reddaemon/calendargrpcsql/models/storage"

	"net"
	"os"
	"time"

	eventpb "github.com/reddaemon/calendargrpcsql/protofiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		err error
	)

	Logger := log.New("-")

	Config := new(config.Config)

	// иницализация конфигурации
	if len(os.Args) == 2 {
		if err = Config.Load(os.Args[1]); err != nil {
			Logger.Fatal("Can't parse config file: ", err)
		}
	} else {
		Logger.Fatal("Set path to config file")
	}

	Logger.SetHeader(
		"${time_rfc3339_nano} " +
			Config.MainServer.ServiceName +
			" {${short_file}:${line}} ${level} -${message}")

	Config.SetLogger(Logger)
	Logger.Info("Listening on port: ", Config.MainServer.Port)
	lis, err := net.Listen("tcp", Config.MainServer.Host+":"+Config.MainServer.Port)
	if err != nil {
		log.Debugf("failed to listen %v", err)
	}
	db, err := database.GetDb(Config)

	if err != nil {
		log.Fatalf("unable to connect, %v", err)
	}
	grpcServer := grpc.NewServer()
	a := app.App{
		Config: Config,
		Logger: Logger,
		Db:     db,
	}

	contextTimeout := time.Millisecond * 120000
	repo := storage.NewPsqlRepository(&a)
	ucs := service.NewEventUseCase(&a, repo, contextTimeout)
	eventpb.RegisterEventServiceServer(grpcServer, server.NewServer(ucs))
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Error(error.Error(err))
	}

}
