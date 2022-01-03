package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ekudinov/hms-push-go/pkg/authention"
	hms_config "github.com/ekudinov/hms-push-go/pkg/config"
	"github.com/ekudinov/hms-push-go/pkg/core"
	"github.com/ekudinov/huawei-push-authorizator/pkg/config"
)

const (

	// below is public address
	// get token address
	authUrl = "https://login.cloud.huawei.com/oauth2/v2/token"
	// send push msg address
	pushUrl = "https://api.push.hicloud.com"
)

type Service struct {
	config *config.Config
	client *core.HttpPushClient
	ticker *time.Ticker
}

func NewService(appConfig *config.Config) *Service {

	clientConfig := &hms_config.Config{
		AppId:     appConfig.ClientID,
		AppSecret: appConfig.ClientSecret,
		AuthUrl:   authUrl,
		PushUrl:   pushUrl,
	}

	client := getPushClient(clientConfig)
	ticker := time.NewTicker(time.Duration(appConfig.CheckInterval) * time.Second)

	return &Service{
		config: appConfig,
		client: client,
		ticker: ticker,
	}
}

func getPushClient(conf *hms_config.Config) *core.HttpPushClient {
	client, err := core.NewHttpClient(conf)
	if err != nil {
		log.Printf("Failed to create new push client! Error is %s\n", err.Error())
		panic(err)
	}

	return client
}

func (s *Service) Start(ctx context.Context) {

	log.Println("Start service.")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-s.ticker.C:

				log.Printf("Start check token at %s.", t.Format("02.01.2006 15:04:05"))

				if isTokenExpired(s.config.EarlyUpdateTime, s.client.GetToken()) {

					log.Println("Token is not valid and try refresh it.")

					err := s.client.RefreshToken()

					if err != nil {
						log.Println("Error refresh token", err)
					}

					log.Println("End try refresh token.")
				}

			}
		}
	}()

}

func (s *Service) Stop(ctx context.Context) {

	s.ticker.Stop()

	ctx.Done()

	log.Println("Stop service.")
}

func (s *Service) GetValidToken() (string, error) {

	token := s.client.GetToken()

	if token.ExpiredAt.Before(time.Now()) {

		log.Println("Token is not valid and expired at", token.ExpiredAt.Format("02.01.2006 15:04:05"))

		return "", errors.New("Token is expired. Check service to update token!")
	}

	return token.Value, nil
}

func isTokenExpired(earlyUpdateTime int, token auth.Token) bool {
	return token.ExpiredAt.Before(time.Now().Add(time.Duration(earlyUpdateTime) * time.Second))
}
