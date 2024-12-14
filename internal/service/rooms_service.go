package service

import "github.com/zhxauda9/StayMate/internal/dal"

type rooms_service struct {
	rooms_repo dal.Rooms_repo_Impl
}

func NewRoomsService(rooms_repo dal.Rooms_repo_Impl) *rooms_service {
	return &rooms_service{rooms_repo}
}
