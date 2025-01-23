package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type roomHandler struct {
	roomService service.RoomServ
	validate    *validator.Validate
}

func NewRoomHandler(roomService service.RoomServ) *roomHandler {
	return &roomHandler{
		roomService: roomService,
		validate:    validator.New(),
	}
}

func (h *roomHandler) PostRoom(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to create a new room.")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	file, fileheader, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoPath := "/static/pictures/storage" + filepath.Base(fileheader.Filename)
	outFile, err := os.Create(photoPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	var room models.Room
	err = json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		l.Log.Error().Err(err).Msg("Error decoding room data")
		http.Error(w, "Error decoding room data", http.StatusBadRequest)
		return
	}

	room.Photo = photoPath
	err = h.roomService.CreateRoom(room)
	if err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error creating room: %v", err))
		http.Error(w, "Error creating room", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Msg("Room created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

func (h *roomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to fetch all rooms.")

	filterStart := r.URL.Query().Get("filterStart")
	filterEnd := r.URL.Query().Get("filterEnd")
	sort := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	rooms, err := h.roomService.GetAllRooms(sort, filterStart, filterEnd, page)
	if err != nil {
		l.Log.Error().Err(err).Msg("Error fetching rooms")
		http.Error(w, "Error fetching rooms", http.StatusInternalServerError)
		return
	}

	l.Log.Info().Msg("Fetched all rooms successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

func (h *roomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to fetch a specific room.")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Str("ID", idStr).Err(err).Msg("Invalid room ID")
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	room, err := h.roomService.GetRoomByID(id)
	if err != nil {
		l.Log.Error().Err(err).Int("RoomID", id).Msg("Error fetching room")
		http.Error(w, "Error fetching room", http.StatusNotFound)
		return
	}
	l.Log.Info().Int("RoomID", id).Msg("Fetched room successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func (h *roomHandler) PutRoom(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to update a room.")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Str("ID", idStr).Err(err).Msg("Invalid room ID")
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	l.Log.Debug().Int("RoomID", id).Msg("Parsed room ID successfully")

	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		l.Log.Error().Err(err).Msg("Error decoding room data")
		http.Error(w, "Error decoding room data", http.StatusBadRequest)
		return
	}

	err = h.roomService.UpdateRoom(id, room)
	if err != nil {
		l.Log.Error().Err(err).Int("RoomID", id).Msg("Error updating room")
		http.Error(w, "Error updating room", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Int("RoomID", id).Msg("Updated room successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func (h *roomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to delete a room.")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Str("ID", idStr).Err(err).Msg("Invalid room ID")
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	err = h.roomService.DeleteRoom(id)
	if err != nil {
		l.Log.Error().Err(err).Int("RoomID", id).Msg("Error deleting room")
		http.Error(w, "Error deleting room", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Int("RoomID", id).Msg("Deleted room successfully")

	w.WriteHeader(http.StatusNoContent)
}
