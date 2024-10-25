package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aidisapp/musiqcity_v2/internal/config"
	"github.com/aidisapp/musiqcity_v2/internal/driver"
	"github.com/aidisapp/musiqcity_v2/internal/forms"
	"github.com/aidisapp/musiqcity_v2/internal/helpers"
	"github.com/aidisapp/musiqcity_v2/internal/models"
	"github.com/aidisapp/musiqcity_v2/internal/render"
	"github.com/aidisapp/musiqcity_v2/internal/repository"
	"github.com/aidisapp/musiqcity_v2/internal/repository/dbrepo"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// This function creates a new repository
func NewRepo(appConfig *config.AppConfig, dbConnectionPool *driver.DB) *Repository {
	return &Repository{
		App: appConfig,
		DB:  dbrepo.NewPostgresRepo(dbConnectionPool.SQL, appConfig),
	}
}

// This function creates a new repository
func NewTestRepo(appConfig *config.AppConfig) *Repository {
	return &Repository{
		App: appConfig,
		DB:  dbrepo.NewTestRepo(appConfig),
	}
}

// This function NewHandlers, sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// This function handles the Home page and renders the template
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.AllArtists()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artists"] = artists
	render.Template(w, r, "home.page.html", &models.TemplateData{
		Data: data,
	})
}

// This function handles the Home page and renders the template
func (m *Repository) ArtistsPage(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.AllArtists()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artists"] = artists
	render.Template(w, r, "artists.page.html", &models.TemplateData{
		Data: data,
	})
}

