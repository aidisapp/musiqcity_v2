# MusiqCity - Music Events Management Platform

## Collaborators

- Prosper Atu
- Roselene Odimgbe
- Idongesit Ekanem

---

## Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/musiqcity.git
   cd musiqcity
   ```

2. **Install dependencies:**
   Make sure you have Go and PostgreSQL installed. Then run:
   ```bash
   go mod tidy
   ```

3. **Set up the database:**
   Configure PostgreSQL and create a database. Could you update the database connection details in your configuration file?

4. **Run the application:**
   ```bash
   go run main.go
   ```

5. **Access the application:**
   Open your web browser and go to `http://localhost:8080` to start using MusiqCity.

---

## Usage Guidelines

- **Admin Dashboard:** Admins can manage users, bookings, and artists via a dedicated dashboard.
- **User Dashboard:** Users can view upcoming events, make bookings, and receive updates.
- **Artiste Dashboard:** Artistes have access to manage their profiles, events, and view bookings.

---

## Architecture Overview

The MusiqCity application uses a simple client-server architecture:

- **Frontend (HTML/CSS):** Provides the user interface for interactions.
- **Backend (Go):** Handles requests, processes data, and interacts with the PostgreSQL database.
- **PostgreSQL Database:** Stores user, event, and booking data securely.
- **Email Notification Service:** Sends automated booking confirmation and updates.

---

## Lessons Learned

- We deepened our understanding of core backend development with Go.
- We improved our skills in managing sessions, security, and event-driven notifications.
- Collaborative problem-solving helped us navigate challenges effectively.

---

## Next Steps

- **Mobile Responsiveness:** A key feature for future development is ensuring the platform is fully optimized for mobile devices.
- **Payment Integration:** In the next version, we plan to implement a secure payment system for booking artistes.

---

## Conclusion

MusiqCity has been a rewarding project that strengthened our bond as a team. We are proud of the progress we've made and excited to continue improving the platform in the future.

---

## Connect

- **GitHub:** [MusiqCity GitHub Repository](https://github.com/aidisapp/musiqcity_v2)
- **Live Project:** [MusiqCity Live Version](https://musiqcity-live-url.com)
```
