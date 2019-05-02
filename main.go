package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	hostAddr = flag.String(
		"host",
		":80",
		"address and port for server to listen on (default: [:80])",
	)
	configFile = flag.String(
		"config",
		"example.config",
		"config file path - contains IPs to whitelist",
	)
)

func main() {
	flag.Parse()

	allowedIPs, err := loadWhiteListIPs(*configFile)
	if err != nil {
		log.Println("error while parsing config file:", err)
		log.Println("Only requests from localhost will be allowed.")
	}
	log.Println("allowed IPs: ", allowedIPs)

	handler := http.NewServeMux()

	handler.HandleFunc("/", handleRequest(allowedIPs))

	log.Fatal(http.ListenAndServe(*hostAddr, handler))
}

func handleRequest(
	allowedIPs []string,
) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if allowIP(r.RemoteAddr, allowedIPs) {
			log.Printf("%v is whitelisted, request allowed\n", r.RemoteAddr)

			remoteURL := strings.TrimPrefix(r.URL.String(), "/")
			remoteURL = strings.Replace(remoteURL, "http:", "http://", 1)
			remoteURL = strings.Replace(remoteURL, "https:", "https://", 1)

			log.Println("requesting URL: ", remoteURL)

			// prepare the request
			defaultClient := &http.Client{
				Timeout: 15 * time.Second,
			}
			req, err := http.NewRequest(r.Method, remoteURL, r.Body)
			if err != nil {
				log.Println("error while preparing request: ", err)
				return
			}
			res, err := defaultClient.Do(req)
			if err != nil {
				log.Println("error while doing the request: ", err)
				return
			}

			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("error while reading response body: ", err)
				return
			}
			if _, err := w.Write(resBody); err != nil {
				log.Println("error while sending response: ", err)
				return
			}
			log.Println("request success!")

		} else {
			log.Printf("%v is not allowed to execute requests\n", r.RemoteAddr)
		}
	}
}

func allowIP(ip string, allowedIPs []string) bool {

	// allow requests from localhost by default
	if strings.HasPrefix(ip, "[::1]") || strings.HasPrefix(ip, "127.0.0.1") {
		return true
	}

	for _, allowed := range allowedIPs {
		if strings.HasPrefix(ip, allowed) {
			return true
		}
	}
	return false
}

func loadWhiteListIPs(filePath string) ([]string, error) {
	var addresses []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			continue
		}
		if len(line) == 0 {
			continue
		}

		addresses = append(addresses, scanner.Text())
	}

	return addresses, nil
}
