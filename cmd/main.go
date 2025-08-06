package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	chrwr "github.com/i-b8o/chromedp_wrapper"
	"github.com/marketconnect/cons_tool_server/app"
)

func main() {
	headless := flag.Bool("headless", true, "run browser in headless mode")
	flag.Parse()

	if !*headless {
		log.Fatal("This application must be run in headless mode to start the web service.")
	}

	http.HandleFunc("/", handleRequest)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	if number == "" {
		http.Error(w, `{"error": "number query parameter is required"}`, http.StatusBadRequest)
		return
	}

	rootUrl := fmt.Sprintf("https://www.consultant.ru/document/cons_doc_LAW_%s/", number)
	log.Printf("Processing request for number: %s, URL: %s", number, rootUrl)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	parseCtx, parseCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer parseCancel()

	c := chrwr.NewChromeWrapper()
	c.SetTimeout(120)

	a, err := app.NewApp(c, rootUrl)
	if err != nil {
		log.Printf("Error creating app: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "failed to create app: %v"}`, err), http.StatusInternalServerError)
		return
	}

	result, err := a.Process(parseCtx)
	if err != nil {
		log.Printf("Error processing request: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "failed to process request: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}