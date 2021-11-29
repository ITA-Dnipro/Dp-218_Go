package main

import (
	"Dp218Go/pkg/grpcserver"
	"Dp218Go/pkg/httpserver"
	"Dp218Go/pkg/pb"
	"Dp218Go/pkg/postgres"
	repo "Dp218Go/repositories"
	"Dp218Go/routing"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var PG_HOST = os.Getenv("PG_HOST")
var PG_PORT = os.Getenv("PG_PORT")
var POSTGRES_DB = os.Getenv("POSTGRES_DB")
var POSTGRES_USER = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
var HTTP_PORT = os.Getenv("HTTP_PORT")
var MIGRATE_DOWN, _ = strconv.ParseBool(os.Getenv("MIGRATE_DOWN"))

func main() {
	var connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		POSTGRES_USER,
		POSTGRES_PASSWORD,
		PG_HOST,
		PG_PORT,
		POSTGRES_DB)

	pg, err := postgres.NewPostgres(connectionString)
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.CloseDB()

	err = doMigrate(connectionString)
	if err != nil {
		log.Printf("app - Run - Migration issues: %v\n", err)
	}

	var userRepo = repo.New(pg)
	var scooterRepo = repo.NewSc(pg)

	handler := routing.NewRouter()
	routing.AddScooterHandler(handler, scooterRepo)
	routing.AddUserHandler(handler, userRepo)
	httpServer := httpserver.New(handler, httpserver.Port(HTTP_PORT))

	scL,err := scooterRepo.GetAllScooters()
	fmt.Println(scL)

	svr := grpcserver.NewServer()
	svr.Run()
	grpcServer := grpc.NewServer()

	pb.RegisterScooterServiceServer(grpcServer, svr)

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("grpc server started: 8000")
		log.Fatal(grpcServer.Serve(listener))
	}()

	http.HandleFunc("/scrun", svr.ScootRun)

	http.HandleFunc("/", grpcserver.GISHandler)
	http.HandleFunc("/scooter", svr.ScooterHandler)

	fmt.Println("http server started: 8000")
	http.ListenAndServe(":9000", nil)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatalf("app - Run - httpServer.Notify: %v", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Fatalf("app - Run - httpServer.Shutdown: %v", err)
	}
}

func doMigrate(connStr string) error {
	migr, err := migrate.New("file:///home/Dp218Go/migrations", connStr + "?sslmode=disable")
	if err!= nil{
		return err
	}

	migr.Force(20211124)

	if MIGRATE_DOWN {
		migr.Down()
	}

	return migr.Up()
}