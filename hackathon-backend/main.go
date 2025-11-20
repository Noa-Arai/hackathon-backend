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
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	connectionName := os.Getenv("INSTANCE_CONNECTION_NAME") // ★Cloud SQL の接続名

	if mysqlUser == "" || mysqlPwd == "" || connectionName == "" || mysqlDatabase == "" {
		log.Fatal("fail: environment variable not set")
	}

	// Cloud SQL の UNIX ソケット接続
	dsn := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
		mysqlUser, mysqlPwd, connectionName, mysqlDatabase)

	_db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}

	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}

	db = _db
}

func main() {
	userDAO := dao.NewUserDAO(db)

	registerUsecase := usecase.NewRegisterUserUsecase(userDAO)
	searchUsecase := usecase.NewSearchUserUsecase(userDAO)

	registerController := controller.NewRegisterUserController(registerUsecase)
	searchController := controller.NewSearchUserController(searchUsecase)

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			registerController.Handle(w, r)
		case http.MethodGet:
			searchController.Handle(w, r)
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
