package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var env map[string]string

func main() {
	loadEnv()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("static/index.html")
	})

	app.Post("/save-mc-progress", func(c *fiber.Ctx) error {
		log(c.FormValue("name"), "save progress minecraft")

		client, err := goph.NewConn(&goph.Config{
			User:     env["SSH_USER"],
			Addr:     env["SSH_HOST"],
			Auth:     goph.Password(env["SSH_PASSWORD"]),
			Port:     2233,
			Callback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			log(c.FormValue("name"), "connect via ssh failed: "+err.Error())
			return c.SendString("connect via ssh failed: " + err.Error())
		} else {
			cmd := fmt.Sprintf("cd /wri/mengkrep && sudo ./progress.sh '%s' '%s'", c.FormValue("name"), time.Now().Format("2006-01-02|15:04:05"))
			out, err := client.Run(cmd)
			if err != nil {
				log(c.FormValue("name"), "run command failed: "+err.Error())
				log(c.FormValue("name"), "command output: "+string(out))
				return c.SendString("run command failed: " + err.Error())
			} else {
				log(c.FormValue("name"), "run command success")
				log(c.FormValue("name"), "command output: "+string(out))
				return c.SendString("success")
			}
		}

		return c.SendString("success")
	})

	app.Listen(":8080")
}

func log(name string, message string) {
	fmt.Println(time.Now().Format("2006-01-02|15:04:05"), "\t", name, "\t", message)
}

func loadEnv() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		panic("cannot read .env file :" + err.Error())
	}

	env = envFile
}
