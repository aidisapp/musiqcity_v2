package dbrepo

import (
	"github.com/aidisapp/musiqcity_v2/internal/models"
)

func (repo *testDBRepo) AllUsers() bool {
	return true
}

// Check if a user exists in the database via email
func (m *testDBRepo) CheckIfUserEmailExist(email string) (bool, error) {
	return true, nil
}

// Inserts a user into the database
func (m *testDBRepo) InsertUser(user models.User) (int, error) {
	var newUserID int

	return newUserID, nil
}

// UpdateUserAccessLevel updates a user access level in the database
func (m *testDBRepo) UpdateUserAccessLevel(user models.User) error {
	return nil
}

// GetUserByID returns a user by id
func (repo *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User

	return user, nil
}

// UpdateUser updates a user in the database
func (repo *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

// Authenticate authenticates a user
func (repo *testDBRepo) Authenticate(email, testPassword string) (int, string, int, error) {
	return 1, "", 0, nil
}

// DeleteTodo deletes a todo
func (m *testDBRepo) DeleteTodo(id int) error {
	return nil
}

//  --------Recent---------- //

func (m *testDBRepo) AllArtists() ([]models.Artist, error) {
	var artists []models.Artist
	return artists, nil
}

// Inserts a new Artist into the database
func (repo *testDBRepo) CreateArtist(artist models.Artist) error {
	return nil
}

func (repo *testDBRepo) GetArtistByID(id int) (models.Artist, error) {
	var artist models.Artist

	return artist, nil
}

func (m *testDBRepo) UpdateArtist(artist models.Artist) error {
	return nil
}

// AllBookingss returns a slice of all bookings
func (repo *testDBRepo) AllBookings() ([]models.Bookings, error) {
	var bookings []models.Bookings
	return bookings, nil
}

// AllNewBookings returns a slice of all Bookings
func (m *testDBRepo) AllNewBookings() ([]models.Bookings, error) {
	var bookings []models.Bookings
	return bookings, nil
}

// Get all Booking Options
func (m *testDBRepo) AllBookingOptions() ([]models.BookingOptions, error) {
	var options []models.BookingOptions
	return options, nil
}

// Get all Booking Options
func (m *testDBRepo) AllArtistBookingOptions(id int) ([]models.BookingOptions, error) {
	var options []models.BookingOptions
	return options, nil
}

// Inserts a new Boking Option into the database
func (repo *testDBRepo) CreateBookingOption(option models.BookingOptions) error {
	return nil
}

// Get a booking option by id
func (repo *testDBRepo) GetBookingOptionByID(id int) (models.BookingOptions, error) {
	var option models.BookingOptions
	return option, nil
}

func (m *testDBRepo) UpdateBookingOption(option models.BookingOptions) error {
	return nil
}
