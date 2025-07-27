package controllers

import (
	"database/sql"
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

func PostCart(ctx *gin.Context) {
	var postCartBody models.PostCartBody

	userId, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if err := ctx.ShouldBindJSON(&postCartBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			cartId := utils.IDGenerator()
			userIdInt, _ := strconv.Atoi(userId)
			createdAt := time.Now()

			for i := range postCartBody.Items {
				postCartBody.Items[i].Id = utils.IDGenerator()
				postCartBody.Items[i].CartId = cartId
			}

			postCartBody.Id = cartId
			postCartBody.UserId = userIdInt
			postCartBody.CreatedAt = createdAt
			postCartBody.PaymentStatus = "Pending"

			repository.CreateCart(postCartBody)
			ctx.JSON(http.StatusCreated, gin.H{
				"message": "cart added",
			})
		}
	}
}

func GetCarts(ctx *gin.Context) {
	carts, err := repository.GetCarts()

	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"carts": carts,
			})
		}
	}
}

func GetCartById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
		} else {
			cart, err := repository.GetCartById(id)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": "Cart doesn't exist",
					})
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
				}
				return
			} else {
				ctx.JSON(http.StatusOK, cart)
			}
		}
	}
}

func GetCartsByUserId(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
		} else {
			carts, err := repository.GetCartsByUserId(id)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": "Cart doesn't exist",
					})
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
				}
				return
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"carts": carts,
				})
			}
		}
	}
}

// func UpdateCart(ctx *gin.Context) {
// 	idParam := ctx.Param("id")
// 	id, err := strconv.ParseInt(idParam, 10, 64)

// 	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

// 	if accessTokenValidation != "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": accessTokenValidation,
// 		})
// 	} else {
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": "Invalid ID",
// 			})
// 			return
// 		} else {
// 			var input []models.CartItemUpdate

// 			if err := ctx.ShouldBindJSON(&input); err != nil {
// 				ctx.JSON(http.StatusBadRequest, gin.H{
// 					"error":   "Invalid input",
// 					"details": err.Error(),
// 				})
// 				return
// 			} else {
// 				err := repository.UpdateCart(id, input)

// 				if err != nil {
// 					fmt.Println(input)
// 					ctx.JSON(http.StatusInternalServerError, gin.H{
// 						"error":   "Failed to update cart",
// 						"details": err.Error(),
// 					})
// 					return
// 				} else {
// 					ctx.JSON(http.StatusOK, gin.H{
// 						"message": "Cart was successfully updated",
// 					})
// 				}
// 			}
// 		}
// 	}
// }

// func DeleteCartItems(ctx *gin.Context) {
// 	idParam := ctx.Param("id")
// 	id, err := strconv.ParseInt(idParam, 10, 64)

// 	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

// 	if accessTokenValidation != "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": accessTokenValidation,
// 		})
// 	} else {
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": "Invalid ID",
// 			})
// 			return
// 		} else {
// 			var itemIds []int

// 			err := repository.DeleteCartItems(id, itemIds)

// 			if err != nil {
// 				ctx.JSON(http.StatusInternalServerError, gin.H{
// 					"error":  "Failed to delete cart item",
// 					"detais": err.Error(),
// 				})
// 				return
// 			} else {
// 				ctx.JSON(http.StatusOK, gin.H{
// 					"message": "Cart items have been deleted",
// 				})
// 			}
// 		}
// 	}
// }

func DeleteCart(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		} else {
			rowsDeleted, err := repository.DeleteCart(id)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":  "Failed to delete cart",
					"detais": err.Error(),
				})
				return
			} else if rowsDeleted == 0 {
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": "No cart found with the given ID",
				})
				return
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Cart has been deleted",
					"rows":    rowsDeleted,
				})
			}
		}
	}
}

func PayCart(ctx *gin.Context) {
	idParam := ctx.Param("cart_id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	_, _, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		} else {
			var input models.CartPayment

			if err := ctx.ShouldBindJSON(&input); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid input",
					"details": err.Error(),
				})
				return
			} else {
				_, err := repository.PayCart(id)

				if err != nil {
					fmt.Println(input)
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Payment Failed",
						"details": err.Error(),
					})
					return
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"message": "Payment Successful",
					})
				}
			}
		}
	}
}