// This function handles the single room(Luxery) page and renders the template
func (m *Repository) SingleArtist(w http.ResponseWriter, r *http.Request) {
	urlParams := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlParams[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	artist, err := m.DB.GetArtistByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	options, err := m.DB.AllArtistBookingOptions(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artist"] = artist
	data["options"] = options

	render.Template(w, r, "single-artist.page.html", &models.TemplateData{
		Data: data,
	})
}

// This function handles the About page and renders the template
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// This function handles the Contact page and renders the template
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// Availability json, to handle availability request and send back json
type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    string `json:"room_id"`
}

// This function handles the Admin Login page and renders the template
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	userExists := m.App.Session.GetInt(r.Context(), "user_id")
	if userExists > 0 {
		m.App.Session.Put(r.Context(), "warning", "You are already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// This function handles user details and authentication
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		stringMap := make(map[string]string)
		stringMap["email"] = email
		m.App.Session.Put(r.Context(), "error", "Invalid inputs")
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form:      form,
			StringMap: stringMap,
		})

		return
	}

	id, _, access_level, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid email/password")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "access_level", access_level)
	m.App.Session.Put(r.Context(), "flash", "Login Successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// This function handles the Signup page and renders the template
func (m *Repository) Signup(w http.ResponseWriter, r *http.Request) {
	userExists := m.App.Session.GetInt(r.Context(), "user_id")
	if userExists > 0 {
		m.App.Session.Put(r.Context(), "warning", "You are already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "signup.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// This function handles the posting of user signup form
func (m *Repository) PostSignup(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	user := models.User{}
	stringMap := make(map[string]string)

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
		return
	}

	user.FirstName = r.Form.Get("first_name")
	user.LastName = r.Form.Get("last_name")
	user.Email = r.Form.Get("email")
	user.Password = r.Form.Get("password")
	user.AccessLevel = 0

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email", "password")
	form.MinLength("first_name", 3, 30)
	form.MinLength("last_name", 3, 30)
	form.IsEmail("email")

	if !form.Valid() {
		stringMap["first_name"] = user.FirstName
		stringMap["last_name"] = user.LastName
		stringMap["email"] = user.Email
		m.App.Session.Put(r.Context(), "error", "Invalid inputs")
		render.Template(w, r, "signup.page.html", &models.TemplateData{
			Form:      form,
			StringMap: stringMap,
		})

		return
	}

	userExist, err := m.DB.CheckIfUserEmailExist(user.Email)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if userExist {
		log.Println("Email exist")
		stringMap["first_name"] = user.FirstName
		stringMap["last_name"] = user.LastName
		stringMap["email"] = user.Email
		m.App.Session.Put(r.Context(), "error", "Email exist. Kindly use the forgot password form to reset your password")
		render.Template(w, r, "signup.page.html", &models.TemplateData{
			Form:      form,
			StringMap: stringMap,
		})

		return
	}

	newUserID, err := m.DB.InsertUser(user)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	jwtToken, err := helpers.GenerateJWTToken(newUserID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	user.Token = jwtToken
	m.App.Session.Put(r.Context(), "user", user)

	// Load the env file and get the frontendURL
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	frontendURL := os.Getenv("FRONTEND_URL")

	// Send email notification to user
	htmlBody := fmt.Sprintf(`
	<strong>Verify Your Account</strong><br />
	<p>Dear %s %s, </p>
	<p>Welcome to MusiqCity.</p>
	<strong>Kindly click the link below</strong>
	<a href="%s/verify-email?userid=%d&token=%s", target="_blank">Verify Account</a>
	<p>We hope to see you soon</p>
	`, user.FirstName, user.LastName, frontendURL, newUserID, jwtToken)

	message := models.MailData{
		To:      user.Email,
		From:    "prosperdevstack@gmail.com",
		Subject: "Verify Your Email",
		Content: htmlBody,
	}

	m.App.MailChannel <- message
	// End of emails

	m.App.Session.Put(r.Context(), "flash", "Sign up Successful!!! <br /> Please, check your email and verify your account to continue")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Load the static page and verify user email
func (m *Repository) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "verfy-email.page.html", &models.TemplateData{})

	// get the user id and JWT token string from the url request
	userID, _ := strconv.Atoi(r.URL.Query().Get("userid"))
	tokenString := r.URL.Query().Get("token")

	user, err := m.DB.GetUserByID(userID)
	if err != nil {
		// helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Unable to fetch your account from our server. <br />Please, contact support")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Load the env file and get the JWT secret
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWTSECRET")

	// Parse and verify the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key used for signing the token
		return []byte(jwtSecret), nil
	})
	if err != nil {
		// Handle token parsing or verification errors
		http.Error(w, "Unable to parse token", http.StatusBadRequest)
		m.App.Session.Put(r.Context(), "error", fmt.Sprintf("Unable to parse token. Error: %s", err))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// If the token is valid
	if token.Valid {
		// Update the user's account status as verified
		user.AccessLevel = 1
		err = m.DB.UpdateUserAccessLevel(user)
		if err != nil {
			// helpers.ServerError(w, err)
			m.App.Session.Put(r.Context(), "error", "Unable to update user's access level. Please contact support")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		m.App.Session.Put(r.Context(), "error", fmt.Sprintf("Invalid token. Error: %s", err))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Handle successful email verification, Send email notification to user
	htmlBody := fmt.Sprintf(`
	<strong>Successful</strong><br />
	<p>Hi %s, </p>
	<h3>Your email has been verified.</h3>
	<p>You can now login to your account.</p>
	<strong>Note:<strong>
	<p>You may still need to verify your address and identity before you can list your services on our website. <br />
	Go to your account dashboard and verify your account by providing the required verification documents.</p>
	<p>We hope to see you soon</p>
	`, user.FirstName)

	message := models.MailData{
		To:      user.Email,
		From:    "prosperdevstack@gmail.com",
		Subject: "Email Verified",
		Content: htmlBody,
	}

	m.App.MailChannel <- message
	// End of emails

	m.App.Session.Put(r.Context(), "flash", "Email Verification Successful!!! <br /> Please, login")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// This function logs out the user
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	m.App.Session.Put(r.Context(), "warning", "You have logged out of your account")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles the admin dashborad
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

// Recent ------------------------------------------

// Handles the all-artists route
func (m *Repository) AdminAllArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.AllArtists()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artists"] = artists

	render.Template(w, r, "admin-all-artists.page.html", &models.TemplateData{
		Data: data,
	})
}

// Handles the new-artist route to create a new artist
func (m *Repository) AdminNewArtist(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-new-artist.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// This function POST the new artist form and store them in the database
func (m *Repository) PostAdminNewArtist(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var artist models.Artist

	artist.Name = r.Form.Get("artist_name")
	artist.Genres = r.Form.Get("genres")
	artist.Description = r.Form.Get("description")
	artist.Phone = r.Form.Get("phone")
	artist.Email = r.Form.Get("email")
	artist.City = r.Form.Get("city")
	artist.Facebook = r.Form.Get("facebook")
	artist.Twitter = r.Form.Get("twitter")
	artist.Youtube = r.Form.Get("youtube")
	artist.Logo = r.Form.Get("logo")
	artist.Banner = r.Form.Get("banner")
	artist.FeaturedImage = r.Form.Get("featured_image")

	// Form validations
	form := forms.New(r.PostForm)
	form.Required("artist_name", "genres", "description", "phone", "email")
	form.MinLength("artist_name", 5, 50)
	form.MinLength("description", 5, 20000)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["artist"] = artist
		m.App.Session.Put(r.Context(), "error", "Invalid form input")
		render.Template(w, r, "admin-new-artist.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Insert new artist here
	err = m.DB.CreateArtist(artist)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert new artist into database")
		helpers.ServerError(w, err)
		http.Redirect(w, r, "/admin/artists", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Artist Created Successfully!!!")
	http.Redirect(w, r, "/admin/artists", http.StatusSeeOther)
}

// Handles the single-room route
func (m *Repository) AdminSingleArtist(w http.ResponseWriter, r *http.Request) {
	urlParams := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlParams[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	artist, err := m.DB.GetArtistByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artist"] = artist

	render.Template(w, r, "admin-single-artist.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Handles the single-artist route for POST
func (m *Repository) PostAdminSingleArtist(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	urlParams := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlParams[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	artist, err := m.DB.GetArtistByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	artist.Name = r.Form.Get("artist_name")
	artist.Genres = r.Form.Get("genres")
	artist.Description = r.Form.Get("description")
	artist.Phone = r.Form.Get("phone")
	artist.Email = r.Form.Get("email")
	artist.City = r.Form.Get("city")
	artist.Facebook = r.Form.Get("facebook")
	artist.Twitter = r.Form.Get("twitter")
	artist.Youtube = r.Form.Get("youtube")
	artist.Logo = r.Form.Get("logo")
	artist.Banner = r.Form.Get("banner")
	artist.FeaturedImage = r.Form.Get("featured_image")

	form := forms.New(r.PostForm)
	form.Required("artist_name", "genres", "description", "phone", "email")
	form.MinLength("artist_name", 5, 50)
	form.MinLength("description", 5, 20000)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["artist"] = artist
		m.App.Session.Put(r.Context(), "error", "Invalid inputs")
		render.Template(w, r, "admin-single-artist.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	err = m.DB.UpdateArtist(artist)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Artist Updated Successsfully!!!")
	http.Redirect(w, r, "/admin/artists", http.StatusSeeOther)
}

// This function POST the new booking form, store them in the database and send emails
func (m *Repository) PostCreateBooking(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var booking models.Bookings

	booking.FirstName = r.Form.Get("first_name")
	booking.LastName = r.Form.Get("last_name")
	booking.Email = r.Form.Get("email")
	booking.Phone = r.Form.Get("phone")
	// booking.StartDate = r.Form.Get("start_date")
	// booking.EndDate = r.Form.Get("end_date")
	artistIDStr := r.Form.Get("artist_id")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		log.Printf("Error converting artist_id to integer: %v", err)
	} else {
		booking.ArtistID = artistID
	}

	// Form validations
	form := forms.New(r.PostForm)
	form.Required("title", "price", "description", "artist_id")
	form.MinLength("title", 5, 50)
	form.MinLength("description", 5, 250)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["booking"] = booking
		m.App.Session.Put(r.Context(), "error", "Invalid form input")
		render.Template(w, r, "admin-new-option.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Insert new artist here
	// err = m.DB.CreateBookingOption(booking)
	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "Can't insert new booking option into database")
	// 	helpers.ServerError(w, err)
	// 	http.Redirect(w, r, "/admin/booking-options/new-option", http.StatusTemporaryRedirect)
	// 	return
	// }

	m.App.Session.Put(r.Context(), "flash", "Booking Option Created Successfully!!!")
	http.Redirect(w, r, "/admin/booking-options/new-option", http.StatusSeeOther)
}

// Handles the all-bookings route
func (m *Repository) AdminAllBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := m.DB.AllBookings()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["bookings"] = bookings

	render.Template(w, r, "admin-all-bookings.page.html", &models.TemplateData{
		Data: data,
	})
}

// Handles the new-reservations route
func (m *Repository) AdminNewBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := m.DB.AllNewBookings()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["bookings"] = bookings

	render.Template(w, r, "admin-new-bookings.page.html", &models.TemplateData{
		Data: data,
	})
}

// Handles the all-artists route
func (m *Repository) AdminAllOptions(w http.ResponseWriter, r *http.Request) {
	options, err := m.DB.AllBookingOptions()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["options"] = options

	render.Template(w, r, "admin-all-options.page.html", &models.TemplateData{
		Data: data,
	})
}

// Handles the new-artist route to create a new artist
func (m *Repository) AdminNewOption(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.AllArtists()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artists"] = artists

	render.Template(w, r, "admin-new-option.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// This function POST the new artist form and store them in the database
func (m *Repository) PostAdminNewOption(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var option models.BookingOptions

	option.Title = r.Form.Get("title")
	option.Price = r.Form.Get("price")
	option.Description = r.Form.Get("description")
	artistIDStr := r.Form.Get("artist_id")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		log.Printf("Error converting artist_id to integer: %v", err)
	} else {
		option.ArtistID = artistID
	}

	// Form validations
	form := forms.New(r.PostForm)
	form.Required("title", "price", "description", "artist_id")
	form.MinLength("title", 5, 50)
	form.MinLength("description", 5, 250)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["option"] = option
		m.App.Session.Put(r.Context(), "error", "Invalid form input")
		render.Template(w, r, "admin-new-option.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Insert new artist here
	err = m.DB.CreateBookingOption(option)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert new booking option into database")
		helpers.ServerError(w, err)
		http.Redirect(w, r, "/admin/booking-options/new-option", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Booking Option Created Successfully!!!")
	http.Redirect(w, r, "/admin/booking-options/new-option", http.StatusSeeOther)
}

func (m *Repository) AdminSingleOption(w http.ResponseWriter, r *http.Request) {
	urlParams := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlParams[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	artists, err := m.DB.AllArtists()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	option, err := m.DB.GetBookingOptionByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["option"] = option
	data["artists"] = artists

	render.Template(w, r, "admin-single-option.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// Handles the single-option route for POST
func (m *Repository) PostAdminSingleOption(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	urlParams := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlParams[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	artists, err := m.DB.AllArtists()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	option, err := m.DB.GetBookingOptionByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["artists"] = artists
	data["option"] = option

	option.Title = r.Form.Get("title")
	option.Price = r.Form.Get("price")
	option.Description = r.Form.Get("description")
	artistIDStr := r.Form.Get("artist_id")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		log.Printf("Error converting artist_id to integer: %v", err)
		m.App.Session.Put(r.Context(), "error", "Invalid Artist selected")
		http.Redirect(w, r, "/admin/booking-options", http.StatusSeeOther)
	} else {
		option.ArtistID = artistID
	}

	// Form validations
	form := forms.New(r.PostForm)
	form.Required("title", "price", "description", "artist_id")
	form.MinLength("title", 5, 50)
	form.MinLength("description", 5, 250)

	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Invalid inputs")
		render.Template(w, r, "admin-single-option.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	err = m.DB.UpdateBookingOption(option)
	if err != nil {
		helpers.ServerError(w, err)
		http.Redirect(w, r, "/admin/booking-options", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Booking Option Updated Successsfully!!!")
	http.Redirect(w, r, "/admin/booking-options", http.StatusSeeOther)
}

// This function handles the ListService page and renders the template
func (m *Repository) ListService(w http.ResponseWriter, r *http.Request) {
	userExists := m.App.Session.GetInt(r.Context(), "user_id")
	if userExists < 1 {
		m.App.Session.Put(r.Context(), "warning", "You must be logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "list-service.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// This function handles the posting of ListService form
func (m *Repository) PostListService(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var artist models.Artist

	artist.Name = r.Form.Get("artist_name")
	artist.Genres = r.Form.Get("genres")
	artist.Description = r.Form.Get("description")
	artist.Phone = r.Form.Get("phone")
	artist.Email = r.Form.Get("email")
	artist.City = r.Form.Get("city")
	artist.Facebook = r.Form.Get("facebook")
	artist.Twitter = r.Form.Get("twitter")
	artist.Youtube = r.Form.Get("youtube")
	artist.Logo = r.Form.Get("logo")
	artist.Banner = r.Form.Get("banner")
	artist.FeaturedImage = r.Form.Get("featured_image")

	// Form validations
	form := forms.New(r.PostForm)
	form.Required("artist_name", "genres", "description", "phone", "email")
	form.MinLength("artist_name", 5, 50)
	form.MinLength("description", 5, 20000)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["artist"] = artist
		m.App.Session.Put(r.Context(), "error", "Invalid form input")
		render.Template(w, r, "list-service.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Load the env file and get the frontendURL
	// err = godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }
	// frontendURL := os.Getenv("FRONTEND_URL")

	// Send email notification to user
	// htmlBody := fmt.Sprintf(`
	// <strong>Verify Your Account</strong><br />
	// <p>Dear %s %s, </p>
	// <p>Welcome to MusiqCity.</p>
	// <strong>Kindly click the link below</strong>
	// <a href="%s/verify-email?userid=%d&token=%s", target="_blank">Verify Account</a>
	// <p>We hope to see you soon</p>
	// `, user.FirstName, user.LastName, frontendURL, newUserID, jwtToken)

	// message := models.MailData{
	// 	To:      user.Email,
	// 	From:    "prosperdevstack@gmail.com",
	// 	Subject: "Verify Your Email",
	// 	Content: htmlBody,
	// }

	// m.App.MailChannel <- message
	// End of emails

	m.App.Session.Put(r.Context(), "flash", "Sign up Successful!!! <br /> Please, check your email and verify your account to continue")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
