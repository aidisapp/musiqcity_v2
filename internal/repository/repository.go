package repository

import (
	"github.com/aidisapp/musiqcity_v2/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	CheckIfUserEmailExist(email string) (bool, error)
	InsertUser(user models.User) (int, error)
	UpdateUserAccessLevel(user models.User) error

	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, int, error)

	AllArtists() ([]models.Artist, error)
	CreateArtist(artist models.Artist) error
	GetArtistByID(id int) (models.Artist, error)
	UpdateArtist(artist models.Artist) error

	AllBookings() ([]models.Bookings, error)
	AllNewBookings() ([]models.Bookings, error)

	AllBookingOptions() ([]models.BookingOptions, error)
	AllArtistBookingOptions(id int) ([]models.BookingOptions, error)
	CreateBookingOption(option models.BookingOptions) error
	GetBookingOptionByID(id int) (models.BookingOptions, error)
	UpdateBookingOption(option models.BookingOptions) error
}
