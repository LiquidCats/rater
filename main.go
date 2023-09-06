package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kyokomi/emoji"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Exchange struct {
	Rates map[string]float32 `json:"rates"`
}

func noErrors(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	err := godotenv.Load() //Load .env file
	noErrors(err)

	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayStr := fmt.Sprintf("%s", yesterday.Format("2006-01-02"))
	today := time.Now()
	todayStr := fmt.Sprintf("%s", today.Format("2006-01-02"))

	currenciesStr := os.Getenv("CURRENCIES")
	currencies := strings.Split(currenciesStr, ",")
	live := os.Getenv("API_URL_LIVE")
	history := os.Getenv("API_URL_HISTORY")
	history = fmt.Sprintf(history, yesterdayStr)

	fmt.Println("URL:", live)
	fmt.Println("CURRENCIES:", currenciesStr)

	respLive, err := http.Get(live)
	noErrors(err)

	respHistory, err := http.Get(history)
	noErrors(err)

	defer respLive.Body.Close()

	bodyLive, err := ioutil.ReadAll(respLive.Body)
	noErrors(err)

	bodyHistory, err := ioutil.ReadAll(respHistory.Body)
	noErrors(err)

	var dataLive Exchange
	var dataHistory Exchange

	err = json.Unmarshal(bodyLive, &dataLive)
	noErrors(err)

	err = json.Unmarshal(bodyHistory, &dataHistory)
	noErrors(err)

	var rates []string
	for _, currency := range currencies {
		previous := dataHistory.Rates[currency]
		current := dataLive.Rates[currency]
		percent := (current/previous - 1) * 100
		var chart string

		if previous > current {
			chart = ":chart_with_downwards_trend:"
		}

		if previous < current {
			chart = ":chart_with_upwards_trend:"
		}

		if previous == current {
			chart = ":bar_chart:"
		}

		r := emoji.Sprintf(":white_check_mark: %s:\n"+
			"$%f\n"+
			"%f%s %s\n"+
			"\n", currency, current, percent, "%", chart)
		fmt.Print(r)
		rates = append(rates, r)
	}

	filename := fmt.Sprintf("rates-%s.txt", todayStr)
	err = ioutil.WriteFile(filename, []byte(strings.Join(rates, "")), 0755)
	noErrors(err)
}
