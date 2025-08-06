package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	chrwr "github.com/i-b8o/chromedp_wrapper"
	"github.com/marketconnect/cons_tool_server/app"
)

const usage = `Usage:
parser -u <url> [-headless=false] [-help]
	
Options:
-u          Required, the document's root url to start processing
-headless   Optional, run browser in headless mode (default: true). Use -headless=false to show browser window.
-help       Optional, Prints this message
 `

func main() {
	rootUrl := flag.String("u", "", "start url")
	headless := flag.Bool("headless", true, "run browser in headless mode")
	help := flag.Bool("help", false, "Optional, prints usage info")
	flag.Parse()

	if *help {
		fmt.Println(usage)
		return
	}

	if *rootUrl == "" {
		log.Fatal("Error: rootUrl is empty. Use -u flag to provide it.\n", usage)
		return
	}

	// cfg := config.GetConfig() // Строка удалена, так как cfg не используется

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", *headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	parseCtx, parseCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer parseCancel()

	log.Println("Chrome wrapper initialisation")
	c := chrwr.NewChromeWrapper()
	c.SetTimeout(120)

	// Исправлен вызов NewApp: убран аргумент cfg
	a, err := app.NewApp(c, *rootUrl)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	log.Println("Running Application...")

	a.Process(parseCtx)
}
