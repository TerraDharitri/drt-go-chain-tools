package main

import (
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
)

var (
	log            = logger.GetOrCreate("main")
	jsonMarshaller = &marshal.JsonMarshalizer{}
)
