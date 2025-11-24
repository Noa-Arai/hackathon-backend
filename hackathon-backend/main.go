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

	// Cloud SQL Proxy 経由で接続
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

	log.Println("DB connected successfully")
}

func main() {
	userDAO := dao.NewUserDAO(db)

	registerUsecase := usecase.NewRegisterUserUsecase(userDAO)

	registerController := controller.NewRegisterUserController(registerUsecase)

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			registerController.Handle(w, r)

		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	closeDBWithSysCall()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Cloud Run のデフォルト
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
		log.Printf("received syscall, %v", s)
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
