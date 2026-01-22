package app

import (
	"context"
	"log"
	"msng/internal/config"
	"msng/internal/db"
	"msng/internal/http/handler"
	"msng/internal/http/middleware"
	"msng/internal/jwtutil"
	"msng/internal/repository"
	"msng/internal/service"
	"msng/internal/token"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppMiddleware struct {
}

type AppHandlers struct {
	UserHandler *handler.UserHandler
}

type AppRepositories struct {
	UserRepository  *repository.UserRepository
	TokenRepository *repository.TokenRepository
}

type AppServices struct {
	UserService *service.UserService
	AuthService *service.AuthService
}

func InitAppDatabase(ctx context.Context, pg_conn_dsn string) *pgxpool.Pool {
	db, err := db.InitDatabase(ctx, pg_conn_dsn)
	log.Println("Server initialized")
	if err != nil {
		log.Fatalf("Error while creating connection to DB: %v", err)
	}
	return db
}

func InitAppRepositories(db *pgxpool.Pool) *AppRepositories {
	userRepository := repository.CreateUserRepository(db)
	tokenRepository := repository.NewTokenRepository(db)
	return &AppRepositories{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}

func InitAppServices(jwt *jwtutil.JWTManager, hashToken *token.HashToken, tokenRepository *repository.TokenRepository, userRepository *repository.UserRepository) *AppServices {
	userService := service.CreateUserService(userRepository)
	authService := service.NewAuthService(jwt, userRepository, hashToken, tokenRepository)
	return &AppServices{
		UserService: userService,
		AuthService: authService,
	}
}

func InitAppHandlers(appServices *AppServices) *AppHandlers {
	userHandler := handler.CreateUserHandler(appServices.UserService, appServices.AuthService)
	return &AppHandlers{
		UserHandler: userHandler,
	}
}

func InitPublicRouters(mux *http.ServeMux, appHandlers *AppHandlers) {

	mux.HandleFunc("/auth/register", appHandlers.UserHandler.RegisterUser)
	mux.HandleFunc("/auth/login", appHandlers.UserHandler.LoginUser)
	mux.HandleFunc("/auth/refresh", appHandlers.UserHandler.RefreshUserToken)
	mux.HandleFunc("/auth/logout", appHandlers.UserHandler.LogoutUser)
}

//func InitPrivateRouters(mux *http.ServeMux, appHandlers *AppHandlers) {
//	mux.HandleFunc("/chat", appHandlers.ChatHandler)
//}

func InitAppMiddlewares(mux *http.ServeMux) *http.Handler {
	loggingMux := middleware.LoggingMiddleware(mux)
	//bodyMux := middleware.RequestSizeMiddleware(loggingMux)
	basic := middleware.Basic(loggingMux)

	return &basic
}

func AppInit(cfg *config.Config, ctx context.Context) (*http.Handler, *pgxpool.Pool) {
	db := InitAppDatabase(ctx, cfg.PGConnDSN)

	repos := InitAppRepositories(db)

	jwt := jwtutil.NewJWTManager([]byte(cfg.JWTTokenSecret), "msng", cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	// New hashtoken
	hashToken := token.NewHashToken(cfg.RefreshTokenSecret)

	services := InitAppServices(jwt, hashToken, repos.TokenRepository, repos.UserRepository)

	handlers := InitAppHandlers(services)

	mux := http.NewServeMux()

	InitPublicRouters(mux, handlers)

	handler := InitAppMiddlewares(mux)

	return handler, db
	//New jwt manager

}
