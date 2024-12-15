package dal

import "database/sql"

type User_repo_Impl interface{}

type user_repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *user_repository {
	return &user_repository{db: db}
}
