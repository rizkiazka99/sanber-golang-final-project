package repository

import (
	"database/sql"
	"fmt"
	"golang-final-project/config"
	"golang-final-project/models"
	"time"
)

func CreateItem(i models.Item) {
	tx, err := config.Db.Begin()
	if err != nil {
		panic(err)
	}

	// Insert item
	query := `
		INSERT INTO items (
			id, item_name, description, price, stock,
			created_by, created_at, modified_by, modified_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var insertedId int
	err = tx.QueryRow(
		query,
		i.Id,
		i.ItemName,
		i.Description,
		i.Price,
		i.Stock,
		i.CreatedBy,
		i.CreatedAt,
		i.ModifiedBy,
		i.ModifiedAt,
	).Scan(&insertedId)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	// Insert item images
	for _, img := range i.Images {
		img.ItemId = insertedId
		query := `INSERT INTO items_images (id, item_id, image_url) VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, img.Id, img.ItemId, img.ImageUrl)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	fmt.Printf("Inserted item with ID %d and %d images\n", insertedId, len(i.Images))
}

func GetItems() ([]models.Item, error) {
	var results []models.Item

	sqlStatement := `
	SELECT
		i.id, i.item_name, i.description, i.price, i.stock,
		i.created_at, i.created_by, i.modified_at, i.modified_by,
		ii.id, ii.item_id, ii.image_url
	FROM items i
	LEFT JOIN items_images ii ON i.id = ii.item_id;
	`

	rows, err := config.Db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	itemMap := make(map[int]*models.Item)

	for rows.Next() {
		var (
			itemId                int
			itemName, description string
			price, stock          int
			createdAt, modifiedAt time.Time
			createdBy, modifiedBy string
			imageId, imageItemId  sql.NullInt64  // use NullInt64
			imageUrl              sql.NullString // already correct
		)

		err := rows.Scan(
			&itemId,
			&itemName,
			&description,
			&price,
			&stock,
			&createdAt,
			&createdBy,
			&modifiedAt,
			&modifiedBy,
			&imageId,
			&imageItemId,
			&imageUrl,
		)
		if err != nil {
			return nil, err
		}

		item, exists := itemMap[itemId]
		if !exists {
			item = &models.Item{
				Id:          itemId,
				ItemName:    itemName,
				Description: description,
				Price:       price,
				Stock:       stock,
				CreatedAt:   &createdAt,
				CreatedBy:   createdBy,
				ModifiedAt:  &modifiedAt,
				ModifiedBy:  modifiedBy,
				Images:      []models.ItemImages{},
			}
			itemMap[itemId] = item
		}

		if imageId.Valid && imageItemId.Valid && imageUrl.Valid {
			imageUrl := config.BaseUrl + imageUrl.String
			image := models.ItemImages{
				Id:       int(imageId.Int64),
				ItemId:   int(imageItemId.Int64),
				ImageUrl: imageUrl,
			}
			item.Images = append(item.Images, image)
		}
	}

	for _, item := range itemMap {
		results = append(results, *item)
	}

	return results, nil
}

func GetItemById(id int64) (*models.Item, error) {
	var result *models.Item

	sqlStatement := `
	SELECT
		i.id, i.item_name, i.description, i.price, i.stock,
		i.created_at, i.created_by, i.modified_at, i.modified_by,
		ii.id, ii.item_id, ii.image_url
	FROM items i
	LEFT JOIN items_images ii ON i.id = ii.item_id
	WHERE i.id = $1;
	`

	rows, err := config.Db.Query(sqlStatement, id)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			itemId                int
			itemName, description string
			price, stock          int
			createdAt, modifiedAt time.Time
			createdBy, modifiedBy string
			imageId, imageItemId  sql.NullInt64
			imageUrl              sql.NullString
		)

		err := rows.Scan(
			&itemId,
			&itemName,
			&description,
			&price,
			&stock,
			&createdAt,
			&createdBy,
			&modifiedAt,
			&modifiedBy,
			&imageId,
			&imageItemId,
			&imageUrl,
		)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = &models.Item{
				Id:          itemId,
				ItemName:    itemName,
				Description: description,
				Price:       price,
				Stock:       stock,
				CreatedAt:   &createdAt,
				CreatedBy:   createdBy,
				ModifiedAt:  &modifiedAt,
				ModifiedBy:  modifiedBy,
				Images:      []models.ItemImages{},
			}
		}

		if imageId.Valid && imageItemId.Valid && imageUrl.Valid {
			imageUrl := config.BaseUrl + imageUrl.String
			image := models.ItemImages{
				Id:       int(imageId.Int64),
				ItemId:   int(imageItemId.Int64),
				ImageUrl: imageUrl,
			}
			result.Images = append(result.Images, image)
		}
	}

	if result == nil {
		return nil, sql.ErrNoRows
	} else {
		return result, nil
	}
}

func UpdateItem(id int64, item models.Item) (int64, error) {
	sqlStatement := `
	UPDATE items
	SET item_name = $2, description = $3, price = $4, stock = $5, modified_by = $6, modified_at = $7
	WHERE id = $1;`

	res, err := config.Db.Exec(
		sqlStatement,
		id,
		item.ItemName,
		item.Description,
		item.Price,
		item.Stock,
		item.ModifiedBy,
		item.ModifiedAt,
	)
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

func DeleteItem(id int64) (int64, error) {
	sqlStatement := `DELETE from items WHERE id = $1`

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
