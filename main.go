package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github/xmapst/free-ss.site/config"
	_ "github/xmapst/free-ss.site/config"
	"github/xmapst/free-ss.site/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Free-SS.site API
// @version 1.0
// @description This is a Free-SS API, contain interfaces such as SS/SSR.
func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	gin.SetMode(config.Config.RunMode)
	router := routers.InitRouter()
	readTimeout := time.Duration(config.Config.ReadTimeout) * time.Second
	writeTimeout := time.Duration(config.Config.WriteTimeout) * time.Second
	endPoint := fmt.Sprintf("%s:%d", config.Config.HTTPAddr, config.Config.HTTPPort)
	maxHeaderBytes := 1 << 33

	server := &http.Server{
		Addr:           endPoint,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logrus.Info("Server Staring...")
	go func() {
		logrus.Info("Listen: ", endPoint)
		if err := server.ListenAndServe(); err != nil {
			logrus.Warningln(err)
		}
	}()

	<-c
	logrus.Warningln("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	os.Exit(0)
}
