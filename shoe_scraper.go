package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type shoeData struct {
	Name  string
	Price float64
}

func main() {
	shoeType := flag.String("shoe-type", "", "the type of shoe to search for")
	outfile := flag.String("outfile", "", "the file to write to")

	flag.Parse()

	if shoeType == nil || len(*shoeType) == 0 {
		panic("Error: specify a non-empty shoe type. Exiting.")
	}
	if outfile == nil || len(*outfile) == 0 {
		panic("Error: specify a non-empty outfile name. Exiting.")
	}

	data := make([]shoeData, 0)
	c := colly.NewCollector()

	c.OnHTML("div.product-card__body", func(e *colly.HTMLElement) {
		var name string
		var price float64

		e.ForEach("div.product-card__title", func(i int, h *colly.HTMLElement) {
			// fmt.Printf("Name: %v\n", h.Text)
			name = h.Text
		})
		e.ForEach("div.product-price.is--current-price", func(i int, h *colly.HTMLElement) {
			// fmt.Printf("Price: %v\n", h.Text)
			var err error
			tmp := strings.Replace(h.Text, "$", "", -1)
			price, err = strconv.ParseFloat(tmp, 64)
			if err != nil {
				fmt.Printf("For price: %s encountered error: %s", h.Text, err.Error())
			}
		})

		data = append(data, shoeData{name, price})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	fmt.Printf("Requesting data for shoe type: %s\n", *shoeType)
	c.Visit("https://www.nike.com/w?q=" + *shoeType)

	dump, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Writing output to %s", *outfile)
	ioutil.WriteFile(*outfile, dump, 0644)
}
