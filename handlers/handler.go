package handlers

import "github.com/Toskosz/everythingreviewed/models"

type Handler struct {
	userService models.InterfaceUserService
}

func NewHandler(user models.InterfaceUserService) Handler {
	return Handler{
		userService: user,
	}
}
