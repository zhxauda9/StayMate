package dal

import "database/sql"

type Booking_repo_Impl interface {
}

type booking_repository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *booking_repository {
	return &booking_repository{db: db}
}
