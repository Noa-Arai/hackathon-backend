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
	"hackathon-backend/middleware"
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

	// --- User: Register ---
	userDAO := dao.NewUserDAO(db)
	registerUserUsecase := usecase.NewRegisterUserUsecase(userDAO)
	registerUserController := controller.NewRegisterUserController(registerUserUsecase)

	// --- User: Login ---
	loginUsecase := usecase.NewLoginUserUsecase(userDAO, os.Getenv("JWT_SECRET"))
	loginController := controller.NewLoginUserController(loginUsecase)

	// --- Items: Register ---
	itemDAO := dao.NewItemDAO(db)
	registerItemUsecase := usecase.NewRegisterItemUsecase(itemDAO)
	registerItemController := controller.NewRegisterItemController(registerItemUsecase)

	// --- Items: Get ---
	getItemsUsecase := usecase.NewGetItemsUsecase(itemDAO)
	getItemsController := controller.NewGetItemsController(getItemsUsecase)

	// --- Purchase: Insert ---
	purchaseDAO := dao.NewPurchaseDAO(db)
	purchaseUsecase := usecase.NewPurchaseUsecase(purchaseDAO)
	purchaseController := controller.NewPurchaseController(purchaseUsecase)

	// --------- Routing -----------

	// Item 登録（POSTのみ Auth required）
	http.Handle("/items", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				registerItemController.Handle(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	// Item 一覧（GET）
	http.HandleFunc("/items/list", getItemsController.Handle)

	// ユーザー登録
	http.HandleFunc("/user", registerUserController.Handle)

	// ログイン
	http.HandleFunc("/login", loginController.Handle)

	// --- Purchase API（POST） Auth 必須 ---
	http.Handle("/purchase", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				purchaseController.Handle(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	// Graceful shutdown
	closeDBWithSysCall()

	// start server
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
