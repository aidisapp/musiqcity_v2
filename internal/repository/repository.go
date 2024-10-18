package repository

import (
	"time"

	"github.com/aidisapp/MusiqCity/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)

	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockByID(id int) error

	AllRooms() ([]models.Room, error)
	UpdateRoom(room models.Room) error
	InsertRoom(room models.Room) error
	DeleteRoom(id int) error

	InsertTodoList(todo models.TodoList) error
	GetTodoListByUserID(id int) ([]models.TodoList, error)
	DeleteTodo(id int) error

	GetRestrictionsForCurrentRoom(roomID int, start, end time.Time) ([]models.RoomRestriction, error)

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
