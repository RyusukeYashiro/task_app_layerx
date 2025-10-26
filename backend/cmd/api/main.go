package main

import (
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"

	"github.com/ryusuke/task_app_layerx/internal/infrastructure/clock"
	"github.com/ryusuke/task_app_layerx/internal/infrastructure/mysql"
	"github.com/ryusuke/task_app_layerx/internal/infrastructure/mysql/repository"
	"github.com/ryusuke/task_app_layerx/internal/presentation/handler"
	"github.com/ryusuke/task_app_layerx/internal/presentation/middleware"
	authuc "github.com/ryusuke/task_app_layerx/internal/usecase/auth"
	taskuc "github.com/ryusuke/task_app_layerx/internal/usecase/task"
	"github.com/ryusuke/task_app_layerx/pkg/auth"
	"github.com/ryusuke/task_app_layerx/pkg/hash"
)

func main() {
	// 環境変数の読み込み
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "task_app_layerx"
	}

	// DB接続
	db, err := mysql.NewDBFromDSN(dbDSN)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	// Infrastructure層の初期化
	txManager := mysql.NewTxManager(db)
	executor := txManager.AsExecutor()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskAssigneeRepo := repository.NewTaskAssigneeRepository()

	// pkg層の初期化
	realClock := clock.New()
	jwtService := auth.NewJWTService(jwtSecret, jwtIssuer, 24*time.Hour, func() time.Time {
		return realClock.Now()
	})
	bcryptService := hash.NewBcryptService(12)

	// UseCase層の初期化
	authUseCase := authuc.NewAuthUseCase(
		userRepo,
		txManager,
		realClock,
		jwtService,
		bcryptService,
	)

	taskUseCase := taskuc.NewTaskUseCase(
		taskRepo,
		taskAssigneeRepo,
		userRepo,
		txManager,
		realClock,
	)

	// Handler層の初期化
	authHandler := handler.NewAuthHandler(authUseCase)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	// Echoの設定
	e := echo.New()

	// ミドルウェアの設定
	e.Use(echoMw.Recover())
	e.Use(echoMw.Logger())
	e.Use(echoMw.CORS())

	// ルーティングの設定
	api := e.Group("/api/v1")

	// 認証不要なエンドポイント
	auth := api.Group("/auth")
	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)

	// 認証が必要なエンドポイント
	jwtMiddleware := middleware.JWTMiddleware(jwtService, userRepo, executor)

	authProtected := api.Group("/auth")
	authProtected.Use(jwtMiddleware)
	authProtected.POST("/logout", authHandler.Logout)

	users := api.Group("/users")
	users.Use(jwtMiddleware)
	users.GET("", authHandler.GetUsers)

	tasks := api.Group("/tasks")
	tasks.Use(jwtMiddleware)
	tasks.GET("", taskHandler.ListTasks)
	tasks.POST("", taskHandler.CreateTask)
	tasks.GET("/:id", taskHandler.GetTask)
	tasks.PATCH("/:id", taskHandler.UpdateTask)
	tasks.DELETE("/:id", taskHandler.DeleteTask)

	// サーバー起動
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
