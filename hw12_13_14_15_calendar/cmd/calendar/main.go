package main

//nolint:depguard
import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/grpc/grpcserver"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/grpc/pb"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/http/handlers"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/server/http"
	"google.golang.org/grpc"
)

const (
	httpTransport = "http"
	grpcTransport = "grpc"
)

var (
	configFile string
	database   string
	transport  string
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
	flag.StringVar(&database, "database", "memory", "What database should we use")
	flag.StringVar(&transport, "transport", "grpc", "What transport we need to use(grpc or http)")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.GetLogger(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage

	switch database {
	case "memory":
		storage = memorystorage.New(logg)
	case "sql":
		storage = sqlstorage.NewPostgresStorage(sqlstorage.PgConfig{
			Host:           config.SQL.Host,
			Username:       config.SQL.Username,
			Password:       config.SQL.Password,
			Port:           config.SQL.Port,
			Database:       config.SQL.Database,
			MigrationsPath: config.SQL.MigrationsPath,
		}, logg, true)
	default:
		logg.Warn("unsupported type of database", map[string]interface{}{"database": database})
	}

	defer storage.Close()
	calendar := app.New(logg, storage)

	switch strings.ToLower(transport) {
	case httpTransport:
		handler := handlers.NewHandler(logg, calendar)
		server := internalhttp.NewServer(logg, config.HTTPServer.Host, config.HTTPServer.Port, handler.InitRoutes())

		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http server: "+err.Error(), nil)
			}
		}()

		logg.Info("calendar is running...", nil)

		if err := server.Start(); err != nil {
			logg.Error("failed to start http server: "+err.Error(), nil)
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	case grpcTransport:
		grpcService := grpcserver.NewServer(calendar, logg)
		server := grpc.NewServer()
		pb.RegisterEventServiceServer(server, grpcService)

		go func() {
			<-ctx.Done()
			logg.Info("stopping grpc server...", nil)
			server.GracefulStop()
		}()

		lsn, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServer.Port))
		if err != nil {
			logg.Fatal(err.Error(), nil)
		}

		logg.Info("starting server on "+lsn.Addr().String(), nil)
		if err := server.Serve(lsn); err != nil {
			logg.Fatal("failed to start grpc server", map[string]interface{}{"error": err})
		}
	default:
		logg.Fatal("unsupported type jf transport", map[string]interface{}{"transport": transport})
	}
}
