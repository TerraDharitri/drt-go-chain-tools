package main

import (
	logger "github.com/TerraDharitri/drt-go-chain-logger"
)

var (
	log             = logger.GetOrCreate("main")
	outputFileName  = "output.json"
	outputFilePerms = 0644
)
