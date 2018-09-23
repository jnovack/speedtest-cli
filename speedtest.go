package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jnovack/speedtest"
)

func version() {
	fmt.Print(speedtest.Version)
	os.Exit(0)
}

func usage() {
	fmt.Fprintf(os.Stderr, "speedtest %s\n", speedtest.Version)
	fmt.Fprint(os.Stderr, "Command line interface for testing internet bandwidth using speedtest.net.\n\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var downloadSpeed, uploadSpeed int

	var hostname, payload string
	hostname, err := os.Hostname()
	var id = flag.String("id", hostname, "The id for this host (e.g. hostname)")

	var host = flag.String("host", "", "Metric server ip address")
	var port = flag.String("port", "80", "Metric server port")
	httpClient := &http.Client{}
	targetURL := fmt.Sprintf("http://%s:%s/public/metrics", *host, *port)

	opts := speedtest.ParseOpts()
	flag.Parse()

	switch {
	case opts.Help:
		usage()
		return
	case opts.Version:
		version()
		return
	}

	client := speedtest.NewClient(opts)

	if opts.List {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatalf("Failed to load server list: %v\n", err)
		}
		fmt.Println(servers)
		return
	}

	config, err := client.Config()
	if err != nil {
		log.Fatal(err)
	}

	client.Log("Testing from %s (%s)...please be patient.\n", config.Client.IP, config.Client.ISP)

	me := fmt.Sprintf(`{ "hostname": "%s", "ip": "%s", "isp": "%s", "latitude": %.4f, "longitude": %.4f, "country": "%s" }`, *id, config.Client.IP, config.Client.ISP, config.Client.Coordinates.Latitude, config.Client.Coordinates.Longitude, config.Client.Country)

	server := client.SelectServer(opts)

	log.Printf("Hosted by %s (%s) [%.2f km]: %d ms\n",
		server.Sponsor,
		server.Name,
		server.Distance,
		server.Latency/time.Millisecond)

	if *host != "" {
		payload = fmt.Sprintf(`{ "metric": { "name": "%s", "value": %d, "units": "ms" }, "client": %s, "server": %s}`, "latency", server.Latency/time.Millisecond, me, server.JSON())
		post(*httpClient, targetURL, payload)
	}

	if opts.Download {
		downloadSpeed = server.DownloadSpeed()
		reportSpeed(opts, "Download", downloadSpeed)
		if *host != "" {
			payload = fmt.Sprintf(`{ "metric": { "name": "%s", "value": %.4f, "units": "Mb" }, "client": "%s", "server": %s}`, "download", float64(downloadSpeed)/(1<<17), me, server.JSON())
			post(*httpClient, targetURL, payload)
		}
	}

	if opts.Upload {
		uploadSpeed = server.UploadSpeed()
		reportSpeed(opts, "Upload", uploadSpeed)
		if *host != "" {
			payload = fmt.Sprintf(`{ "metric": { "name": "%s", "value": %.4f, "units": "Mb" }, "client": "%s", "server": %s}`, "upload", float64(uploadSpeed)/(1<<17), me, server.JSON())
			post(*httpClient, targetURL, payload)
		}
	}
}

func post(httpClient http.Client, targetUrl, payload string) {
	jsonStr := []byte(payload)
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func reportSpeed(opts *speedtest.Opts, prefix string, speed int) {
	if !opts.Quiet {
		if opts.SpeedInBytes {
			log.Printf("%s: %.2f MB/s\n", prefix, float64(speed)/(1<<20))
		} else {
			log.Printf("%s: %.2f Mb/s\n", prefix, float64(speed)/(1<<17))
		}
	}
}

/*
func selectServer(opts *speedtest.Opts, client *speedtest.Client) (selected *speedtest.Server) {
	if opts.Server != 0 {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatalf("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.Find(opts.Server)
		if selected == nil {
			log.Fatalf("Server not found: %d\n", opts.Server)
			return nil
		}
		selected.MeasureLatency(speedtest.DefaultLatencyMeasureTimes, speedtest.DefaultErrorLatency)
	} else {
		servers, err := client.ClosestServers()
		if err != nil {
			log.Fatalf("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.MeasureLatencies(
			speedtest.DefaultLatencyMeasureTimes,
			speedtest.DefaultErrorLatency).First()
	}

	if !opts.Quiet {
		client.Log("Hosted by %s (%s) [%.2f km]: %d ms\n",
			selected.Sponsor,
			selected.Name,
			selected.Distance,
			selected.Latency/time.Millisecond)
	}

	return selected
}
*/
