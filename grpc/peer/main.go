package main

import (
	"time"

	server "peer/server"
)

func main() {
	master := &server.Master{}
	go master.Start("10008")
	time.Sleep(1 * time.Second)

	worker := &server.Worker{}
	worker.Start("127.0.0.1:10008")
	defer worker.Destroy()
}
