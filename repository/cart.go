package repository

import (
	"database/sql"
	"fmt"
	"golang-final-project/config"
	"golang-final-project/models"
	"log"
	"time"
)

func CreateCart(body models.PostCartBody) {
	var cart models.PostCartBody

	tx, err := config.Db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	cartQuery := `
	INSERT INTO carts (id, user_id, created_at, total_price, payment_method, payment_status)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
	`

	err = tx.QueryRow(
		cartQuery,
		body.Id,
		body.UserId,
		body.CreatedAt,
		body.TotalPrice,
		body.PaymentMethod,
		body.PaymentStatus,
	).Scan(
		&cart.Id,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.TotalPrice,
		&cart.PaymentMethod,
		&cart.PaymentStatus,
	)

	if err != nil {
		panic(err)
	} else {
		cartItemQuery := `
		INSERT INTO cart_items (id, cart_id, item_id, quantity)
		VALUES ($1, $2, $3, $4)
		`

		for _, item := range body.Items {
			_, err = tx.Exec(
				cartItemQuery,
				item.Id,
				body.Id,
				item.ItemId,
				item.Quantity,
			)
			if err != nil {
				panic(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func GetCarts() ([]models.Cart, error) {
	var results []models.Cart

	query := `
	SELECT
		i.id, i.user_id, i.created_at, i.total_price, i.payment_method, i.payment_status,
		ii.id, ii.cart_id, ii.item_id, ii.quantity,
		iii.id, iii.item_name, iii.price,
		iiii.id, iiii.item_id, iiii.image_url
	FROM carts i
	LEFT JOIN cart_items ii ON i.id = ii.cart_id
	LEFT JOIN items iii ON iii.id = ii.item_id
	LEFT JOIN items_images iiii ON iiii.item_id = iii.id
	`

	rows, err := config.Db.Query(query)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Maps to track grouping
	cartMap := make(map[int]*models.Cart)
	cartItemMap := make(map[int]*models.CartItem)

	for rows.Next() {
		var (
			cartID, userID, cartItemID, cartItemCartID, cartItemItemID, quantity int
			itemID, price                                                        int
			itemName                                                             string
			paymentMethod, paymentStatus                                         string
			totalPrice                                                           sql.NullInt64
			imageID, imageItemID                                                 sql.NullInt64
			imageURL                                                             sql.NullString
			createdAt                                                            time.Time
		)

		err := rows.Scan(
			&cartID, &userID, &createdAt, &totalPrice, &paymentMethod, &paymentStatus,
			&cartItemID, &cartItemCartID, &cartItemItemID, &quantity,
			&itemID, &itemName, &price,
			&imageID, &imageItemID, &imageURL,
		)
		if err != nil {
			return nil, err
		}

		// Handle Cart
		cart, exists := cartMap[cartID]
		if !exists {
			cart = &models.Cart{
				Id:            cartID,
				UserId:        userID,
				CreatedAt:     createdAt,
				TotalPrice:    int(totalPrice.Int64),
				CartItems:     []models.CartItem{},
				PaymentMethod: paymentMethod,
				PaymentStatus: paymentStatus,
			}
			cartMap[cartID] = cart
		}

		// Handle CartItem
		cartItem, exists := cartItemMap[cartItemID]
		if !exists && cartItemID != 0 {
			item := &models.Item{
				Id:       itemID,
				ItemName: itemName,
				Price:    price,
				Images:   []models.ItemImages{},
			}

			cartItem = &models.CartItem{
				Id:       cartItemID,
				CartId:   cartItemCartID,
				ItemId:   cartItemItemID,
				Quantity: quantity,
				Item:     item,
			}

			cart.CartItems = append(cart.CartItems, *cartItem)
			cartItemMap[cartItemID] = cartItem
		}

		// Handle Images
		if cartItem != nil && imageID.Valid && imageURL.Valid {
			cartItem.Item.Images = append(cartItem.Item.Images, models.ItemImages{
				Id:       int(imageID.Int64),
				ItemId:   int(imageItemID.Int64),
				ImageUrl: imageURL.String,
			})
		}
	}

	for _, cart := range cartMap {
		results = append(results, *cart)
	}

	return results, nil
}

func GetCartById(id int64) (*models.Cart, error) {
	var cart *models.Cart

	query := `
	SELECT
		i.id, i.user_id, i.created_at, i.total_price, i.payment_method, i.payment_status,
		ii.id, ii.cart_id, ii.item_id, ii.quantity,
		iii.id, iii.item_name, iii.price,
		iiii.id, iiii.item_id, iiii.image_url
	FROM carts i
	LEFT JOIN cart_items ii ON i.id = ii.cart_id
	LEFT JOIN items iii ON iii.id = ii.item_id
	LEFT JOIN items_images iiii ON iiii.item_id = iii.id
	WHERE i.id = $1
	`

	rows, err := config.Db.Query(query, id)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cartItemMap := make(map[int]*models.CartItem)

	for rows.Next() {
		var (
			cartID, userID, cartItemID, cartItemCartID, cartItemItemID, quantity int
			itemID, price                                                        int
			itemName                                                             string
			paymentMethod, paymentStatus                                         string
			totalPrice                                                           sql.NullInt64
			imageID, imageItemID                                                 sql.NullInt64
			imageURL                                                             sql.NullString
			createdAt                                                            time.Time
		)

		err := rows.Scan(
			&cartID, &userID, &createdAt, &totalPrice, &paymentMethod, &paymentStatus,
			&cartItemID, &cartItemCartID, &cartItemItemID, &quantity,
			&itemID, &itemName, &price,
			&imageID, &imageItemID, &imageURL,
		)
		if err != nil {
			return nil, err
		}

		// Handle Cart
		cart = &models.Cart{
			Id:            cartID,
			UserId:        userID,
			CreatedAt:     createdAt,
			TotalPrice:    int(totalPrice.Int64),
			CartItems:     []models.CartItem{},
			PaymentMethod: paymentMethod,
			PaymentStatus: paymentStatus,
		}

		// Handle CartItem
		cartItem, exists := cartItemMap[cartItemID]
		if !exists && cartItemID != 0 {
			item := &models.Item{
				Id:       itemID,
				ItemName: itemName,
				Price:    price,
				Images:   []models.ItemImages{},
			}

			cartItem = &models.CartItem{
				Id:       cartItemID,
				CartId:   cartItemCartID,
				ItemId:   cartItemItemID,
				Quantity: quantity,
				Item:     item,
			}

			cart.CartItems = append(cart.CartItems, *cartItem)
			cartItemMap[cartItemID] = cartItem
		}

		// Handle Images
		if cartItem != nil && imageID.Valid && imageURL.Valid {
			cartItem.Item.Images = append(cartItem.Item.Images, models.ItemImages{
				Id:       int(imageID.Int64),
				ItemId:   int(imageItemID.Int64),
				ImageUrl: imageURL.String,
			})
		}
	}

	if cart == nil {
		return nil, sql.ErrNoRows
	} else {
		return cart, nil
	}
}

func GetCartsByUserId(id int64) ([]models.Cart, error) {
	var results []models.Cart

	query := `
	SELECT
		i.id, i.user_id, i.created_at, i.total_price, i.payment_method, i.payment_status,
		ii.id, ii.cart_id, ii.item_id, ii.quantity,
		iii.id, iii.item_name, iii.price,
		iiii.id, iiii.item_id, iiii.image_url
	FROM carts i
	LEFT JOIN cart_items ii ON i.id = ii.cart_id
	LEFT JOIN items iii ON iii.id = ii.item_id
	LEFT JOIN items_images iiii ON iiii.item_id = iii.id
	WHERE i.user_id = $1
	`

	rows, err := config.Db.Query(query, id)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Maps to track grouping
	cartMap := make(map[int]*models.Cart)
	cartItemMap := make(map[int]*models.CartItem)

	for rows.Next() {
		var (
			cartID, userID, cartItemID, cartItemCartID, cartItemItemID, quantity int
			itemID, price                                                        int
			itemName                                                             string
			paymentMethod, paymentStatus                                         string
			totalPrice                                                           sql.NullInt64
			imageID, imageItemID                                                 sql.NullInt64
			imageURL                                                             sql.NullString
			createdAt                                                            time.Time
		)

		err := rows.Scan(
			&cartID, &userID, &createdAt, &totalPrice, &paymentMethod, &paymentStatus,
			&cartItemID, &cartItemCartID, &cartItemItemID, &quantity,
			&itemID, &itemName, &price,
			&imageID, &imageItemID, &imageURL,
		)
		if err != nil {
			return nil, err
		}

		// Handle Cart
		cart, exists := cartMap[cartID]
		if !exists {
			cart = &models.Cart{
				Id:            cartID,
				UserId:        userID,
				CreatedAt:     createdAt,
				TotalPrice:    int(totalPrice.Int64),
				CartItems:     []models.CartItem{},
				PaymentMethod: paymentMethod,
				PaymentStatus: paymentStatus,
			}
			cartMap[cartID] = cart
		}

		// Handle CartItem
		cartItem, exists := cartItemMap[cartItemID]
		if !exists && cartItemID != 0 {
			item := &models.Item{
				Id:       itemID,
				ItemName: itemName,
				Price:    price,
				Images:   []models.ItemImages{},
			}

			cartItem = &models.CartItem{
				Id:       cartItemID,
				CartId:   cartItemCartID,
				ItemId:   cartItemItemID,
				Quantity: quantity,
				Item:     item,
			}

			cart.CartItems = append(cart.CartItems, *cartItem)
			cartItemMap[cartItemID] = cartItem
		}

		// Handle Images
		if cartItem != nil && imageID.Valid && imageURL.Valid {
			cartItem.Item.Images = append(cartItem.Item.Images, models.ItemImages{
				Id:       int(imageID.Int64),
				ItemId:   int(imageItemID.Int64),
				ImageUrl: imageURL.String,
			})
		}
	}

	for _, cart := range cartMap {
		results = append(results, *cart)
	}

	return results, nil
}

// func UpdateCart(cartId int64, updates []models.CartItemUpdate) error {
// 	if len(updates) == 0 {
// 		return nil
// 	} else {
// 		query := `UPDATE cart_items SET quantity = CASE item_id`
// 		args := []interface{}{cartId}
// 		idPlaceholders := ""

// 		for i, u := range updates {
// 			argPos := i*2 + 2
// 			query += fmt.Sprintf(" WHEN $%d THEN $%d", argPos, argPos+1)
// 			args = append(args, u.ItemId, u.Quantity)

// 			if i > 0 {
// 				idPlaceholders += ", "
// 			}
// 			idPlaceholders += fmt.Sprintf("$%d", argPos)
// 		}

// 		query += fmt.Sprintf(" ELSE quantity END WHERE cart_id = $1 AND id IN (%s)", idPlaceholders)

// 		fmt.Println("Final query:", query)
// 		fmt.Println("Args:", args)

// 		_, err := config.Db.Exec(query, args...)
// 		return err
// 	}
// }

// func DeleteCartItems(cartId int64, itemIds []int) error {
// 	query := `DELETE FROM cart_items WHERE cart_id = $1 AND item_id IN (`
// 	args := []interface{}{cartId}
// 	placeholders := ""

// 	for i, id := range itemIds {
// 		if i > 0 {
// 			placeholders += ", "
// 		}
// 		placeholders += fmt.Sprintf("$%d", i+2)
// 		args = append(args, id)
// 	}

// 	query += placeholders + ")"

// 	_, err := config.Db.Exec(query, args...)
// 	return err
// }

func DeleteCart(id int64) (int64, error) {
	sqlStatement := `DELETE from carts WHERE id = $1`

	res, err := config.Db.Exec(sqlStatement, id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

func PayCart(id int64) (int64, error) {
	tx, err := config.Db.Begin()
	if err != nil {
		return 0, err
	}

	updateCartQuery := `
	UPDATE carts
	SET payment_status = $2
	WHERE id = $1
	`

	_, err = tx.Exec(
		updateCartQuery,
		id,
		"Paid",
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	updateStockQuery := `
	UPDATE items
	SET stock = stock - ci.quantity
	FROM cart_items ci
	WHERE items.id = ci.item_id AND ci.cart_id = $1
	`

	_, err = tx.Exec(updateStockQuery, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if any stock went negative (business rule check)
	stockCheckQuery := `
	SELECT COUNT(*) FROM items
	WHERE stock < 0
	`

	var negativeCount int
	err = tx.QueryRow(stockCheckQuery).Scan(&negativeCount)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if negativeCount > 0 {
		tx.Rollback()
		return 0, fmt.Errorf("not enough stock for one or more items")
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return 1, nil
}
