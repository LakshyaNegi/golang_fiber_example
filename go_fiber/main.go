package main

import (
	"log"
	"strconv"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DBConn *gorm.DB
)

type Shoe struct {
	gorm.Model
	Name  string  `json:"name,omitempty"`
	Size  int     `json:"size,omitempty"`
	Price float64 `json:"price,omitempty"`
	Email string  `json:"email,omitempty"`
}

type Users struct {
	gorm.Model
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// var shoes = []*Shoe{
// 	{Id: 1, Name: "Air max 97", Size: 10, Price: 16499, Email: "admin"},
// 	{Id: 2, Name: "Air force 1", Size: 9, Price: 8999, Email: "admin"},
// 	{Id: 3, Name: "Cortez", Size: 8, Price: 7999, Email: "admin2"},
// }

func initDB() {
	//var err error
	DBConn, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}
	DBConn2, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}

	log.Printf("Database connected")

	DBConn.AutoMigrate(&Shoe{})
	DBConn2.AutoMigrate(&Users{})
	log.Printf("Database migrated")
}

func main() {
	app := fiber.New()
	initDB()
	//defer DBConn.Close()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!!")
	})

	app.Post("/login", login)
	app.Post("/register", register)

	app.Get("/shoes", GetShoes)
	app.Get("/shoes/:id", GetShoe)
	app.Post("/shoes", Protected(), CreateShoe)
	app.Delete("/shoes/:id", Protected(), DeleteShoe)
	app.Put("/shoes/:id", Protected(), UpdateShoe)

	err := app.Listen(":4000")
	if err != nil {
		panic(err)
	}

}

func GetShoes(c *fiber.Ctx) error {
	db, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}
	var shoes []Shoe
	db.Find(&shoes)
	return c.Status(fiber.StatusOK).JSON(shoes)
}

func GetShoe(c *fiber.Ctx) error {
	Paramid := c.Params("id")
	id, _ := strconv.Atoi(Paramid)

	db, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}
	var shoe Shoe
	db.Find(&shoe, id)

	return c.Status(fiber.StatusOK).JSON(shoe)
}

func CreateShoe(c *fiber.Ctx) error {
	//db := DBConn
	db, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	shoe := new(Shoe)
	if err := c.BodyParser(shoe); err != nil {
		log.Panicf("Create shoe err %v", err)
		return c.Status(400).SendString(err.Error())
	}

	shoe.Email = email
	//shoe.Id = 123
	log.Println(shoe)

	db.Create(&shoe)

	return c.Status(fiber.StatusOK).JSON(shoe)

}

func DeleteShoe(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	ParamId := c.Params("id")
	id, _ := strconv.Atoi(ParamId)

	db, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}
	var shoe Shoe
	db.Find(&shoe, id)

	// if shoe.Name == "" {
	// 	return c.Status(fiber.StatusNotFound).SendString("Shoe not found")
	// }

	if shoe.Email == email {
		db.Delete(&shoe)
		return c.Status(fiber.StatusOK).JSON(shoe)
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Status(fiber.StatusNotFound).SendString("Shoe not found")
}

func UpdateShoe(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	Paramid := c.Params("id")
	id, _ := strconv.Atoi(Paramid)

	db, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}

	var shoe Shoe
	db.Find(&shoe, id)

	s := new(Shoe)
	if err := c.BodyParser(s); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if shoe.Email == email {
		db.Model(&shoe).Update("Price", s.Price)
		db.Model(&shoe).Update("Name", s.Name)
		db.Model(&shoe).Update("Size", s.Size)

		return c.Status(fiber.StatusOK).JSON(shoe)
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Status(fiber.StatusNotFound).SendString("Shoe not found")
}

func register(c *fiber.Ctx) error {
	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}
	creds := new(Users)
	if err := c.BodyParser(creds); err != nil {
		log.Printf("Parse error %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Println(creds)

	db.Create(&creds)

	return c.Status(fiber.StatusOK).SendString("Registered")

}

func login(c *fiber.Ctx) error {
	input := new(Users)
	if err := c.BodyParser(input); err != nil {
		log.Printf("Parse error %v", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	email := input.Email
	pass := input.Password

	//log.Println(input)
	//log.Println(input.Password)

	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}

	var user Users
	db.Where("Email = ?", email).First(&user)

	if pass != user.Password {
		log.Printf("Password error")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})

}

func Protected() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("secret"),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		log.Printf("Missing jwt error %v", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		log.Printf("Unathorized error %v", err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
