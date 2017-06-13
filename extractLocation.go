package main

import (
	"net/http"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/yhat/scrape"
	"fmt"
	"net/url"
	"strings"
	"log"
	"strconv"
	"sync"
)

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getSource(url string) *html.Node {
	response, err := http.Get(url)
	handleError(err)
	defer response.Body.Close()

	source, err := html.Parse(response.Body)
	handleError(err)

	return source
}

func getLocationUrl(source *html.Node) string {
	//lat&long are hidden in src of <img class="_a3f img"/>
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Img && scrape.Attr(n, "class") == "_a3f img"
	}

	img, ok := scrape.Find(source, matcher)

	if ok == true {
		for _, attribute := range img.Attr {
			if attribute.Key == "src" {
				return attribute.Val
				break
			}
		}
	}

	return ""
}

func locationUrlToCoordinate(uri string) Coordinate {
	//example URL: https://external-waw1-1.xx.fbcdn.net/static_map.php?v=29&osm_provider=2&size=820x242&bbox=51.111570%2C16.981860%7C51.115570%2C16.997860&markers=51.11357000%2C16.99386000&language=pl_PL
	u, err := url.Parse(uri)
	handleError(err)

	p, err := url.ParseQuery(u.RawQuery)
	handleError(err)

	//lat&long are saved in &markers= parameter
	markers := strings.Split(p["markers"][0], ",")
	var cords [2]float64

	for i, marker := range markers {
		cords[i], err = strconv.ParseFloat(marker, 32)
		handleError(err)
	}

	return Coordinate{cords[0], cords[1]}
}

func getCoordinates(url string) {
	source := getSource(url)
	location := getLocationUrl(source)
	cords := locationUrlToCoordinate(location)

	fmt.Println(cords)
}

type Coordinate struct {
	lat float64
	long float64
}

var wg sync.WaitGroup

func main() {
	var foodTruck [9]string
	foodTruck[0] = "https://facebook.com/Tentego-food-drinks-383056401888734/"
	foodTruck[1] = "http://facebook.com/Pasibus/about/"
	foodTruck[2] = "http://facebook.com/nienazartybus/about"
	foodTruck[3] = "http://facebook.com/66AmericanBurger/about"
	foodTruck[4] = "http://facebook.com/balkanburgerPL/about"
	foodTruck[5] = "http://facebook.com/totutruck /about"
	foodTruck[6] = "http://facebook.com/bratwursty/about"
	foodTruck[7] = "http://facebook.com/Mojosandwiches/about"
	foodTruck[8] = "http://facebook.com/bagietyzfurgonety/about"

	for _, url := range foodTruck {
		wg.Add(1)
		go getCoordinates(url)
	}

	wg.Wait()
}
