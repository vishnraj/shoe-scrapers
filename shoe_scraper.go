package main

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	shoeType := flag.String("shoe-type", "", "the type of shoe to search for")

	flag.Parse()

	if shoeType == nil || len(*shoeType) == 0 {
		panic("Error: specify a non-empty shoe type. Exiting.")
	}

	c := colly.NewCollector()

	c.OnHTML("div.product-card__body", func(e *colly.HTMLElement) {
		e.ForEach("div.product-card__title", func(i int, h *colly.HTMLElement) {
			fmt.Printf("Name: %v\n", h.Text)
		})
		e.ForEach("div.product-price.is--current-price", func(i int, h *colly.HTMLElement) {
			fmt.Printf("Price: %v\n", h.Text)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	fmt.Printf("Requesting data for shoe type: %s\n", *shoeType)
	c.Visit("https://www.nike.com/w?q=" + *shoeType)
}
