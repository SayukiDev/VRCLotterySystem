package main

import (
	"flag"
	"time"

	"github.com/SayukiDev/VRCLotterySystem/config"
	"github.com/SayukiDev/VRCLotterySystem/internal/http"
	"github.com/SayukiDev/VRCLotterySystem/internal/provider"
	"github.com/SayukiDev/VRCLotterySystem/internal/task"
	"github.com/SayukiDev/VRCLotterySystem/log"

	"go.uber.org/zap"
)

var configPath = flag.String("c", "config.json", "config file path")

// Version will Replace with git hash or git tag
var Version = "v0.0.1"

func printHelloMessage() {
	log.Info("Dawing System")
	log.Info("Version: " + Version)
	log.Info("Author: https://github.com/SayukiDev")
	log.Info("Repository: SayukiDev/DrawingSystem")
}

func main() {
	flag.Parse()
	printHelloMessage()
	c, err := config.LoadConfig(*configPath)
	if err != nil {
		log.ErrorE("Loading config file failed", zap.Error(err))
	}

	p := provider.NewProvider(c)
	err = p.Init()
	if err != nil {
		log.ErrorE("Init provider failed", zap.Error(err))
	}

	task.NewTask(p, time.Minute*2).Start()

	h := http.NewHttp(c, p)
	err = h.Start(":8080")
	if err != nil {
		log.ErrorE("Start http server failed", zap.Error(err))
	}

	// Todo: start discord bot
}
