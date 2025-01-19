package main

import (
	"fit-byte/internal/config"
	"fit-byte/pkg/dotenv"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	env, err := dotenv.LoadEnv()
	if err != nil {
		log.Fatal("failed to load env", err.Error())
		return
	}

	log := config.NewLogger()
	validator := config.NewValidator()
	app := echo.New()
	s3Uploader := config.NewS3Uploader(env)
	pg := config.NewDatabase(log)
	defer pg.Pool.Close()

	config.Bootstrap(&config.BootstrapConfig{
		App:        app,
		DB:         pg,
		Log:        log,
		Validator:  validator,
		S3Uploader: s3Uploader,
		Env:        env,
	})

	PORT := os.Getenv("PORT")
	log.Fatal(app.Start(PORT))
}
