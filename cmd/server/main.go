package main

import (
	"context"
	"fmt"
	"net/http"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ekudinov/huawei-push-authorizator/pkg/config"
	"github.com/ekudinov/huawei-push-authorizator/pkg/service"

	"time"

	"github.com/gorilla/mux"
)

var (
	appService *service.Service

	appConfig *config.Config
)

func main() {

	appConfig = config.NewConfigFromEnv()

	appService = service.NewService(appConfig)

	ctx := context.Background()

	appService.Start(ctx)

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	router := mux.NewRouter()
	router.HandleFunc("/auth", auth).Methods("GET")

	serverAddr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)
		}
	}()

	log.Printf("Config host %s, check interval %d, early update time:%d", serverAddr, appConfig.CheckInterval, appConfig.EarlyUpdateTime)
	log.Print("Service and server started...")

	<-term

	appService.Stop(ctx)

	log.Print("Server start to stop...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server stop error:%+v", err)
	}

	log.Print("Server stop ok.")

}

func auth(resp http.ResponseWriter, req *http.Request) {
	if isApiTokenValid(req) {
		accessToken, err := appService.GetValidToken()
		if err != nil {
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte(err.Error()))

		} else {
			resp.Write([]byte(accessToken))
		}
	} else {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("Access deny\n"))
	}

}

func isApiTokenValid(req *http.Request) bool {

	token := req.Header.Get("x-api-token")

	return token == appConfig.ApiToken

}
