package gcplogger

import (
	"context"
	"errors"
	"fmt"
	"github.com/harness/runner/delegateshell/client"
	"github.com/patrickmn/go-cache"
	"github/meenaravichandran1/runner-logger/logger"
	"golang.org/x/oauth2"
	"time"
)

var NoTokenError = errors.New("no log token defined on the server")
var ExpiredError = errors.New("expired token")

const (
	defaultCacheExpirationDuration = 30 * time.Minute
	expiryMargin                   = 2 * time.Minute
	tokenKey                       = "loggingToken"
)

type TokenManager struct {
	// holds the oauth2 token
	cache         *cache.Cache
	ManagerClient *client.ManagerClient
	ProjectID     string
}

func NewTokenManager(ctx context.Context, managerClient *client.ManagerClient) (*TokenManager, error) {
	tokenCache := cache.New(defaultCacheExpirationDuration, -1)
	tokenManager := &TokenManager{
		cache:         tokenCache,
		ManagerClient: managerClient,
	}
	_, err := tokenManager.SetToken(ctx)
	if err != nil {
		return nil, err
	}
	return tokenManager, nil
}

func (tokenManager *TokenManager) Token() (*oauth2.Token, error) {
	token, found := tokenManager.cache.Get(tokenKey)
	if found {
		if token, ok := token.(*oauth2.Token); ok {
			return token, nil
		}
		// If typecast fails, log an error and proceed to refresh
		logger.Errorln("Invalid token type in cache")
	}
	// TODO add retries
	logger.Infoln("refreshing logging token")
	token, err := tokenManager.SetToken(context.Background())
	if err != nil {
		logger.Errorln("cannot refresh logging token:", err)
		return nil, err
	}
	if token, ok := token.(*oauth2.Token); ok {
		return token, nil
	}
	return nil, fmt.Errorf("invalid token type in cache after refresh")
}

func (tokenManager *TokenManager) SetToken(ctx context.Context) (*oauth2.Token, error) {
	logCredentials, err := tokenManager.fetchLoggingToken(ctx)
	if err != nil {
		return nil, err
	}
	token := mapToOauthToken(logCredentials)

	if token.AccessToken == "" {
		return nil, NoTokenError
	}
	tokenManager.ProjectID = logCredentials.ProjectId
	durationUntilExpiration := time.Duration(logCredentials.ExpirationTimeMillis-time.Now().UnixMilli()-expiryMargin.Milliseconds()) * time.Millisecond
	tokenManager.cache.Set(tokenKey, token, durationUntilExpiration)
	logger.Printf("Logging token set for: %v", durationUntilExpiration)
	return token, nil
}

func (tokenManager *TokenManager) fetchLoggingToken(ctx context.Context) (*client.AccessTokenBean, error) {
	logCredentials, err := tokenManager.ManagerClient.GetLoggingToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logger credentials from the server: %w", err)
	}

	if logCredentials.ProjectId == "" || len(logCredentials.TokenValue) == 0 {
		return nil, fmt.Errorf("logging credentials are missing from the server")
	}

	if time.UnixMilli(logCredentials.ExpirationTimeMillis).Before(time.Now()) {
		return nil, ExpiredError
	}

	return logCredentials, nil
}

func mapToOauthToken(credentials *client.AccessTokenBean) *oauth2.Token {
	token := &oauth2.Token{
		AccessToken: credentials.TokenValue,
		Expiry:      time.UnixMilli(credentials.ExpirationTimeMillis),
	}
	return token
}
