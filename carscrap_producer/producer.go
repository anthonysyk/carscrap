package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"carscrap/pkg/carscrap_redis"
	"strconv"
	"strings"
	"sync"
)

type ScrapingTracker struct {
	sync.Mutex
	PagesCounter   int
	ElementCounter int
}

func (t *ScrapingTracker) IncrementElementCounter() {
	t.Lock()
	t.ElementCounter++
	t.Unlock()
}

func (t *ScrapingTracker) IncrementPageCounter() {
	t.Lock()
	t.PagesCounter++
	t.Unlock()
}

func main() {
	// Init Redis
	svc := carscrap_redis.New(&carscrap_redis.NewInput{
		RedisURL: "0.0.0.0:6379",
	})
	defer svc.Close()

	// Find Last Page
	lastPage, err := FindLastPage()
	if err != nil {
		log.Fatalf("Could not find last page: %v", err)
	}
	fmt.Printf("Last Page is: %d \n", lastPage)

	// Start Scraping
	parallelism := 3
	fmt.Printf("Starting Scraping with %d Threads \n", parallelism)

	tracker := &ScrapingTracker{}
	StartScraping(svc, tracker, lastPage, parallelism)

	// Close Producer
	fmt.Println("Scraping Job finished:")
	fmt.Printf("Elements found: %d", tracker.PagesCounter)
	fmt.Printf("Pages screped: %d", tracker.ElementCounter)
	_ = svc.Publish("report", "done")
}

func FindLastPage() (int, error) {
	prefix := "137.0.-1.-1.-1.0.999999.1900.999999.-1.99.0."
	suffix := "?fulltext=&geoban=M137R99"

	c := colly.NewCollector()

	var lastPage int
	var err error
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, prefix) && strings.Contains(href, suffix) {
			href = strings.Replace(href, prefix, "", 1)
			href = strings.Replace(href, suffix, "", 1)
			page, err := strconv.Atoi(href)
			if err != nil {
				return
			}
			if page > lastPage {
				lastPage = page
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Finding Last Page: ", r.URL)
	})

	c.Visit("http://www.autoreflex.com/137.0.-1.-1.-1.0.999999.1900.999999.-1.99.0.1?fulltext=&geoban=M137R99")
	c.Wait()

	return lastPage, err
}

func StartScraping(svc *carscrap_redis.Service, tracker *ScrapingTracker, lastPage int, threads int) {
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "**",
		Parallelism: threads,
		//Delay:      5 * time.Second,
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "pic" {
			err := svc.Publish("worker", e.Request.AbsoluteURL(e.Attr("href")))
			if err != nil {
				log.Fatal(err)
			}
			tracker.IncrementElementCounter()
			return
		}
	})

	c.OnResponse(func(f *colly.Response) {
		tracker.PagesCounter++
		fmt.Printf("Total elements found: %d \n", tracker.ElementCounter)
		fmt.Printf("Total pages scraped: %d \n", tracker.PagesCounter)
		_ = svc.Publish("report", fmt.Sprintf("Pages: %v, Elements: %v", tracker.PagesCounter, tracker.ElementCounter))
	})

	for i := 1; i <= lastPage; i++ {
		c.Visit("http://www.autoreflex.com/137.0.-1.-1.-1.0.999999.1900.999999.-1.99.0." + strconv.Itoa(i))
	}

	c.Wait()
}
