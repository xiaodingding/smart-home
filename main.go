package main

import (
	"fmt"
	"github.com/e154/smart-home/api/mobile"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/websocket"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/initial"
	l "github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/stream_proxy"
	"github.com/op/go-logging"
	"os"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(shortVersionBanner, GetHumanVersion())
			return
		case "-backup":
			container := BuildContainer()
			container.Invoke(func(
				backup *backup.Backup,
				graceful *graceful_service.GracefulService) {

				if err := backup.New(); err != nil {
					log.Error(err.Error())
				}

				graceful.Shutdown()
			})
			return
		case "-restore":
			if len(os.Args) < 3 {
				log.Error("need backup name")
				return
			}
			container := BuildContainer()
			container.Invoke(func(
				backup *backup.Backup,
				graceful *graceful_service.GracefulService) {

				if err := backup.Restore(os.Args[2]); err != nil {
					log.Error(err.Error())
				}

				graceful.Shutdown()
			})
			return
		case "-reset":
			container := BuildContainer()
			container.Invoke(func(
				initialService *initial.InitialService) {

				initialService.Reset()
			})
			return
		default:
			fmt.Printf(verboseVersionBanner, "v2", os.Args[0])
			return
		}
	}

	start()
}

func start() {

	fmt.Printf(shortVersionBanner, "")

	container := BuildContainer()
	err := container.Invoke(func(server *server.Server,
		graceful *graceful_service.GracefulService,
		back *l.LogBackend,
		initialService *initial.InitialService,
		ws *websocket.WebSocket,
		streamProxy *stream_proxy.StreamProxy,
		mobileServer *mobile.MobileServer) {

		l.Initialize(back)
		go server.Start()
		go mobileServer.Start()
		go ws.Start()
		go streamProxy.Start()

		graceful.Wait()
	})

	if err != nil {
		panic(err.Error())
	}
}
