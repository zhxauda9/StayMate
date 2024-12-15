package handler

import (
	"net/http"

	"github.com/zhxauda9/StayMate/internal/service"
)

type userHandler struct {
	booking_service service.User_serveice_Impl
}

func NewUserHandler(booking_service service.User_serveice_Impl) *userHandler {
	return &userHandler{booking_service: booking_service}
}

func (h *userHandler) PostUser(w http.ResponseWriter, r *http.Request) {

}
func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

}
func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
