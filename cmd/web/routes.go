package main

import (
	"net/http"

	"github.com/aidisapp/musiqcity_v2/internal/config"
	"github.com/aidisapp/musiqcity_v2/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(_ *config.AppConfig) http.Handler {
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

	mux.Post("/create-booking", handlers.Repo.PostCreateBooking)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/signup", handlers.Repo.Signup)
	mux.Post("/user/signup", handlers.Repo.PostSignup)
	mux.Get("/verify-email", handlers.Repo.VerifyUserEmail)
	mux.Get("/user/logout", handlers.Repo.Logout)
	mux.Get("/user/list-service", handlers.Repo.ListService)
	mux.Post("/user/list-service", handlers.Repo.PostListService)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// Use the Auth middleware
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

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
	})

	return mux
}
