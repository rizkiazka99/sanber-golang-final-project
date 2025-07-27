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

func PostItem(ctx *gin.Context) {
	var item models.Item

	userId, role, accessTokenValidation := middleware.ValidateAccessToken(ctx)

	if accessTokenValidation != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": accessTokenValidation,
		})
	} else {
		if role == "user" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Only admins are allowed to perform this action",
			})
		} else {
			if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to parse form data",
				})
				return
			} else {
				now := time.Now()
				price, err := strconv.Atoi(ctx.PostForm("price"))
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Price must be a number"})
					return
				}
				stock, err := strconv.Atoi(ctx.PostForm("stock"))
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Stock must be a number"})
					return
				}

				item = models.Item{
					Id:          utils.IDGenerator(),
					ItemName:    ctx.PostForm("item_name"),
					Description: ctx.PostForm("desc"),
					Price:       price,
					Stock:       stock,
					CreatedAt:   &now,
					CreatedBy:   userId,
					ModifiedAt:  &now,
					ModifiedBy:  userId,
				}

				form, _ := ctx.MultipartForm()
				files := form.File["images"]
				for _, file := range files {
					filename := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), file.Filename)
					if err := ctx.SaveUploadedFile(file, filename); err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"error": "Failed to save image",
						})

						return
					} else {
						item.Images = append(item.Images, models.ItemImages{
							Id:       utils.IDGenerator(),
							ImageUrl: filename,
						})
					}
				}

				repository.CreateItem(item)
				ctx.JSON(http.StatusCreated, gin.H{
					"message": "item created",
				})
			}
		}
	}
}

func GetItems(ctx *gin.Context) {
	items, err := repository.GetItems()

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
				"items": items,
			})
		}
	}
}

func GetItemById(ctx *gin.Context) {
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
			item, err := repository.GetItemById(id)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": "Item doesn't exist",
					})
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
				}
				return
			} else {
				ctx.JSON(http.StatusOK, item)
			}
		}
	}
}

func UpdateItem(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	userId, role, accessTokenValidation := middleware.ValidateAccessToken(ctx)

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
			var input models.Item

			if err := ctx.ShouldBindJSON(&input); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid input",
					"details": err.Error(),
				})
				return
			} else if role == "user" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Only admins are allowed to perform this action",
				})
			} else {
				now := time.Now()

				input.ModifiedAt = &now
				input.ModifiedBy = userId
				rowsUpdated, err := repository.UpdateItem(id, input)

				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to update item",
						"details": err.Error(),
					})
					return
				} else if rowsUpdated == 0 {
					ctx.JSON(http.StatusNotFound, gin.H{
						"message": "No item found with the given ID",
					})
					return
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"message": "Item was successfully updated",
						"rows":    rowsUpdated,
					})
				}
			}
		}
	}
}

func DeleteItem(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	_, role, accessTokenValidation := middleware.ValidateAccessToken(ctx)

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
		} else if role == "user" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Only admins are allowed to perform this action",
			})
		} else {
			rowsDeleted, err := repository.DeleteItem(id)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":  "Failed to delete item",
					"detais": err.Error(),
				})
				return
			} else if rowsDeleted == 0 {
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": "No item found with the given ID",
				})
				return
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Item has been deleted",
					"rows":    rowsDeleted,
				})
			}
		}
	}
}
