package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/skip"
)

type LoginInput struct {
	Email string `json:"email"`
}

func main() {

	server := fiber.New()
	server.Use(requestid.New())
	server.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	basicAuthHandler := basicauth.New(basicauth.Config{
		Users: map[string]string{
			"kwonsukmin@naver.com": "1234",
		},
	})
	server.Use(skip.New(basicAuthHandler, func(ctx *fiber.Ctx) bool {
		method := ctx.Method()
		path := ctx.OriginalURL()
		if path == "/login" && method == fiber.MethodPost { //로그인
			return true
		}

		return ctx.Method() == fiber.MethodGet && path == "/health" //기본인증 제외정보
		//return true
	}))

	server.Get("/health", func(c *fiber.Ctx) error {
		body := `{"health":"OK"}`
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(body)
	})

	server.Get("/boards", func(c *fiber.Ctx) error {
		body := `{"boards":[
{
	"id":"done",
	"name":"돈냉면"
},
{
	"id":"susi",
	"name":"스시"
}
]}`
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(body)
	})

	server.Post("/login", func(c *fiber.Ctx) error {
		loginInput := new(LoginInput)
		if err := c.BodyParser(loginInput); err != nil {
			return err
		}

		//여기서 정보를 비교한다고 치고

		body := `
{"accessToken":"a3dvbnN1a21pbkBuYXZlci5jb206MTIzNA=="}
`
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(body)
	})

	server.Get("/test", func(c *fiber.Ctx) error {
		body := `{"test":"test"}`
		c.Set("content-type", "application/json; charset=utf-8")
		return c.SendString(body)
	})

	server.Get("/metrics", monitor.New())

	port := ":3000"
	fmt.Printf("server start fail. port:%s\n", port)
	err := server.Listen(port)
	if err != nil {
		fmt.Printf("server start fail. err:%s\n", err)
	}

}
