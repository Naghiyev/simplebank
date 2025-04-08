package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	db "simple-banking/db/sqlc"
	"simple-banking/util"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(context *gin.Context) {
	var req createUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.Fullname,
		Email:          req.Email,
		HashedPassword: hashedPassword, //fake

	}

	user, err := server.store.CreateUser(context, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := createUserResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	context.JSON(http.StatusCreated, rsp)
}
