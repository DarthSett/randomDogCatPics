package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

//DateScrape returns the quote on a specified date
//Date should be written in mmm-dd-yyyy format with month being a 3 lettered abbreviation and 1st letter capital
//Date shouldn't preced more than a week from present date
func DateScrape(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page.
	//time.Sleep(10 * time.Second)
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

//TodayScrape returns the quote of the day
func TodayScrape(w http.ResponseWriter, _ *http.Request) {
	// Request the HTML page.
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
	router.HandleFunc("/{date}", DateScrape)
	router.HandleFunc("/", TodayScrape)

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), router))
}
