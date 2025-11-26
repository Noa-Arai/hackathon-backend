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

	// ======== mux + CORS Wrapping ========
	mux := http.NewServeMux()
	handler := middleware.CORS(mux)

	// ======== Usecases & Controllers ========

	// User
	userDAO := dao.NewUserDAO(db)
	registerUserUsecase := usecase.NewRegisterUserUsecase(userDAO)
	registerUserController := controller.NewRegisterUserController(registerUserUsecase)

	loginUsecase := usecase.NewLoginUserUsecase(userDAO, os.Getenv("JWT_SECRET"))
	loginController := controller.NewLoginUserController(loginUsecase)

	// Items
	itemDAO := dao.NewItemDAO(db)
	registerItemUsecase := usecase.NewRegisterItemUsecase(itemDAO)
	registerItemController := controller.NewRegisterItemController(registerItemUsecase)

	getItemsUsecase := usecase.NewGetItemsUsecase(itemDAO)
	getItemsController := controller.NewGetItemsController(getItemsUsecase)

	// Purchase
	purchaseDAO := dao.NewPurchaseDAO(db)
	purchaseUsecase := usecase.NewPurchaseUsecase(purchaseDAO)
	purchaseController := controller.NewPurchaseController(purchaseUsecase)

	// Messages
	messageDAO := dao.NewMessageDAO(db)
	messageUsecase := usecase.NewMessageUsecase(messageDAO)
	messageController := controller.NewMessageController(messageUsecase)

	// AI (Gemini)
	aiUsecase := usecase.NewAIUsecase()
	aiController := controller.NewAIController(aiUsecase)

	// ======== Routing (全部 mux.Handle に変更) ========

	// User register & login
	mux.HandleFunc("/user", registerUserController.Handle)
	mux.HandleFunc("/login", loginController.Handle)

	// Items
	mux.Handle("/items", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				registerItemController.Handle(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	mux.HandleFunc("/items/list", getItemsController.Handle)

	// Purchase
	mux.Handle("/purchase", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				purchaseController.Handle(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	// Messages
	mux.Handle("/messages", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				messageController.Send(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	mux.HandleFunc("/messages/list", messageController.List)

	// AI
	mux.HandleFunc("/ai/describe", aiController.Describe)
	mux.HandleFunc("/ai/ask", aiController.Ask)

	// ========== Graceful shutdown ==========
	closeDBWithSysCall()

	// ========== Start server ==========
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on :" + port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
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
