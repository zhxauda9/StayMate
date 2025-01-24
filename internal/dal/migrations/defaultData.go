package migrations

import (
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

func Fill(db *gorm.DB) {
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Status: "active", Photo: ""},
		{Name: "Jane Smith", Email: "jane@example.com", Status: "active", Photo: ""},
		{Name: "Michael Johnson", Email: "michael@example.com", Status: "inactive", Photo: ""},
		{Name: "Emily Davis", Email: "emily@example.com", Status: "active", Photo: ""},
		{Name: "Sarah Williams", Email: "sarah@example.com", Status: "inactive", Photo: ""},
		{Name: "David Brown", Email: "david@example.com", Status: "active", Photo: ""},
		{Name: "Linda Garcia", Email: "linda@example.com", Status: "active", Photo: ""},
		{Name: "James Martinez", Email: "james@example.com", Status: "inactive", Photo: ""},
		{Name: "Patricia Lee", Email: "patricia@example.com", Status: "active", Photo: ""},
		{Name: "Robert Wilson", Email: "robert@example.com", Status: "inactive", Photo: ""},
		{Name: "Maria Lopez", Email: "maria@example.com", Status: "active", Photo: ""},
		{Name: "William Walker", Email: "william@example.com", Status: "inactive", Photo: ""},
		{Name: "Elizabeth Moore", Email: "elizabeth@example.com", Status: "active", Photo: ""},
		{Name: "Joseph Taylor", Email: "joseph@example.com", Status: "inactive", Photo: ""},
		{Name: "Linda Anderson", Email: "linda.anderson@example.com", Status: "active", Photo: ""},
		{Name: "Charles Thomas", Email: "charles@example.com", Status: "inactive", Photo: ""},
		{Name: "Nancy Jackson", Email: "nancy@example.com", Status: "active", Photo: ""},
		{Name: "Donald Harris", Email: "donald@example.com", Status: "inactive", Photo: ""},
		{Name: "Barbara Clark", Email: "barbara@example.com", Status: "active", Photo: ""},
	}

	for _, user := range users {
		var existingUser models.User
		if err := db.First(&existingUser, "email = ?", user.Email).Error; err != nil {
			if user.Photo == "" {
				user.Photo = "static/pictures/default/user.jpg"
			}
			db.Create(&user)
		}
	}

	rooms := []models.Room{
		{Number: 101, Class: "single", Price: 100.00, Status: "available", Photo: "", Description: ""},
		{Number: 102, Class: "double", Price: 150.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 103, Class: "suite", Price: 250.00, Status: "available", Photo: "", Description: ""},
		{Number: 104, Class: "single", Price: 120.00, Status: "available", Photo: "", Description: ""},
		{Number: 105, Class: "double", Price: 180.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 106, Class: "suite", Price: 300.00, Status: "available", Photo: "", Description: ""},
		{Number: 107, Class: "single", Price: 90.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 108, Class: "double", Price: 160.00, Status: "available", Photo: "", Description: ""},
		{Number: 109, Class: "suite", Price: 280.00, Status: "available", Photo: "", Description: ""},
		{Number: 110, Class: "single", Price: 110.00, Status: "available", Photo: "", Description: ""},
		{Number: 111, Class: "double", Price: 170.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 112, Class: "suite", Price: 350.00, Status: "available", Photo: "", Description: ""},
		{Number: 113, Class: "single", Price: 95.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 114, Class: "double", Price: 140.00, Status: "available", Photo: "", Description: ""},
		{Number: 115, Class: "suite", Price: 290.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 116, Class: "single", Price: 105.00, Status: "available", Photo: "", Description: ""},
		{Number: 117, Class: "double", Price: 155.00, Status: "available", Photo: "", Description: ""},
		{Number: 118, Class: "suite", Price: 320.00, Status: "available", Photo: "", Description: ""},
		{Number: 119, Class: "single", Price: 115.00, Status: "occupied", Photo: "", Description: ""},
		{Number: 120, Class: "double", Price: 165.00, Status: "available", Photo: "", Description: ""},
	}

	for _, room := range rooms {
		var existingRoom models.Room
		if err := db.First(&existingRoom, "number = ?", room.Number).Error; err != nil {
			if room.Photo == "" {
				room.Photo = "static/pictures/default/room.jpg"
			}
			if room.Description == "" {
				room.Description = "VERY BEAUTIFULLLLL ROOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOm"
			}
			db.Create(&room)
		}
	}

	bookings := []models.Booking{
		{UserID: 1, RoomID: 1, CheckIn: "2025-01-05 14:00:00", CheckOut: "2025-01-07 11:00:00"},
		{UserID: 2, RoomID: 2, CheckIn: "2025-01-06 15:00:00", CheckOut: "2025-01-08 12:00:00"},
		{UserID: 3, RoomID: 3, CheckIn: "2025-01-07 13:00:00", CheckOut: "2025-01-09 10:00:00"},
		{UserID: 4, RoomID: 4, CheckIn: "2025-01-08 14:30:00", CheckOut: "2025-01-10 11:30:00"},
		{UserID: 5, RoomID: 5, CheckIn: "2025-01-09 16:00:00", CheckOut: "2025-01-11 12:30:00"},
		{UserID: 6, RoomID: 6, CheckIn: "2025-01-10 14:00:00", CheckOut: "2025-01-12 10:00:00"},
		{UserID: 7, RoomID: 7, CheckIn: "2025-01-11 15:00:00", CheckOut: "2025-01-13 11:00:00"},
		{UserID: 8, RoomID: 8, CheckIn: "2025-01-12 16:00:00", CheckOut: "2025-01-14 12:00:00"},
		{UserID: 9, RoomID: 9, CheckIn: "2025-01-13 14:30:00", CheckOut: "2025-01-15 10:00:00"},
		{UserID: 10, RoomID: 10, CheckIn: "2025-01-14 15:00:00", CheckOut: "2025-01-16 11:30:00"},
		{UserID: 11, RoomID: 11, CheckIn: "2025-01-15 14:00:00", CheckOut: "2025-01-17 10:00:00"},
		{UserID: 12, RoomID: 12, CheckIn: "2025-01-16 15:30:00", CheckOut: "2025-01-18 11:00:00"},
		{UserID: 13, RoomID: 13, CheckIn: "2025-01-17 14:00:00", CheckOut: "2025-01-19 12:00:00"},
		{UserID: 14, RoomID: 14, CheckIn: "2025-01-18 16:00:00", CheckOut: "2025-01-20 10:00:00"},
		{UserID: 15, RoomID: 15, CheckIn: "2025-01-19 14:00:00", CheckOut: "2025-01-21 11:00:00"},
		{UserID: 16, RoomID: 16, CheckIn: "2025-01-20 15:00:00", CheckOut: "2025-01-22 12:00:00"},
		{UserID: 17, RoomID: 17, CheckIn: "2025-01-21 14:30:00", CheckOut: "2025-01-23 11:00:00"},
		{UserID: 18, RoomID: 18, CheckIn: "2025-01-22 16:00:00", CheckOut: "2025-01-24 10:30:00"},
		{UserID: 19, RoomID: 19, CheckIn: "2025-01-23 14:00:00", CheckOut: "2025-01-25 11:30:00"},
		{UserID: 20, RoomID: 20, CheckIn: "2025-01-24 15:00:00", CheckOut: "2025-01-26 12:00:00"},
	}

	for _, booking := range bookings {
		var existingBooking models.Booking
		if err := db.First(&existingBooking, "user_id = ? AND room_id = ?", booking.UserID, booking.RoomID).Error; err != nil {
			db.Create(&booking)
		}
	}
}
