package dal

import "database/sql"

type Rooms_repo_Impl interface{}

type rooms_repository struct {
	db *sql.DB
}

func NewRoomsRepository(db *sql.DB) *rooms_repository {
	return &rooms_repository{db: db}
}
