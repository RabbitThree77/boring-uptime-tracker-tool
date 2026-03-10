package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/joho/godotenv/autoload"
)

type Website struct {
	Name     string  `toml:"name"`
	URL      string  `toml:"url"`
	Interval float64 `toml:"interval"`
	Retry    uint    `toml:"retry"` //0 is for no retry, 1 is for 2 requests before error etc...
}

type ServerConfig struct {
	Verbose bool    `toml:"verbose"`
	Timeout float64 `toml:"timeout"`
}

type NotificationsConfig struct {
	DiscordWebhook string `toml:"discord_webhook"`
}

type Config struct {
	Websites      []Website
	Server        ServerConfig
	Notifications NotificationsConfig
}

func loadConfig(path string) Config {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		log.Fatal("Could not read TOML config")
	}
	return conf
}

func validateServer(conf *Config) {
	if conf.Server.Timeout <= 0 {
		conf.Server.Timeout = 1
	}
	conf.Notifications.DiscordWebhook = os.Getenv(conf.Notifications.DiscordWebhook)
	if conf.Notifications.DiscordWebhook == "" {
		fmt.Println("Discord Webhook not found!")
	}
}

func validateWebsite(site *Website) {
	if site.Interval <= 0.0 {
		site.Interval = 15.0
	}
	if site.URL == "" || site.Name == "" {
		log.Fatal("URL or name not found for a website!")
	}
}

func main() {
	conf := loadConfig("Config.toml")
	validateServer(&conf)
	for _, site := range conf.Websites {
		validateWebsite(&site)
		c := make(chan bool)
		go handleCheck(site, c, conf)
		defer close(c)
	}
	select {}
}

func handleCheck(website Website, done chan bool, conf Config) {
	ticker := time.NewTicker(time.Duration(float64(time.Second) * website.Interval))
	defer ticker.Stop()

	client := &http.Client{
		Timeout: time.Duration(float64(time.Second) * conf.Server.Timeout),
	}

	var attempts uint = 0
	var isDown bool = false


	for {
		select {
		case <-done:
			return
		case <-ticker.C:

			success := func() bool {
				resp, err := client.Get(website.URL)
				if err != nil {
					return false
				}
				defer resp.Body.Close()

				return resp.StatusCode == 200
			}()

			if !success {
				attempts++
				if conf.Server.Verbose {
					fmt.Println("an attempt failed")
				}
			} else {
				if conf.Server.Verbose {
					fmt.Printf("%s is Okay\n", website.Name)
				}
				attempts = 0
				if isDown {
					//TODO
					go sendDiscordEmbed(conf.Notifications.DiscordWebhook, website.Name, "Website is Operational!", 5763719)
					isDown = false
				}
			}
			if attempts >= website.Retry+1 {
				fmt.Printf("ALERT: %s is DOWN\n", website.Name)
				attempts = 0
				if !isDown {
					//TODO
					go sendDiscordEmbed(conf.Notifications.DiscordWebhook, website.Name, "Website is Down!", 15548997)
					isDown = true
				}
			}

		}
	}
}

type DiscordEmbed struct {
    Title string `json:"title"`
	Message string `json:"description"`
	Color int `json:"color"`
}

type DiscordPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

func sendDiscordEmbed(webhookUrl string, name string, message string, color int) {
	if webhookUrl == "" {
		return
	}

	embed := DiscordEmbed{
		Title: name,
		Message: message,
		Color: color,
	}
	list := [1]DiscordEmbed{embed}
	payload := DiscordPayload{
		Embeds: list[:],
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error Marhsaling: %v\n", err)
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("failed to send alert to discord: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
        fmt.Printf("Discord returned non-ok status: %d\n", resp.StatusCode)
    }

}
