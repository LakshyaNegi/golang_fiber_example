package main

import (
	"log"
	"strconv"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/jinzhu/gorm"
)

// var postCache PostCache

var rc *redisCache = NewRedisCache("localhost:6379", 0, 20)

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
	var rcshoe *Shoe = rc.Get(Paramid)
	if rcshoe == nil {
		log.Println("Shoe from DB")
		shoe := new(Shoe)
		db, err := gorm.Open("sqlite3", "shoes.db")
		if err != nil {
			panic("Failed to connect database")
		}

		db.Find(&shoe, id)
		rc.Set(Paramid, shoe)
		return c.Status(fiber.StatusOK).JSON(shoe)
	}

	return c.Status(fiber.StatusOK).JSON(rcshoe)
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
