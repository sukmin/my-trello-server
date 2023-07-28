package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/skip"
)

func main() {

	server := fiber.New()
	server.Use(requestid.New())
	server.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	basicAuthHandler := basicauth.New(basicauth.Config{
		Users: map[string]string{
			"test": "test",
		},
	})
	server.Use(skip.New(basicAuthHandler, func(ctx *fiber.Ctx) bool {
		return ctx.Method() == fiber.MethodGet && ctx.OriginalURL() == "/health" //기본인증 제외정보
	}))

	server.Get("/health", func(c *fiber.Ctx) error {
		body := `{"health":"OK"}`
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(body)
	})

	server.Get("/test", func(c *fiber.Ctx) error {
		body := `{"test":"test"}`
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(body)
	})

	port := ":3000"
	fmt.Printf("server start fail. port:%s\n", port)
	err := server.Listen(port)
	if err != nil {
		fmt.Printf("server start fail. err:%s\n", err)
	}

}
