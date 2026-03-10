package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/BurntSushi/toml"
)

type Website struct {
	Name     string  `toml:"name"`
	URL      string  `toml:"url"`
	Interval float64 `toml:"interval"`
}

type ServerConfig struct {
	Verbose bool `toml:"verbose"`
}

type Config struct {
	Websites []Website
	Server   ServerConfig
}

func loadConfig(path string) Config {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		log.Fatal("Could not read TOML config")
	}
	return conf
}

func validateWebsite(site *Website) {
	if site.Interval <= 0.0 {
		site.Interval = 5.0
	}
	if site.URL == "" || site.Name == "" {
		log.Fatal("URL or name not found for a website!")
	}
}

func main() {
	conf := loadConfig("Config.toml")
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
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			resp, err := http.Get(website.URL)
			if err != nil || resp.StatusCode != 200 {
				fmt.Printf("ALERT: %s is DOWN\n", website.Name)
			} else {
				if conf.Server.Verbose {
					fmt.Printf("%s is Okay\n", website.Name)
				}

			}
		}
	}
}
