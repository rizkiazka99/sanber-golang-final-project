package controllers

import (
	"fmt"
	"golang-final-project/middleware"
	"golang-final-project/models"
	"golang-final-project/repository"
	"golang-final-project/utils"
	"net/http"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else if user.Username == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Username and/or password fields cannot be empty",
		})
	} else if user.Role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please specify a role",
		})
	} else {
		user.Id = utils.IDGenerator()
		user.Password, _ = middleware.HashPassword(user.Password)

		fmt.Println(user)

		e := repository.CreateUser(user)
		if e != "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": e,
			})
		} else {
			userResponse := models.BuildUserResponse(user)

			ctx.JSON(http.StatusCreated, gin.H{
				"user": userResponse,
			})
		}
	}
}

func Login(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else if user.Username == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Username and/or password fields cannot be empty",
		})
	} else {
		userData, err := repository.Login(user.Username, user.Password)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			var data models.LoggedIn
			idStr := strconv.Itoa(userData.Id)

			accessToken, e := middleware.GenerateJwt(idStr, userData.Role)

			if e != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(),
				})
			} else {
				data.Id = userData.Id
				data.Username = userData.Username
				data.Role = userData.Role
				data.AccessToken = accessToken
				data.TokenExpirationTime = time.Now().Add(time.Hour * 1)

				_, e := repository.AssignAccessToken(data.Id, data.AccessToken, data.TokenExpirationTime)

				if e != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": e.Error(),
					})
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"data": data,
					})
				}
			}
		}
	}
}
