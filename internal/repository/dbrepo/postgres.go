package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/aidisapp/musiqcity_v2/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate authenticates a user
func (repo *postgresDBRepo) Authenticate(email, testPassword string) (int, string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	var access_level int

	row := repo.DB.QueryRowContext(ctx, "SELECT id, password, access_level FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword, &access_level)
	if err != nil {
		return id, "", 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, access_level, nil
}

// Check if a user exists in the database via email
func (m *postgresDBRepo) CheckIfUserEmailExist(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select
			count(id)
		from
			users
		where
			email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return false, nil
	}

	return true, nil
}

// Inserts a user into the database
func (m *postgresDBRepo) InsertUser(user models.User) (int, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newUserID int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return 0, err
	}

	query := `insert into users (first_name, last_name, email, password, access_level, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = m.DB.QueryRowContext(context, query, user.FirstName, user.LastName, user.Email, hashedPassword, user.AccessLevel, time.Now(), time.Now()).Scan(&newUserID)

	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// UpdateUserAccessLevel updates a user access level in the database
func (m *postgresDBRepo) UpdateUserAccessLevel(user models.User) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update users set access_level = $1, updated_at = $2 where id = $3
	`

	_, err := m.DB.ExecContext(context, query, user.AccessLevel, time.Now(), user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (repo *postgresDBRepo) InsertReservation(res models.Bookings) (int, error) {
	// Close this transaction if unable to run this statement within 3 seconds
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	insertStatement := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := repo.DB.QueryRowContext(context, insertStatement, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.ArtistID, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (repo *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select
			count(id)
		from
			room_restrictions
		where
			room_id = $1
			and $2 <= end_date and $3 >= start_date;`

	row := repo.DB.QueryRowContext(context, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// GetUserByID returns a user by id
func (repo *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			from users where id = $1`

	row := repo.DB.QueryRowContext(context, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (repo *postgresDBRepo) UpdateUser(user models.User) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
	`

	_, err := repo.DB.ExecContext(context, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "update reservations set processed = $1 where id = $2"

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

//  --------Recent---------- //

// Get all artists
func (m *postgresDBRepo) AllArtists() ([]models.Artist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var artists []models.Artist

	query := `select id, name, genres, description, phone, email, city, facebook, twitter, youtube, logo, banner, featured_image, created_at, updated_at from artists order by created_at asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return artists, err
	}
	defer rows.Close()

	for rows.Next() {
		var artist models.Artist
		err := rows.Scan(
			&artist.ID,
			&artist.Name,
			&artist.Genres,
			&artist.Description,
			&artist.Phone,
			&artist.Email,
			&artist.City,
			&artist.Facebook,
			&artist.Twitter,
			&artist.Youtube,
			&artist.Logo,
			&artist.Banner,
			&artist.FeaturedImage,
			&artist.CreatedAt,
			&artist.UpdatedAt,
		)
		if err != nil {
			return artists, err
		}
		artists = append(artists, artist)
	}

	if err = rows.Err(); err != nil {
		return artists, err
	}

	return artists, nil
}

// Inserts a new Artist into the database
func (repo *postgresDBRepo) CreateArtist(artist models.Artist) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into artists (name, genres, description, phone, email, city, facebook, twitter, youtube, logo, banner, featured_image, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := repo.DB.ExecContext(context, query, artist.Name, artist.Genres, artist.Description, artist.Phone, artist.Email, artist.City, artist.Facebook, artist.Twitter, artist.Youtube, artist.Logo, artist.Banner, artist.FeaturedImage, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

// Get an artist by id
func (repo *postgresDBRepo) GetArtistByID(id int) (models.Artist, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var artist models.Artist

	query := `
		select id, name, genres, description, phone, email, city, facebook, twitter, youtube, logo, banner, featured_image, created_at, updated_at from artists where id = $1
	`

	row := repo.DB.QueryRowContext(context, query, id)
	err := row.Scan(
		&artist.ID,
		&artist.Name,
		&artist.Genres,
		&artist.Description,
		&artist.Phone,
		&artist.Email,
		&artist.City,
		&artist.Facebook,
		&artist.Twitter,
		&artist.Youtube,
		&artist.Logo,
		&artist.Banner,
		&artist.FeaturedImage,
		&artist.CreatedAt,
		&artist.UpdatedAt,
	)

	if err != nil {
		return artist, err
	}

	return artist, nil
}

// UpdateArtist updates an artist in the database
func (m *postgresDBRepo) UpdateArtist(artist models.Artist) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update artists set name = $1, genres = $2, description = $3, phone = $4, email = $5, city = $6, facebook = $7, twitter = $8, youtube = $9, logo = $10, banner = $11, featured_image = $12, updated_at = $13
		where id = $14		
	`

	_, err := m.DB.ExecContext(ctx, query,
		artist.Name,
		artist.Genres,
		artist.Description,
		artist.Phone,
		artist.Email,
		artist.City,
		artist.Facebook,
		artist.Twitter,
		artist.Youtube,
		artist.Logo,
		artist.Banner,
		artist.FeaturedImage,
		time.Now(),
		artist.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// AllBookingss returns a slice of all bookings
func (repo *postgresDBRepo) AllBookings() ([]models.Bookings, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bookings []models.Bookings

	query := `
		select b.id, b.first_name, b.last_name, b.email, b.phone, b.start_date, 
		b.end_date, b.processed, b.artist_id, b.created_at, b.updated_at,
		ar.id, ar.name, ar.genres, ar.description, ar.city
		from bookings b
		left join artists ar on (b.artist_id = ar.id)
		order by b.start_date asc
	`

	rows, err := repo.DB.QueryContext(context, query)
	if err != nil {
		return bookings, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Bookings
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.Processed,
			&i.ArtistID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Artist.ID,
			&i.Artist.Name,
			&i.Artist.Genres,
			&i.Artist.Description,
			&i.Artist.City,
		)

		if err != nil {
			return bookings, err
		}
		bookings = append(bookings, i)
	}

	if err = rows.Err(); err != nil {
		return bookings, err
	}

	return bookings, nil
}

// AllNewBookings returns a slice of all Bookings
func (m *postgresDBRepo) AllNewBookings() ([]models.Bookings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bookings []models.Bookings

	query := `
		select b.id, b.first_name, b.last_name, b.email, b.phone, b.start_date, 
		b.end_date, b.processed, b.artist_id, b.created_at, b.updated_at,
		ar.id, ar.name, ar.genres, ar.description, ar.city
		from bookings b
		left join artists ar on (b.artist_id = ar.id)
		where processed = 0
		order by b.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return bookings, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Bookings
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.Processed,
			&i.ArtistID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Artist.ID,
			&i.Artist.Name,
			&i.Artist.Genres,
			&i.Artist.Description,
			&i.Artist.City,
		)

		if err != nil {
			return bookings, err
		}
		bookings = append(bookings, i)
	}

	if err = rows.Err(); err != nil {
		return bookings, err
	}

	return bookings, nil
}

// Get all Booking Options
func (m *postgresDBRepo) AllBookingOptions() ([]models.BookingOptions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var options []models.BookingOptions

	query := `select id, title, description, price, artist_id, created_at, updated_at from booking_options order by created_at asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return options, err
	}
	defer rows.Close()

	for rows.Next() {
		var option models.BookingOptions
		err := rows.Scan(
			&option.ID,
			&option.Title,
			&option.Description,
			&option.Price,
			&option.ArtistID,
			&option.CreatedAt,
			&option.UpdatedAt,
		)
		if err != nil {
			return options, err
		}
		options = append(options, option)
	}

	if err = rows.Err(); err != nil {
		return options, err
	}

	return options, nil
}

// Get all Booking Options
func (m *postgresDBRepo) AllArtistBookingOptions(id int) ([]models.BookingOptions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var options []models.BookingOptions

	query := `select id, title, description, price, artist_id, created_at, updated_at from booking_options where artist_id = $1 order by created_at asc`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return options, err
	}
	defer rows.Close()

	for rows.Next() {
		var option models.BookingOptions
		err := rows.Scan(
			&option.ID,
			&option.Title,
			&option.Description,
			&option.Price,
			&option.ArtistID,
			&option.CreatedAt,
			&option.UpdatedAt,
		)
		if err != nil {
			return options, err
		}
		options = append(options, option)
	}

	if err = rows.Err(); err != nil {
		return options, err
	}

	return options, nil
}

// Inserts a new Boking Option into the database
func (repo *postgresDBRepo) CreateBookingOption(option models.BookingOptions) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into booking_options (title, description, price, artist_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := repo.DB.ExecContext(context, query, option.Title, option.Description, option.Price, option.ArtistID, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

// Get a booking option by id
func (repo *postgresDBRepo) GetBookingOptionByID(id int) (models.BookingOptions, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var option models.BookingOptions

	query := `
		select id, title, description, price, artist_id, created_at, updated_at from booking_options where id = $1
	`

	row := repo.DB.QueryRowContext(context, query, id)
	err := row.Scan(
		&option.ID,
		&option.Title,
		&option.Description,
		&option.Price,
		&option.ArtistID,
		&option.CreatedAt,
		&option.UpdatedAt,
	)

	if err != nil {
		return option, err
	}

	return option, nil
}

// UpdateBookingOption updates an option in the database
func (m *postgresDBRepo) UpdateBookingOption(option models.BookingOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update booking_options set title = $1, description = $2, price = $3, artist_id = $4, updated_at = $5
		where id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		option.Title,
		option.Description,
		option.Price,
		option.ArtistID,
		time.Now(),
		option.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
