package app

import (
	"net/http"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers"
	"github.com/JulioZittei/wsrs-ama-go/internal/exception_handler"
	"github.com/JulioZittei/wsrs-ama-go/internal/mappers"
	"github.com/JulioZittei/wsrs-ama-go/internal/middlewares"
	"github.com/JulioZittei/wsrs-ama-go/internal/repositories"
	"github.com/JulioZittei/wsrs-ama-go/internal/services"
	"github.com/JulioZittei/wsrs-ama-go/internal/store/pgstore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

type App struct {
	db      *pgstore.Queries
	handler *chi.Mux
}

func NewApplication(queries *pgstore.Queries) App {
	return App{
		db: queries,
	}
}

func (app *App) Init() {
	// init handler
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.LanguageMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	messageMapper := mappers.MessageMapper{}
	roomMapper := mappers.RoomMapper{}

	// init repositories
	roomsRepository := repositories.NewRoomsRepository(app.db, &roomMapper, &messageMapper)

	// init services
	roomService := services.NewRoomsService(roomsRepository, &roomMapper, &messageMapper)

	// init controllers
	roomsController := controllers.NewRoomsController(roomService, websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	// config routes and handlers
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/subscribe/{room_id}", roomsController.SubscribeRoom)
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", exception_handler.ExceptionHandler(roomsController.CreateRoom))
			r.Get("/", exception_handler.ExceptionHandler(roomsController.GetRooms))

			r.Route("/{room_id}/messages", func(r chi.Router) {
				r.Get("/", exception_handler.ExceptionHandler(roomsController.GetRoomMessages))
				r.Post("/", exception_handler.ExceptionHandler(roomsController.CreateRoomMessage))

				r.Route("/{message_id}", func(r chi.Router) {
					r.Get("/", exception_handler.ExceptionHandler(roomsController.GetRoomMessage))
					r.Patch("/like", exception_handler.ExceptionHandler(roomsController.LikeRoomMessage))
					r.Delete("/like", exception_handler.ExceptionHandler(roomsController.RemoveLikeRoomMessage))
					r.Patch("/answer", exception_handler.ExceptionHandler(roomsController.AnswerRoomMessage))
				})
			})
		})
	})

	app.handler = router
}

func (app *App) GetHandler() http.Handler {
	return app.handler
}
