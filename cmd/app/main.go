package main

import (
	"Dp218Go/configs"
	"Dp218Go/repositories/postgres"
	"Dp218Go/routing"
	"Dp218Go/routing/httpserver"
	"Dp218Go/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
)

var sessionKey = "secretkey"

func main() {
/*
	var connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configs.POSTGRES_USER,
		configs.POSTGRES_PASSWORD,
		configs.PG_HOST,
		configs.PG_PORT,
		configs.POSTGRES_DB)
 */
	var connectionString = "postgres://scooteradmin:Megascooter!@localhost:5444/scooterdb"

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
	var accService = services.NewAccountService(accRepoDb, accRepoDb, accRepoDb)

	var supplierRepoDB = postgres.NewSupplierRepoDB(db)
	var supplierService = services.NewSupplierService(supplierRepoDB, supplierRepoDB, supplierRepoDB)

	sessStore := sessions.NewCookieStore([]byte(sessionKey))
	authService := services.NewAuthService(userRoleRepoDB, sessStore)

	handler := routing.NewRouter(authService)
	routing.AddUserHandler(handler, userService)
	routing.AddAccountHandler(handler, accService)
	routing.AddSupplierHandler(handler, supplierService)
	httpServer := httpserver.New(handler, httpserver.Port(configs.HTTP_PORT))

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
