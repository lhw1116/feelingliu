package main

//  程序启动入口
import (
	"context"
	"feelingliu/middleware"
	"feelingliu/modles"
	"feelingliu/routers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gin.SetMode(modles.ServerInfo.RunMode)
	r := routers.InitRouter()

	defer middleware.CloseLogFile()

	s := &http.Server{
		Addr:           modles.ServerInfo.ServerAddr,
		Handler:        r,
		ReadTimeout:    modles.ServerInfo.ReadTimeout,
		WriteTimeout:   modles.ServerInfo.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if e := s.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", e)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}