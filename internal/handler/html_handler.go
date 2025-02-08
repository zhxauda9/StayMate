package handler

import (
	"net/http"
	"path/filepath"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "home.html"))
}

func ServeAdmin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "admin.html"))
}

func ServeBookings(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "bookings.html"))
}

func ServeRooms(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "rooms.html"))
}

func ServeUsers(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "users.html"))
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "login.html"))
}

func ServeRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "register.html"))
}

func ServeProfile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "profile.html"))
}

func ServeEmailVerify(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "email-confirm.html"))
}

// /admin/chats
func ServeAdminChats(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "admin-chats.html"))
}
