package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

const AvailableTaps = "8000"
const TapCount = "8000"

const SleepTime = "10m"

var wg sync.WaitGroup

func DoTap(authorizationToken string) error {
	url := "https://api.hamsterkombat.io/clicker/tap"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"count":%s,"availableTaps":%s,"timestamp":%s}`, AvailableTaps, TapCount, GetCurrentTimestamp()))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Accept-Language", "en,en-US;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Origin", "https://hamsterkombat.io")
	req.Header.Add("Referer", "https://hamsterkombat.io/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 13; 2201117TY Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/124.0.6367.179 Mobile Safari/537.36")
	req.Header.Add("X-Requested-With", "org.telegram.messenger")
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+authorizationToken)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"124\", \"Android WebView\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?1")
	req.Header.Add("sec-ch-ua-platform", "\"Android\"")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Println("successful auto tap")
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		log.Println(string(body))
	}

	return nil
}

func BoostTaps(authorizationToken string) error {
	url := "https://api.hamsterkombat.io/clicker/buy-boost"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"boostId":"BoostFullAvailableTaps","timestamp":%s}`, GetCurrentTimestamp()))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Accept-Language", "en,en-US;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Origin", "https://hamsterkombat.io")
	req.Header.Add("Referer", "https://hamsterkombat.io/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 13; 2201117TY Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/124.0.6367.179 Mobile Safari/537.36")
	req.Header.Add("X-Requested-With", "org.telegram.messenger")
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+authorizationToken)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"124\", \"Android WebView\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?1")
	req.Header.Add("sec-ch-ua-platform", "\"Android\"")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Println("successful boost taps")
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		log.Println(string(body))
	}

	return nil
}

func RunJob(authorizationToken string) {
	err := DoTap(authorizationToken)
	if err != nil {
		log.Println(err)
	}

	err = BoostTaps(authorizationToken)
	if err != nil {
		log.Println(err)
	}

	err = DoTap(authorizationToken)
	if err != nil {
		log.Println(err)
	}
}

func GetCurrentTimestamp() string {
	currentTime := time.Now()

	timestamp := currentTime.Unix()

	timestampStr := strconv.FormatInt(timestamp, 10)

	return timestampStr
}

func RunScheduler(authorizationToken, sleepTime string) {
	defer wg.Done()

	timeZone := "Asia/Tehran"
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal(err)
	}

	s := gocron.NewScheduler(loc)
	if err != nil {
		log.Println(err)
	}

	_, err = s.Every(sleepTime).Do(RunJob, authorizationToken)
	if err != nil {
		log.Fatal(err)
	}

	s.StartBlocking()
}

func main() {
	authorizationTokens := flag.String("t", "", "sleep times.(separate with ,)")
	flag.Parse()

	authorizationTokensList := strings.Split(*authorizationTokens, ",")

	for _, authToken := range authorizationTokensList {
		wg.Add(1)
		go RunScheduler(authToken, SleepTime)
	}

	wg.Wait()
}
