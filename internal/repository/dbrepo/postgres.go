package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aidisapp/MusiqCity/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate authenticates a user
func (repo *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := repo.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (repo *postgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (repo *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// Close this transaction if unable to run this statement within 3 seconds
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	insertStatement := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := repo.DB.QueryRowContext(context, insertStatement, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (repo *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insertStatement := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.DB.ExecContext(context, insertStatement, res.StartDate, res.EndDate, res.RoomID, res.ReservationID, time.Now(), time.Now(), res.RestrictionID)

	if err != nil {
		return err
	}

	return nil
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

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (repo *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select
			r.id, r.room_name
		from
			rooms r
		where r.id not in 
		(select room_id from room_restrictions rr where $1 <= rr.end_date and $2 >= rr.start_date);
	`

	rows, err := repo.DB.QueryContext(context, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// Get all rooms
func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `SELECT id, room_name, price, image_src, description, created_at, updated_at FROM rooms ORDER BY room_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.Price,
			&room.ImageSource,
			&room.Description,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (repo *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, price, image_src, description, created_at, updated_at from rooms where id = $1
	`

	row := repo.DB.QueryRowContext(context, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.Price,
		&room.ImageSource,
		&room.Description,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

// UpdateRoom updates a room in the database
func (m *postgresDBRepo) UpdateRoom(room models.Room) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update rooms set room_name = $1, price = $2, image_src = $3, description = $4, updated_at = $5
		where id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		room.RoomName,
		room.Price,
		room.ImageSource,
		room.Description,
		time.Now(),
		room.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Inserts a room into the database
func (repo *postgresDBRepo) InsertRoom(room models.Room) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into rooms (room_name, price, image_src, description, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := repo.DB.ExecContext(context, query, room.RoomName, room.Price, room.ImageSource, room.Description, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

// DeleteRoom deletes a room
func (m *postgresDBRepo) DeleteRoom(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from rooms where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
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

// AllReservations returns a slice of all reservations
func (repo *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
	`

	rows, err := repo.DB.QueryContext(context, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at, 
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where processed = 0
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationByID returns one reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, 
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)

	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

// UpdateReservation updates a reservation in the database
func (m *postgresDBRepo) UpdateReservation(user models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
		where id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation deletes one reservation by id
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "delete from reservations where id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
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

// GetRestrictionsForCurrentRoom returns restrictions for a room by date range
func (m *postgresDBRepo) GetRestrictionsForCurrentRoom(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `
		select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
		from room_restrictions where $1 < end_date and $2 >= start_date
		and room_id = $3
`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id,
		created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate, id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// InsertTodoList inserts a new todo list into the database
func (repo *postgresDBRepo) InsertTodoList(todo models.TodoList) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into todo_list (todo, user_id, created_at, updated_at) values($1, $2, $3, $4)`

	_, err := repo.DB.ExecContext(context, query, todo.Todo, todo.UserID, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

// GetTodoListByUserID gets all todo for a user by user_id
func (m *postgresDBRepo) GetTodoListByUserID(id int) ([]models.TodoList, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todoList []models.TodoList

	query := `
		select id, todo, user_id, created_at, updated_at from todo_list where user_id = $1 
		order by created_at asc
	`

	rows, err := m.DB.QueryContext(context, query, id)
	if err != nil {
		return todoList, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.TodoList
		err := rows.Scan(
			&i.ID,
			&i.Todo,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		)

		if err != nil {
			return todoList, err
		}
		todoList = append(todoList, i)
	}

	if err = rows.Err(); err != nil {
		return todoList, err
	}

	return todoList, nil
}

// DeleteTodo deletes a todo
func (m *postgresDBRepo) DeleteTodo(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from todo_list where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
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
