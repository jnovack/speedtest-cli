package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/surol/speedtest-cli/speedtest"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func version() {
	fmt.Print(speedtest.Version)
}

func usage() {
	fmt.Fprint(os.Stderr, "Command line interface for testing internet bandwidth using speedtest.net.\n\n")
	flag.PrintDefaults()
}

func main() {
	var id = flag.String("id", "", "The id for this host (e.g. hostname)")
	var host = flag.String("host", "", "Server where metrics are collected")
	var port = flag.String("port", "", "Port where sever is listening on")
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

	client.Log("Testing from %s (%s)...\n", config.Client.ISP, config.Client.IP)

	server := selectServer(opts, client);

	downloadSpeed := server.DownloadSpeed()
	reportSpeed(opts, "Download", downloadSpeed)

	uploadSpeed := server.UploadSpeed()
	reportSpeed(opts, "Upload", uploadSpeed)

	payload := fmt.Sprintf(`{"host": "%s", "metric_name":"download", "value": "%d"}`, *id, downloadSpeed)
	var jsonStr = []byte(payload)
	targetURL := fmt.Sprintf("http://%s:%s/public/metrics", *host, *port)
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	payload = fmt.Sprintf(`{"host": "%s", "metric_name":"upload", "value": "%d"}`, *id, uploadSpeed)
	jsonStr = []byte(payload)
	req, err = http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func reportSpeed(opts *speedtest.Opts, prefix string, speed int) {
	if opts.SpeedInBytes {
		fmt.Printf("%s: %.2f MiB/s\n", prefix, float64(speed) / (1 << 20))
	} else {
		fmt.Printf("%s: %.2f Mib/s\n", prefix, float64(speed) / (1 << 17))
	}
}

func selectServer(opts *speedtest.Opts, client *speedtest.Client) (selected *speedtest.Server) {
	if opts.Server != 0 {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatal("Failed to load server list: %v\n", err)
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
			log.Fatal("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.MeasureLatencies(
			speedtest.DefaultLatencyMeasureTimes,
			speedtest.DefaultErrorLatency).First()
	}

	if opts.Quiet {
		log.Printf("Ping: %d ms\n", selected.Latency / time.Millisecond)
	} else {
		client.Log("Hosted by %s (%s) [%.2f km]: %d ms\n",
			selected.Sponsor,
			selected.Name,
			selected.Distance,
			selected.Latency / time.Millisecond)
	}

	return selected
}
