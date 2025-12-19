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

	// ======== mux ========
	mux := http.NewServeMux()

	// ======== DAO ========
	userDAO := dao.NewUserDAO(db)
	itemDAO := dao.NewItemDAO(db)
	messageDAO := dao.NewMessageDAO(db)
	purchaseDAO := dao.NewPurchaseDAO(db)

	// ======== USECASE ========

	// User
	registerUserUsecase := usecase.NewRegisterUserUsecase(userDAO)
	loginUsecase := usecase.NewLoginUserUsecase(userDAO, os.Getenv("JWT_SECRET"))

	// Profile
	getMyProfileUsecase := usecase.NewGetMyProfileUsecase(userDAO)
	updateProfileUsecase := usecase.NewUpdateProfileUsecase(userDAO)
	updateAvatarUsecase := usecase.NewUpdateAvatarUsecase(userDAO)

	// Items
	registerItemUsecase := usecase.NewRegisterItemUsecase(itemDAO)
	getItemsUsecase := usecase.NewGetItemsUsecase(itemDAO)
	getItemImageUsecase := usecase.NewGetItemImageUsecase(itemDAO)
	updateItemUsecase := usecase.NewUpdateItemUsecase(itemDAO)

	// Messages
	messageUsecase := usecase.NewMessageUsecase(messageDAO)

	// Purchase
	purchaseUsecase := usecase.NewPurchaseUsecase(purchaseDAO)

	// AI
	aiUsecase := usecase.NewAIUsecase()

	// ======== CONTROLLER ========

	// Auth
	registerUserController := controller.NewRegisterUserController(registerUserUsecase)
	loginController := controller.NewLoginUserController(loginUsecase)

	// Profile
	getMyProfileController := controller.NewGetMyProfileController(getMyProfileUsecase)
	updateProfileController := controller.NewUpdateProfileController(updateProfileUsecase)
	updateAvatarController := controller.NewUpdateAvatarController(updateAvatarUsecase)
	getAvatarController := controller.NewGetAvatarController(userDAO)

	// Items
	registerItemController := controller.NewRegisterItemController(registerItemUsecase)
	getItemsController := controller.NewGetItemsController(getItemsUsecase, getItemImageUsecase)
	updateItemController := controller.NewUpdateItemController(updateItemUsecase)

	// Messages
	messageController := controller.NewMessageController(messageUsecase)

	// Purchase
	purchaseController := controller.NewPurchaseController(purchaseUsecase)

	// AI
	aiController := controller.NewAIController(aiUsecase)

	// ======== ROUTING ========

	// ---------------- USER ----------------
	mux.HandleFunc("/user", registerUserController.Handle)
	mux.HandleFunc("/login", loginController.Handle)

	// ---------------- PROFILE (JWT) ----------------
	mux.Handle("/users/me", middleware.AuthMiddleware(
		http.HandlerFunc(getMyProfileController.Handle),
	))
	mux.Handle("/users/me/update", middleware.AuthMiddleware(
		http.HandlerFunc(updateProfileController.Handle),
	))
	mux.Handle("/users/me/avatar", middleware.AuthMiddleware(
		http.HandlerFunc(updateAvatarController.Handle),
	))

	// public avatar fetch
	mux.HandleFunc("/users/avatar", getAvatarController.Handle)

	// ---------------- ITEMS ----------------
	mux.Handle("/items", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				registerItemController.Handle(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	// 商品一覧 / 画像
	mux.HandleFunc("/items/list", getItemsController.Handle)
	mux.HandleFunc("/items/image1", getItemsController.HandleImage(1))
	mux.HandleFunc("/items/image2", getItemsController.HandleImage(2))
	mux.HandleFunc("/items/image3", getItemsController.HandleImage(3))

	// ⭐ 商品編集（JWT必要）
	mux.Handle("/items/update", middleware.AuthMiddleware(
		http.HandlerFunc(updateItemController.Handle),
	))

	//AI感情検索機能追加
	mux.HandleFunc("/items/search/emotion", aiController.HandleEmotionSearch)
	// ---------------- MESSAGES ----------------

	// ⭐ DM（メッセージ送信）
	mux.Handle("/messages", middleware.AuthMiddleware(
		http.HandlerFunc(messageController.Send),
	))

	// ⭐ DM（メッセージ一覧 /rooms?item_id=...&partner_id=...）
	mux.Handle("/messages/room", middleware.AuthMiddleware(
		http.HandlerFunc(messageController.List),
	))

	// ⭐ DMルームリスト
	mux.Handle("/messages/rooms", middleware.AuthMiddleware(
		http.HandlerFunc(messageController.ListRooms),
	))

	// ---------------- PURCHASE ----------------
	mux.Handle("/purchase", middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				purchaseController.Handle(w, r)
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}),
	))

	// ---------------- AI ----------------
	mux.HandleFunc("/ai/describe", aiController.Describe)
	mux.HandleFunc("/ai/ask", aiController.Ask)

	// graceful shutdown
	closeDBWithSysCall()

	// ======== ⭐ CORS を mux.ServeHTTP に適用（決定的な修正） ========
	handler := middleware.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))

	// Start server
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
