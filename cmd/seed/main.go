package main

import (
	"context"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yorukot/knocker/db"
	"github.com/yorukot/knocker/helpers/config"
	"github.com/yorukot/knocker/helpers/logger"
	"go.uber.org/zap"
)

var sampleURLs = []string{
	"https://api.github.com",
	"https://www.google.com",
	"https://www.reddit.com",
	"https://www.amazon.com",
	"https://www.twitter.com",
	"https://www.facebook.com",
	"https://www.youtube.com",
	"https://www.wikipedia.org",
	"https://www.linkedin.com",
	"https://www.instagram.com",
	"https://www.netflix.com",
	"https://www.twitch.tv",
	"https://www.discord.com",
	"https://www.stackoverflow.com",
	"https://www.github.com",
	"https://www.gitlab.com",
	"https://www.bitbucket.org",
	"https://www.npmjs.com",
	"https://www.pypi.org",
	"https://www.docker.com",
	"https://www.stripe.com",
	"https://www.shopify.com",
	"https://www.cloudflare.com",
	"https://www.digitalocean.com",
	"https://www.heroku.com",
	"https://www.vercel.com",
	"https://www.netlify.com",
	"https://www.aws.amazon.com",
	"https://cloud.google.com",
	"https://azure.microsoft.com",
	"https://www.mongodb.com",
	"https://www.postgresql.org",
	"https://www.redis.io",
	"https://www.elastic.co",
	"https://www.jenkins.io",
	"https://www.travis-ci.com",
	"https://circleci.com",
	"https://www.atlassian.com",
	"https://www.slack.com",
	"https://www.zoom.us",
	"https://www.notion.so",
	"https://www.figma.com",
	"https://www.canva.com",
	"https://www.adobe.com",
	"https://www.apple.com",
	"https://www.microsoft.com",
	"https://www.oracle.com",
	"https://www.ibm.com",
	"https://www.salesforce.com",
	"https://www.dropbox.com",
	"https://www.box.com",
	"https://www.trello.com",
	"https://www.asana.com",
	"https://www.jira.atlassian.com",
	"https://www.confluence.atlassian.com",
	"https://www.medium.com",
	"https://www.dev.to",
	"https://www.hashnode.com",
	"https://www.producthunt.com",
	"https://www.hackernews.com",
	"https://www.techcrunch.com",
	"https://www.theverge.com",
	"https://www.wired.com",
	"https://www.cnet.com",
	"https://www.engadget.com",
	"https://www.arstechnica.com",
	"https://www.zdnet.com",
	"https://www.coinbase.com",
	"https://www.binance.com",
	"https://www.kraken.com",
	"https://www.paypal.com",
	"https://www.square.com",
	"https://www.revolut.com",
	"https://www.n26.com",
	"https://www.spotify.com",
	"https://www.soundcloud.com",
	"https://www.bandcamp.com",
	"https://www.deezer.com",
	"https://www.tidal.com",
	"https://www.pandora.com",
	"https://www.hulu.com",
	"https://www.disneyplus.com",
	"https://www.hbomax.com",
	"https://www.crunchyroll.com",
	"https://www.paramount.com",
	"https://www.peacocktv.com",
	"https://www.espn.com",
	"https://www.nba.com",
	"https://www.nfl.com",
	"https://www.mlb.com",
	"https://www.nhl.com",
	"https://www.fifa.com",
	"https://www.uefa.com",
	"https://www.weather.com",
	"https://www.accuweather.com",
	"https://www.bbc.com",
	"https://www.cnn.com",
	"https://www.nytimes.com",
}

var intervals = []int32{30, 60, 120} // 30s, 1m, 2m, 5m, 10m

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load environment variables
	_, err := config.InitConfig()
	if err != nil {
		zap.L().Fatal("Error initializing config", zap.Error(err))
	}

	// Initialize database
	pgsql, err := db.InitDatabase()
	if err != nil {
		zap.L().Fatal("Error initializing Postgres", zap.Error(err))
	}
	defer pgsql.Close()

	ctx := context.Background()
	now := time.Now()

	zap.L().Info("Starting to insert monitors...")

	monitorsInserted := 0
	monitorID := int64(1)

	// Create 20 monitors for testing
	for i := 0; i < 20; i++ {
		url := sampleURLs[i%len(sampleURLs)]
		interval := intervals[i%len(intervals)]

		// Set next_check to now so they get picked up immediately
		lastCheck := now.Add(-time.Duration(interval) * time.Second)
		nextCheck := now

		query := `
			INSERT INTO monitors (id, url, interval, last_check, next_check)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING
		`

		_, err := pgsql.Exec(ctx, query, monitorID, url, interval, lastCheck, nextCheck)
		if err != nil {
			zap.L().Error("Failed to insert monitor",
				zap.Int64("id", monitorID),
				zap.String("url", url),
				zap.Error(err))
			continue
		}

		monitorsInserted++
		monitorID++
	}

	zap.L().Info("Successfully inserted monitors", zap.Int("count", monitorsInserted))

	// Show stats
	var total int64
	err = pgsql.QueryRow(ctx, "SELECT COUNT(*) FROM monitors").Scan(&total)
	if err != nil {
		zap.L().Error("Failed to count monitors", zap.Error(err))
	} else {
		zap.L().Info("Total monitors in database", zap.Int64("total", total))
	}

	var ready int64
	err = pgsql.QueryRow(ctx, "SELECT COUNT(*) FROM monitors WHERE next_check <= $1", now).Scan(&ready)
	if err != nil {
		zap.L().Error("Failed to count ready monitors", zap.Error(err))
	} else {
		zap.L().Info("Monitors ready to be checked", zap.Int64("ready", ready))
	}
}
