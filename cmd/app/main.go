package main

import (
	"github.com/amiosamu/vk-internship/internal/app"
)

// @title Marketplace Backend API
// @version 1.0.0
// @description Test assignment from VK for a Backend Developer position
// @host localhost:8000
// @BasePath /

const configPath = "config/config.yaml"

func main() {

	app.Run(configPath)
}
