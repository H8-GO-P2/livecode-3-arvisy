package handler

import (
	"livecode-3-arvisy/model"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("your-secret-key")

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return UserHandler{DB: db}
}

func (u *UserHandler) Register(c echo.Context) error {
	user := new(model.Users)

	err := c.Bind(user)
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid request",
		})
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		return c.JSON(400, echo.Map{
			"message": "name, email, and password are required",
		})
	}

	existingUser := new(model.Users)
	result := u.DB.Where("email = ?", user.Email).First(existingUser)
	if result.RowsAffected > 0 {
		return c.JSON(400, echo.Map{
			"message": "email already registered",
		})
	}

	var lastUserID int
	u.DB.Model(&model.Users{}).Select("user_id").Order("user_id desc").Limit(1).Scan(&lastUserID)

	user.UserID = lastUserID + 1

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to hash password",
			"detail":  err.Error(),
		})
	}

	user.Password = string(hashedPassword)

	token, err := GenerateJWTToken(user)
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to generate JWT token",
			"detail":  err.Error(),
		})
	}

	user.JWTToken = token

	result = u.DB.Create(&user)
	if result.Error != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to register user",
			"detail":  result.Error.Error(),
		})
	}

	responseBody := model.Users{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
	}

	return c.JSON(201, echo.Map{
		"message":   "success register",
		"user_info": responseBody,
	})
}

func (uh *UserHandler) Login(c echo.Context) error {
	loginRequest := new(model.Users)
	if err := c.Bind(loginRequest); err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid request",
		})
	}

	user := new(model.Users)
	result := uh.DB.Where("email = ?", loginRequest.Email).First(user)
	if result.Error != nil {
		return c.JSON(401, echo.Map{
			"message": "invalid email or password",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.JSON(401, echo.Map{
			"message": "invalid email or password",
		})
	}

	token, err := GenerateJWTToken(user)
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to generate JWT token",
			"detail":  err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "login success",
		"token":   token,
	})
}

func GenerateJWTToken(user *model.Users) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = user.UserID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (u *UserHandler) GetCart(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(401, echo.Map{
			"message": "unauthorized",
		})
	}

	var carts []model.Carts
	if err := u.DB.Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to retrieve cart data",
			"detail":  err.Error(),
		})
	}

	for i := range carts {
		carts[i].CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	return c.JSON(200, echo.Map{
		"carts": carts,
	})
}

func (u *UserHandler) CreateCart(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(401, echo.Map{
			"message": "unauthorized",
		})
	}

	var cartData model.Carts
	if err := c.Bind(&cartData); err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid request data",
		})
	}

	var lastCartID int
	u.DB.Model(&model.Carts{}).Select("cart_id").Order("cart_id desc").Limit(1).Scan(&lastCartID)

	cartData.CartID = lastCartID + 1

	cartData.UserID = userID
	cartData.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	result := u.DB.Create(&cartData)
	if result.Error != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to create cart",
			"detail":  result.Error.Error(),
		})
	}

	return c.JSON(201, echo.Map{
		"message": "cart created successfully",
		"cart":    cartData,
	})
}

func (u *UserHandler) DeleteCart(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(401, echo.Map{
			"message": "unauthorized",
		})
	}

	cartID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid cart ID",
		})
	}

	result := u.DB.Where("user_id = ? AND cart_id = ?", userID, cartID).Delete(&model.Carts{})
	if result.Error != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to delete cart",
			"detail":  result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(404, echo.Map{
			"message": "cart not found",
		})
	}

	return c.JSON(200, echo.Map{
		"message": "cart deleted successfully",
	})
}

func (u *UserHandler) GetOrder(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(401, echo.Map{
			"message": "unauthorized",
		})
	}

	var orders []model.Orders
	if err := u.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to retrieve order data",
			"detail":  err.Error(),
		})
	}

	for i := range orders {
		orders[i].CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	return c.JSON(200, echo.Map{
		"orders": orders,
	})
}

func (u *UserHandler) CreateOrder(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(401, echo.Map{
			"message": "unauthorized",
		})
	}

	var orderData model.Orders
	if err := c.Bind(&orderData); err != nil {
		return c.JSON(400, echo.Map{
			"message": "invalid request data",
		})
	}

	var totalCartPrice float64
	cartItems := []model.Carts{}
	u.DB.Where("user_id = ?", userID).Find(&cartItems)

	for _, cartItem := range cartItems {
		product := model.Products{}
		u.DB.Where("product_id = ?", cartItem.ProductID).First(&product)
		totalCartPrice += float64(cartItem.Quantity) * product.Price
	}

	orderData.TotalPrice = totalCartPrice

	var lastOrderID int
	u.DB.Model(&model.Orders{}).Select("order_id").Order("order_id desc").Limit(1).Scan(&lastOrderID)

	orderData.OrderID = lastOrderID + 1
	orderData.UserID = userID
	orderData.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	result := u.DB.Create(&orderData)
	if result.Error != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to create order",
			"detail":  result.Error.Error(),
		})
	}

	deleteCartResult := u.DB.Where("user_id = ?", userID).Delete(&model.Carts{})
	if deleteCartResult.Error != nil {
		return c.JSON(500, echo.Map{
			"message": "failed to delete user carts",
			"detail":  deleteCartResult.Error.Error(),
		})
	}

	return c.JSON(201, echo.Map{
		"message":     "order created successfully",
		"order":       orderData,
		"total_price": orderData.TotalPrice,
	})
}
