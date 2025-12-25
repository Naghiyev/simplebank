package api

import (
	"database/sql"
	"net/http"
	db "simple-banking/db/sqlc"
	"simple-banking/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
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
	rsp := newUserResponse(user)
	context.JSON(http.StatusCreated, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(context *gin.Context) {
	var req loginUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUser(context, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accesstoken, err := server.tokenMaker.CreateToken(req.Username,
		server.config.AccesTokenDuration)

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		accesstoken,
		newUserResponse(user),
	}
	context.JSON(http.StatusOK, rsp)

}
