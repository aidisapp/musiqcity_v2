package models

import (
	"time"
)

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	Token       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Informations for sending mail
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}

// ----------------Recent ---------------

// Artist model
type Artist struct {
	ID            int
	Name          string
	Genres        string
	Description   string
	Phone         string
	Email         string
	City          string
	Facebook      string
	Twitter       string
	Youtube       string
	Logo          string
	Banner        string
	FeaturedImage string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Bookings model
type Bookings struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	Processed int
	ArtistID  int
	Artist    Artist
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Booking options model
type BookingOptions struct {
	ID          int
	Title       string
	Description string
	Price       string
	ArtistID    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
