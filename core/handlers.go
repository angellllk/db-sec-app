package core

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func registerBuyer(c *fiber.Ctx, db *sql.DB) error {
	response := httpResponse{
		Error:   true,
		Message: "Can't parse body data",
	}

	var data RegisterBuyer

	errParse := c.BodyParser(&data)
	if errParse != nil {
		log.Printf("got error: %v\n", errParse)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	record, errParseData := ParseData(data)
	if errParseData != nil {
		response.Message = errParseData.Error()
		log.Printf("got error: %v\n", errParseData)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	errAdd := AddBuyer(db, record)
	if errAdd != nil {
		response.Message = errAdd.Error()
		log.Printf("got error: %v\n", errAdd)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse{
		Error:   false,
		Message: fmt.Sprintf("successfully registered buyer %s", record.Name),
	})
}

func importCar() error {
	return nil
}

func createSale() error {
	return nil
}

func salesDetails(sale string) error {
	return nil
}
