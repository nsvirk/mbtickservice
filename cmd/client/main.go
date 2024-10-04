// client is a simple client to test the server
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

type ClientConfig struct {
	KiteUserID     string
	KitePassword   string
	KiteTotpSecret string
	RedisHost      string
	RedisPort      string
	RedisPassword  string
}

func ClientConfigLoad() (*ClientConfig, error) {

	config := &ClientConfig{
		KiteUserID:     getEnv("KITE_USER_ID", ""),
		KitePassword:   getEnv("KITE_PASSWORD", ""),
		KiteTotpSecret: getEnv("KITE_TOTP_SECRET", ""),
		RedisHost:      getEnv("TS_CLIENT_REDIS_HOST", ""),
		RedisPort:      getEnv("TS_CLIENT_REDIS_PORT", ""),
		RedisPassword:  getEnv("TS_CLIENT_REDIS_PASSWORD", ""),
	}

	if config.KiteUserID == "" {
		return nil, fmt.Errorf("KITE_USER_ID is required")
	}

	if config.KitePassword == "" {
		return nil, fmt.Errorf("KITE_PASSWORD is required")
	}

	if config.KiteTotpSecret == "" {
		return nil, fmt.Errorf("KITE_TOTP_SECRET is required")
	}

	if config.RedisHost == "" {
		return nil, fmt.Errorf("TS_CLIENT_REDIS_HOST is required")
	}

	if config.RedisPort == "" {
		return nil, fmt.Errorf("TS_CLIENT_REDIS_PORT is required")
	}

	if config.RedisPassword == "" {
		return nil, fmt.Errorf("TS_CLIENT_REDIS_PASSWORD is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(cfg *ClientConfig) (*RedisClient, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{rdb: rdb}, nil
}

type UserSessionResponse struct {
	Status string `json:"status"`
	Data   struct {
		UserID        string `json:"user_id"`
		UserName      string `json:"user_name"`
		UserShortname string `json:"user_shortname"`
		AvatarURL     string `json:"avatar_url"`
		PublicToken   string `json:"public_token"`
		KFSession     string `json:"kf_session"`
		Enctoken      string `json:"enctoken"`
		LoginTime     string `json:"login_time"`
	} `json:"data"`
}

// getUserSessionResponse makes a POST request to the moneybots api to get the user_id and enctoken
func getUserSessionResponse(cfg *ClientConfig) (*UserSessionResponse, error) {
	apiURL := "https://www.moneybots.app/api/session/login"
	requestBody := map[string]string{
		"user_id":     cfg.KiteUserID,
		"password":    cfg.KitePassword,
		"totp_secret": cfg.KiteTotpSecret,
	}

	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make post request: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userSessionResponse UserSessionResponse
	if err := json.Unmarshal(bodyBytes, &userSessionResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if userSessionResponse.Status != "ok" {
		return nil, fmt.Errorf("API returned non-ok status: %s", userSessionResponse.Status)
	}

	return &userSessionResponse, nil
}

// Make a POST request to
// POST https://www.moneybots.app/ticks/publish/start
// with the following body
// {
//     "bot_id" :  "BOT1",
//     "ticker_instruments": [
//         "NSE:RELIANCE",
//         "MCX:GOLDM24DECFUT",
//         "MCX:GOLDM24NOVFUT",
//         "MCX:SILVERMIC25APRFUT",
//         "MCX:SILVERM25JUNFUT",
//         "MCX:GOLDM24OCTFUT"
//     ]
// }

// Response
// {
//     "status": "ok",
//     "data": {
//         "published_channel": "CH:TICKS:SA0846:BOT1",
//         "subscribed_count": 6
//     }
// }

type PublishStartRequest struct {
	BotID             string   `json:"bot_id"`
	TickerInstruments []string `json:"ticker_instruments"`
}

type PublishStartResponse struct {
	Status string `json:"status"`
	Data   struct {
		PublishedChannel string `json:"published_channel"`
		SubscribedCount  int    `json:"subscribed_count"`
	} `json:"data"`
}

func startPublishing(userID, enctoken string) (*PublishStartResponse, error) {
	apiURL := "https://www.moneybots.app/ticks/publish/start"
	requestBody := PublishStartRequest{
		BotID:             "BOT1",
		TickerInstruments: []string{"NSE:RELIANCE", "MCX:GOLDM24DECFUT", "MCX:GOLDM24NOVFUT", "MCX:SILVERMIC25APRFUT", "MCX:SILVERM25JUNFUT", "MCX:GOLDM24OCTFUT"},
	}

	// Marshal the request body
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	fmt.Println("Request body:", string(jsonRequestBody))

	// Create a new request with the Authorization header
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	authHeader := fmt.Sprintf("%s:%s", userID, enctoken)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make post request: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var publishStartResponse PublishStartResponse
	if err := json.Unmarshal(bodyBytes, &publishStartResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if publishStartResponse.Status != "ok" {
		return nil, fmt.Errorf("API returned non-ok status with code: %d", resp.StatusCode)
	}

	return &publishStartResponse, nil
}

func main() {
	fmt.Println("----------------------------------------")
	fmt.Println("Client started")
	fmt.Println("----------------------------------------")

	// Load config
	cfg, err := ClientConfigLoad()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	// Get user_id and enctoken
	userSessionResponse, err := getUserSessionResponse(cfg)
	if err != nil {
		fmt.Println("Failed to get user session:", err)
		return

	}

	// Get user_id and enctoken from userSession
	userID := userSessionResponse.Data.UserID
	enctoken := userSessionResponse.Data.Enctoken
	loginTime := userSessionResponse.Data.LoginTime

	fmt.Println("User ID:", userID)
	fmt.Println("Enctoken:", enctoken)
	fmt.Println("LoginTime:", loginTime)
	fmt.Println("----------------------------------------")

	// start publishing
	publishStartResponse, err := startPublishing(userID, enctoken)
	if err != nil {
		fmt.Println("Failed to start publishing:", err)
		return
	}

	fmt.Println("Published channel:", publishStartResponse.Data.PublishedChannel)
	fmt.Println("Subscribed count:", publishStartResponse.Data.SubscribedCount)
	fmt.Println("----------------------------------------")

	// Create Redis client
	redisClient, err := NewRedisClient(cfg)
	if err != nil {
		fmt.Println("Failed to create Redis client:", err)
		return
	}

	// Subscribe to the published channel
	pubsub := redisClient.rdb.Subscribe(context.Background(), publishStartResponse.Data.PublishedChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println("Received message:", msg.Payload)
		fmt.Println("----------------------------------------")
	}

}
