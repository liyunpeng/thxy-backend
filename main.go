package main

import (
	"os"
	"os/signal"
	"syscall"
	"thxy/logger"
	"thxy/model"
	"thxy/router"
	"thxy/setting"
)

func main() {
	logger.InitCustLogger()
	setting.InitConfig("conf/app.toml")
	r := router.InitRouter()
	model.Setup()

	//addr := fmt.Sprintf(":%d", setting.TomlConfig.Test.Server.Host)
	//maxHeaderBytes := 1 << 20

	go r.Run(":8082")


	//srv := &http.Server{
	//	Addr:    addr,
	//	Handler: r,
	//	//ReadTimeout:    time.Duration(conf.Server.ReadTimeout * 1e9),
	//	//WriteTimeout:   time.Duration(conf.Server.WriteTimeout * 1e9),
	//	MaxHeaderBytes: maxHeaderBytes,
	//}
	//
	//go func() {
	//	// service connections
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		logger.Info.Fatalf("listen: %s\n", err)
	//	}
	//}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
	)

Loop:
	for {
		select {
		case <-signals:
			logger.Info.Println("开始执行关闭服务")
			//r.Stop()


			break Loop
		}
	}

}
