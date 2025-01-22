package service

import (
	"errors"
	"fmt"

	"github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/models"
)

type roomService struct {
	roomRepo postgres.RoomRepo
}

type RoomServ interface {
	CreateRoom(room models.Room) error
	GetRoomByID(id int) (models.Room, error)
	GetAllRooms(sort, filterStart, filterEnd string, page int) ([]models.Room, error)
	UpdateRoom(id int, room models.Room) error
	DeleteRoom(id int) error
}

func NewRoomService(roomRepo postgres.RoomRepo) RoomServ {
	return &roomService{roomRepo: roomRepo}
}

func (s *roomService) CreateRoom(room models.Room) error {
	if s.roomRepo.RoomExists(room.ID) {
		return errors.New("room already created")
	}
	return s.roomRepo.CreateRoom(room)
}

func (s *roomService) GetRoomByID(id int) (models.Room, error) {
	room, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return models.Room{}, fmt.Errorf("error in service layer while fetching room by ID: %v", err)
	}
	return room, nil
}

func (s *roomService) GetAllRooms(sort, filterStart, filterEnd string, page int) ([]models.Room, error) {
	const limit = 8
	offset := (page - 1) * limit

	rooms, err := s.roomRepo.GetAllRooms(sort, filterStart, filterEnd, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error in service layer while fetching all rooms: %v", err)
	}

	return rooms, nil
}

func (s *roomService) UpdateRoom(id int, room models.Room) error {
	err := s.roomRepo.UpdateRoom(id, room)
	if err != nil {
		return fmt.Errorf("error in service layer while updating room: %v", err)
	}
	return nil
}

func (s *roomService) DeleteRoom(id int) error {
	err := s.roomRepo.DeleteRoom(id)
	if err != nil {
		return fmt.Errorf("error in service layer while deleting room: %v", err)
	}
	return nil
}
