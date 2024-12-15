package service

import "github.com/zhxauda9/StayMate/internal/dal"

type User_serveice_Impl interface {
}

type user_service struct {
	rooms_repo dal.User_repo_Impl
}

func NewUserService(rooms_repo dal.User_repo_Impl) *user_service {
	return &user_service{rooms_repo}
}
