package handler

import (
	"errors"
	"livecode-3-arvisy/model"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) ProductHandler {
	return ProductHandler{DB: db}
}

func (p *ProductHandler) GetAllProduct(c echo.Context) error {
	var products []model.Products
	result := p.DB.Find(&products)
	if result.Error != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to retrieve products",
			"detail":  result.Error.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"products": products,
	})
}

func (p *ProductHandler) GetProductByID(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid product ID",
		})
	}

	var product model.Products
	result := p.DB.First(&product, productID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(404, echo.Map{
				"message": "product not found",
			})
		}
		return c.JSON(500, echo.Map{
			"message": "failed to retrieve product",
		})
	}

	return c.JSON(200, echo.Map{
		"product": product,
	})
}
