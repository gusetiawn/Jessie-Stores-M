package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/handler"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/middleware/authtoken"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server func(dependency *ServerHandlers) error

type ServerHandlers struct {
	handler.Handler
	*authtoken.AuthTokenMiddleware
}

func Serve(dependency *ServerHandlers) error {
	file, err := os.Open("./panic.log")
	if err != nil {
		logrus.Info(err)
		file, err = os.Create("./panic.log")
		if err != nil {
			logrus.Error(err)
		}
	}

	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.RecoveryWithWriter(file))
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/metrics"}}))

	settingCors := cors.DefaultConfig()
	settingCors.AllowOrigins = viper.GetStringSlice("CORS_ALLOWED")
	settingCors.AllowCredentials = true
	settingCors.AddAllowHeaders(viper.GetStringSlice("HEADERS_ALLOWED")...)
	r.Use(cors.New(settingCors))

	setupRouter(r, dependency)
	port := fmt.Sprintf(":%s", viper.GetString("APP_PORT"))

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	q := <-quit
	logrus.Info("Shutdown server ...")
	logrus.Info("Exit signal: ", q)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server shutdown: ", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()

	return err
}
