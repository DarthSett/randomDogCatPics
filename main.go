package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//DateScrape returns the quote on a specified date
//Date should be written in mmm-dd-yyyy format with month being a 3 lettered abbreviation and 1st letter capital
//Date shouldn't preced more than a week from present date
func DateScrape(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page.
	time.Sleep(10 * time.Second)
	res, err := http.Get("http://eduro.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	date := strings.Split(vars["date"], "-")
	quote := ""
	// Find quote of the day
	doc.Find(".singlemeta").Each(func(i int, s *goquery.Selection) {
		if (date[0] == s.Find(".months").Text() || strings.Title(strings.ToLower(date[0])) == s.Find(".months").Text()) && date[1] == s.Find(".dates").Text() && date[2] == s.Find(".years").Text() {
			// get the quote
			quote = s.Find("p").Text()
			w.Write([]byte(quote))
			return
		}
	})

}

func Cat(w http.ResponseWriter, r *http.Request) {
	println(1)
	res, err := http.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %v", res.StatusCode)
	}
	var i []map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("err: ", err)
	}
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Fatalf("err: ", err)
	}
	url := fmt.Sprint(i[0]["url"])
	res.Body.Close()
	println(2, ": Url:", url)
	res, err = http.Get(url)
	if err != nil {
		log.Fatalf("err: ", err)
	}
	println("3: ContentLength: ", fmt.Sprint(res.ContentLength))
	w.Header().Set("Content-Length", fmt.Sprint(res.ContentLength))
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	println("4: ContentType: ", res.Header.Get("Content-Type"))
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Fatalf("err: ", err)
	}
	res.Body.Close()
	return
}

//TodayScrape returns the quote of the day
func TodayScrape(w http.ResponseWriter, _ *http.Request) {
	// Request the HTML page.
	println(1)
	res, err := http.Get("http://eduro.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	quote := ""
	//find quote of the day
	doc.Find(".singlemeta").Each(func(i int, s *goquery.Selection) {
		if i < 1 {
			// get the quote
			quote = s.Find("p").Text()
			w.Write([]byte(quote))
		}
	})
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/date/{date}", DateScrape)
	router.HandleFunc("/", TodayScrape)
	router.HandleFunc("/cat", Cat)

	log.Fatal(http.ListenAndServe(":"+ os.Getenv("PORT"), router))
}
