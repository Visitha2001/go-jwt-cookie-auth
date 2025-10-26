package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/visitha2001/go-jwt-auth/database"
	"github.com/visitha2001/go-jwt-auth/models"
	"gorm.io/gorm"
)

type ItemInput struct {
	UserID uint    `json:"user_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func GetItems(c *fiber.Ctx) error {
	var items []models.Item
	if err := database.DB.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not fetch items",
		})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}

func GetItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Item ID is required",
		})
	}

	var item models.Item
	err := database.DB.First(&item, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error", "message": "Item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not fetch item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

func CreateItem(c *fiber.Ctx) error {
	var input ItemInput
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid request body",
		})
	}
	if input.Name == "" || input.Price <= 0 || input.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid input data",
		})
	}
	newItem := models.Item{
		UserID: input.UserID,
		Name:   input.Name,
		Price:  input.Price,
	}
	if err := database.DB.Create(&newItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not create item",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "message": "Item created successfully",
		"item": newItem,
	})
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Item ID is required",
		})
	}
	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error", "message": "Item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not fetch item",
		})
	}
	var input ItemInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid request body",
		})
	}
	item.Name = input.Name
	item.Price = input.Price
	item.UserID = input.UserID
	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not update item",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", "message": "Item updated successfully",
		"item": item,
	})
}

func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Item ID is required",
		})
	}
	result := database.DB.Delete(&models.Item{}, id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not delete item",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "Item not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", "message": "Item deleted successfully",
		"item": id,
	})
}

func GetItemsByUserID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "User ID is required",
		})
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error", "message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not fetch user",
		})
	}
	var items []models.Item
	if err := database.DB.Where("user_id = ?", id).Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not fetch items",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Items fetched successfully",
		"user_id":   user.ID,
		"user_name": user.Username,
		"items":     items,
	})
}
