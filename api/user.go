package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	db "simple-banking/db/sqlc"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname"`
	Email    string `json:"email" binding:"required"`
}

func (server *Server) createUser(context *gin.Context) {
	var req createUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.Fullname,
		Email:          req.Email,
		HashedPassword: "asdf", //fake

	}

	user, err := server.store.CreateUser(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	context.JSON(http.StatusCreated, user)
}
