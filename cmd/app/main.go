package main

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"seg-red-file/internal/app/config"
)

func init() {
	_ = godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")
	certs := os.Getenv("CERTS_FOLDER")
	app := config.SetupRouter()
	err := app.RunTLS(":"+port, filepath.Join(certs, "file.crt"), filepath.Join(certs, "file.key"))
	if err != nil {
		panic(err)
	}
}
