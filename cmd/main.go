package main

import (
	"api/configs"
	"api/internal/auth"
	"api/internal/link"
	"api/internal/stat"
	"api/internal/user"
	"api/pkg/db"
	"api/pkg/event"
	"api/pkg/middleware"
	"fmt"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus: eventBus,
		StatRepo: statRepository,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		Config:         conf,
		LinkRepository: linkRepository,
		EventBus:       eventBus,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick()

	// Middlewares
	middlewares := middleware.Chain(middleware.CORS, middleware.Logging)
	return middlewares(router)
}

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
