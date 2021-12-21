package main

import (
	"Dp218Go/configs"
	"Dp218Go/protos"
	"Dp218Go/repositories/postgres"
	"Dp218Go/routing"
	"Dp218Go/routing/grpcserver"
	"Dp218Go/routing/httpserver"
	"Dp218Go/services"
	"Dp218Go/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
)

func main() {

	var connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configs.POSTGRES_USER,
		configs.POSTGRES_PASSWORD,
		configs.PG_HOST,
		configs.PG_PORT,
		configs.POSTGRES_DB)

	db, err := postgres.NewConnection(connectionString)
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer db.CloseDB()

	err = doMigrate(connectionString)
	if err != nil {
		log.Printf("app - Run - Migration issues: %v\n", err)
	}

	var userRoleRepoDB = postgres.NewUserRepoDB(db)
	var userService = services.NewUserService(userRoleRepoDB, userRoleRepoDB)

	var accRepoDb = postgres.NewAccountRepoDB(userRoleRepoDB, db)
	var clock = services.NewClock()
	var accService = services.NewAccountService(accRepoDb, accRepoDb, accRepoDb, clock)
	var stationRepoDb = postgres.NewStationRepoDB(db)
	var stationService = services.NewStationService(stationRepoDb)

	var scooterRepo = postgres.NewScooterRepoDB(db)
	var grpcScooterService = services.NewGrpcScooterService(scooterRepo)
	var scooterService = services.NewScooterService(scooterRepo)

	var supplierRepoDB = postgres.NewSupplierRepoDB(db)
	var supplierService = services.NewSupplierService(supplierRepoDB)

	var problemRepoDb = postgres.NewProblemRepoDB(userRoleRepoDB, scooterRepo, db)
	var problemService = services.NewProblemService(problemRepoDb)
	var solutionRepoDb = postgres.NewSolutionRepoDB(db)
	var solutionService = services.NewSolutionService(solutionRepoDb)

	var orderRepoDB = postgres.NewOrderRepoDB(db)
	var orderService = services.NewOrderService(orderRepoDB)

	sessStore := sessions.NewCookieStore([]byte(configs.SESSION_SECRET))
	authService := services.NewAuthService(userRoleRepoDB, sessStore)

	custService := services.NewCustomerService(stationRepoDb)

	handler := routing.NewRouter()
	routing.AddAuthHandler(handler, authService)
	routing.AddCustomerHandler(handler, custService)
	routing.AddUserHandler(handler, userService)
	routing.AddStationHandler(handler, stationService)
	routing.AddAccountHandler(handler, accService)
	routing.AddScooterHandler(handler, scooterService)
	routing.AddProblemHandler(handler, problemService, solutionService)
	routing.AddGrpcScooterHandler(handler, grpcScooterService)
	routing.AddOrderHandler(handler, orderService)
	routing.AddSupplierHandler(handler, supplierService)
	httpServer := httpserver.New(handler, httpserver.Port(configs.HTTP_PORT))
	handler.HandleFunc("/scooter", httpServer.ScooterHandler)

	grpcServer := grpcserver.NewGrpcServer()
	protos.RegisterScooterServiceServer(grpcServer, httpServer)
	http.ListenAndServe(":8080", handler)

	utils.CheckKafka() //TODO: delete after checking

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
	migr, err := migrate.New("file://"+configs.MIGRATIONS_PATH, connStr+"?sslmode=disable")
	if err != nil {
		return err
	}

	if configs.MIGRATE_VERSION_FORCE > 0 {
		migr.Force(configs.MIGRATE_VERSION_FORCE)
	}

	if configs.MIGRATE_DOWN {
		migr.Down()
	}

	return migr.Up()
}
