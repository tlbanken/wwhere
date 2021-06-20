package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/tlbanken/wwhere/cmd/wwhere/keys"

	//"encoding/json"
	"net/http"
)

// Auto generated from https://mholt.github.io/json-to-go/
type Response struct {
	Message string `json:"message"`
	Cod     string `json:"cod"`
	Count   int    `json:"count"`
	List    []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   int     `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Dt   int `json:"dt"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Sys struct {
			Country string `json:"country"`
		} `json:"sys"`
		Rain   interface{} `json:"rain"`
		Snow   interface{} `json:"snow"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"list"`
}

var cities = [...]string{
	"phoenix",
	"london",
	"toronto",
	"paris",
	"rome",
	"tokyo",
	"seoul",
	"seattle",
	"helsinki",
	"moscow",
	"beijing",
	"singapore",
	"amsterdam",
	"berlin",
	"madrid",
	"dubai",
	"delhi",
	"cairo",
	"shanghai",
	"boston",
	"sydney",
}

const URLBase = "https://community-open-weather-map.p.rapidapi.com/find?"

const URLRest = "&cnt=1&mode=null&lon=0&type=link%2C%20accurate&lat=0&units=imperial"

func main() {
	rand.Seed(time.Now().UTC().Unix())

	for {
		correctCity := getRandomCity()
		url := buildURL(correctCity)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-rapidapi-key", keys.APIKey)
		fmt.Println("Fetching Next Question...")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer resp.Body.Close()

		// printBody(resp)
		responseObject := buildResponseObject(resp)
		correct, ans := askQuestion(responseObject)
		if correct {
			fmt.Println("Correct!")
		} else {
			fmt.Printf("Wrong! (Answer: %v)\n", ans)
		}
		fmt.Println()
	}

	// fmt.Println("------------------------")
	// fmt.Printf("%+v\n", responseObject)
}

// func printBody(resp *http.Response) {
// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	fmt.Println(string(b))
// }

func buildResponseObject(resp *http.Response) Response {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	var responseObject Response
	json.Unmarshal(bodyBytes, &responseObject)
	return responseObject
}

func buildURL(city string) string {
	return URLBase + "q=" + city + URLRest
}

func getRandomCity() string {
	index := rand.Int() % len(cities)
	return cities[index]
}

func askQuestion(resp Response) (bool, string) {
	ans := strings.ToLower(resp.List[0].Name)
	var choices []string
	choices = append(choices, ans)
	for i := 0; i < 3; i++ {
		city := ""
		for city == "" || stringInSlice(city, choices) {
			city = getRandomCity()
		}
		choices = append(choices, city)
	}
	// shuffle answers
	choices = shuffle(choices)

	fmt.Printf("-------------------------------------\n")
	// fmt.Printf("Answer: %v\n", ans)
	fmt.Printf("Which city corresponds to the following weather conditions?\n")
	fmt.Printf("Temp: %v\n", resp.List[0].Main.Temp)
	fmt.Printf("Humidity: %v\n", resp.List[0].Main.Humidity)
	fmt.Printf("Pressure: %v\n", resp.List[0].Main.Pressure)
	fmt.Printf("Wind Speed: %v\n", resp.List[0].Wind.Speed)
	fmt.Printf("Description: %v\n", resp.List[0].Weather[0].Description)
	fmt.Printf("-------------------------------------\n")
	fmt.Printf("a) %v\n", choices[0])
	fmt.Printf("b) %v\n", choices[1])
	fmt.Printf("c) %v\n", choices[2])
	fmt.Printf("d) %v\n", choices[3])
	fmt.Println()

	// user input
	guessIndex := -1
	for guessIndex == -1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter City Choice (a,b,c,d): ")
		guess, _ := reader.ReadString('\n')
		guess = strings.TrimRight(guess, "\n")
		guess = strings.ToLower(guess)
		guessIndex = toIndex(guess)
		if guessIndex == -1 {
			fmt.Printf("Invalid option!\n")
		}
	}
	return choices[guessIndex] == ans, ans
}

func stringInSlice(s string, l []string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

func shuffle(list []string) []string {
	for i := 0; i < 10; i++ {
		r1 := rand.Int() % len(list)
		r2 := rand.Int() % len(list)
		list[r1], list[r2] = list[r2], list[r1]
	}
	return list
}

func toIndex(guess string) int {
	if guess == "a" {
		return 0
	}
	if guess == "b" {
		return 1
	}
	if guess == "c" {
		return 2
	}
	if guess == "d" {
		return 3
	}
	return -1
}
