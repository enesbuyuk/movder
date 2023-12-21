package main

import (
	websocketmovder "github.com/enesbuyuk/movder/pkg/websocket"
	"github.com/enesbuyuk/movder/routers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	routers.StartServer()
	go websocketmovder.H.Run()
}
