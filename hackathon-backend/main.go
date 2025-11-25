package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"hackathon-backend/controller"
	"hackathon-backend/dao"
	"hackathon-backend/usecase"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlDB := os.Getenv("MYSQL_DATABASE")
	connName := os.Getenv("INSTANCE_CONNECTION_NAME")

	if mysqlUser == "" || mysqlPwd == "" || mysqlDB == "" || connName == "" {
		log.Fatal("environment variables not set")
	}

	socketDir := "/cloudsql"

	dsn := fmt.Sprintf("%s:%s@unix(%s/%s)/%s?parseTime=true",
		mysqlUser, mysqlPwd, socketDir, connName, mysqlDB)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db.Ping error: %v", err)
	}

	log.Println("DB Connected")
}

func main() {
	userDAO := dao.NewUserDAO(db)
	registerUsecase := usecase.NewRegisterUserUsecase(userDAO)
	registerController := controller.NewRegisterUserController(registerUsecase)

	loginUsecase := usecase.NewLoginUserUsecase(userDAO, os.Getenv("JWT_SECRET"))
	loginController := controller.NewLoginUserController(loginUsecase)

	http.HandleFunc("/user", registerController.Handle)
	http.HandleFunc("/login", loginController.Handle)

	closeDBWithSysCall()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall: %v", s)
		db.Close()
		os.Exit(0)
	}()
}
