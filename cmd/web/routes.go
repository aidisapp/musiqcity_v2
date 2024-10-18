package main

import (
	"net/http"

	"github.com/aidisapp/MusiqCity/internal/config"
	"github.com/aidisapp/MusiqCity/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Add all our middlewares here
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/artists", handlers.Repo.ArtistsPage)
	mux.Get("/artists/{id}", handlers.Repo.SingleArtist)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Post("/reservation-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// Use the Auth middleware
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)

		mux.Get("/rooms", handlers.Repo.AdminAllRooms)
		mux.Get("/rooms/{id}", handlers.Repo.AdminSingleRoom)
		mux.Post("/rooms/{id}", handlers.Repo.PostAdminSingleRoom)
		mux.Get("/rooms/new-room", handlers.Repo.AdminNewRoom)
		mux.Post("/rooms/new-room", handlers.Repo.PostAdminNewRoom)
		mux.Get("/delete-room/{id}", handlers.Repo.AdminDeleteRoom)

		mux.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		mux.Get("/todo-list", handlers.Repo.AdminTodoList)
		mux.Post("/todo-list", handlers.Repo.PostAdminTodoList)
		mux.Get("/delete-todo/{id}", handlers.Repo.AdminDeleteTodo)

		// Recent ----------
		mux.Get("/artists", handlers.Repo.AdminAllArtists)
		mux.Get("/artists/new-artist", handlers.Repo.AdminNewArtist)
		mux.Post("/artists/new-artist", handlers.Repo.PostAdminNewArtist)
		mux.Get("/artists/{id}", handlers.Repo.AdminSingleArtist)
		mux.Post("/artists/{id}", handlers.Repo.PostAdminSingleArtist)

		mux.Get("/all-bookings", handlers.Repo.AdminAllBookings)
		mux.Get("/new-bookings", handlers.Repo.AdminNewBookings)

		mux.Get("/booking-options", handlers.Repo.AdminAllOptions)
		mux.Get("/booking-options/new-option", handlers.Repo.AdminNewOption)
		mux.Post("/booking-options/new-option", handlers.Repo.PostAdminNewOption)
		mux.Get("/booking-options/{id}", handlers.Repo.AdminSingleOption)
		mux.Post("/booking-options/{id}", handlers.Repo.PostAdminSingleOption)

		mux.Post("/create-booking", handlers.Repo.PostMakeReservation)
	})

	return mux
}
