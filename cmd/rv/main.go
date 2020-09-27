package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jwilner/rv/internal/platform"
)

func main() {
	var c platform.Config

	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	stdoutLog, _ := strconv.ParseBool(os.Getenv("STDOUT_LOG"))

	log.SetFlags(log.LstdFlags)
	if c.Debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if stdoutLog {
		log.SetOutput(os.Stdout)
	}

	c.DBURL = os.Getenv("DATABASE_URL")

	if port := os.Getenv("PORT"); port != "" {
		c.Addr = ":" + port
	}
	if grpcPort := os.Getenv("GRPC_PORT"); grpcPort != "" && !strings.Contains(grpcPort, ":") {
		c.GRPCAddr = ":" + grpcPort
	}

	c.StaticDir = os.Getenv("STATIC_DIR") // where to serve static assets from, if at all

	if c.SigningKey = os.Getenv("TOKEN_SIGNING_KEY"); c.SigningKey == "" {
		log.Fatalln("TOKEN_SIGNING_KEY must be provided")
	}

	{
		tokDurS, ok := os.LookupEnv("TOKEN_DURATION")
		if !ok {
			log.Fatalln("TOKEN_DURATION must be provided")
		}
		length, err := time.ParseDuration(tokDurS)
		if err != nil || length == 0 {
			log.Fatalf("Invalid TOKEN_DURATION %q: %v", tokDurS, err)
		}
		c.TokenLength = length
	}

	// if set, we'll serve the slack endpoints
	c.Slack.Token = os.Getenv("SLACK_TOKEN")
	c.Slack.SigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	if info, ok := os.LookupEnv("SLACK_CLIENT_INFO"); ok {
		clientInfo := strings.SplitN(info, ":", 2)
		if len(clientInfo) != 2 {
			log.Fatalf("Expected colon delimited CLIENT_ID:CLIENT_SECRET for SLACK_CLIENT_INFO")
		}
		c.Slack.ClientID, c.Slack.ClientSecret = clientInfo[0], clientInfo[1]
	}
	if c.Slack.Token != "" && (c.Slack.SigningSecret == "" || c.Slack.ClientID == "" || c.Slack.ClientSecret == "") {
		log.Fatal("All of SLACK_TOKEN, SLACK_SIGNING_SECRET, and SLACK_CLIENT_INFO are required to run slack component")
	}
	if err := platform.Run(&c); err != nil {
		log.Fatal(err)
	}
}
