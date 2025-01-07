package postgres

import (
	"fmt"

	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

type RoomRepo interface {
	CreateRoom(room models.Room) error
	GetRoomByID(id int) (models.Room, error)
	GetAllRooms(sort, filterStart, filterEnd string, limit, offset int) ([]models.Room, error)
	UpdateRoom(id int, room models.Room) error
	DeleteRoom(id int) error
	RoomExists(roomID int) bool
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepo {
	return &roomRepository{db: db}
}

func (r *roomRepository) CreateRoom(room models.Room) error {
	if err := r.db.Create(&room).Error; err != nil {
		return fmt.Errorf("error inserting room %v", err)
	}
	return nil
}

func (r *roomRepository) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if err := r.db.First(&room, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Room{}, fmt.Errorf("room with ID %d not found", id)
		}
		return models.Room{}, fmt.Errorf("error fetching room by ID: %v", err)
	}
	return room, nil
}

func (r *roomRepository) GetAllRooms(sort, filterStart, filterEnd string, limit, offset int) ([]models.Room, error) {
	var rooms []models.Room
	query := r.db.Model(&models.Room{})

	if filterStart != "" && filterEnd != "" {
		query = query.Where("price >= ? AND price <= ?", filterStart, filterEnd)
	}

	if sort != "" {
		query = query.Order(sort)
	}

	if err := query.Limit(limit).Offset(offset).Find(&rooms).Error; err != nil {
		return nil, fmt.Errorf("error fetching rooms: %v", err)
	}

	return rooms, nil
}

func (r *roomRepository) UpdateRoom(id int, room models.Room) error {
	var existingRoom models.Room
	if err := r.db.First(&existingRoom, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("room with ID %d not found", id)
		}
		return fmt.Errorf("error fetching room for update: %v", err)
	}
	if err := r.db.Model(&existingRoom).Updates(room).Error; err != nil {
		return fmt.Errorf("error updating room: %v", err)
	}
	return nil
}

func (r *roomRepository) DeleteRoom(id int) error {
	var room models.Room
	if err := r.db.First(&room, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("room with ID %d not found", id)
		}
		return fmt.Errorf("error fetching room for deletion: %v", err)
	}
	if err := r.db.Delete(&room).Error; err != nil {
		return fmt.Errorf("error deleting room: %v", err)
	}
	return nil
}

func (r *roomRepository) RoomExists(roomID int) bool {
	var count int64
	r.db.Model(&models.Room{}).
		Where("room_id = ?", roomID).
		Count(&count)
	return count > 0
}
