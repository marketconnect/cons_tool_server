package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	chrwr "github.com/i-b8o/chromedp_wrapper"
	// "github.com/marketconnect/cons_parser/config" // Больше не используется
)

// Result определяет структуру для итогового JSON-ответа.
type Result struct {
	RootURL string `json:"rootUrl"`
	DocName string `json:"doc_name"`
	SecURL  string `json:"sec_url"`
}

// App содержит конфигурацию и состояние приложения.
type App struct {
	// cfg     *config.Config // Поле cfg удалено
	rootUrl string
	c       *chrwr.Chrome
}

// NewApp создает новый экземпляр приложения.
// Параметр config удален из сигнатуры.
func NewApp(c *chrwr.Chrome, rootUrl string) (App, error) {
	if rootUrl == "" {
		return App{}, fmt.Errorf("rootUrl cannot be empty")
	}
	app := App{
		// cfg:     config, // Строка удалена
		rootUrl: rootUrl,
		c:       c,
	}
	return app, nil
}

// Process выполняет основную логику согласно требованиям.
func (a *App) Process(ctx context.Context) {
	log.Printf("Starting processing for rootUrl: %s\n", a.rootUrl)

	// 1. Открыть страницу, переданную в rootUrl.
	if err := a.c.OpenURL(ctx, a.rootUrl); err != nil {
		log.Fatalf("Failed to open root URL: %v", err)
		return
	}
	if err := a.c.WaitLoaded(ctx); err != nil {
		log.Fatalf("Failed to wait for root URL to load: %v", err)
		return
	}

	// 2. Сохранить doc_name из текста.
	docNameSelector := `body > div.content.document-page > section > div.external-block > div.external-block__content > div.document-page__title`
	jsGetDocName := fmt.Sprintf(`document.querySelector('%s').innerText`, docNameSelector)

	docName, err := a.c.GetString(ctx, jsGetDocName)
	if err != nil {
		log.Fatalf("Failed to get document name using selector '%s': %v", docNameSelector, err)
		return
	}
	docName = strings.TrimSpace(docName)
	log.Printf("Found document name: '%s'\n", docName)

	// 3. Перейти по ссылке из первого элемента оглавления.
	firstTocLinkSelector := `body > div.content.document-page > section > div.external-block > div.external-block__content > div.document-page__toc > ul > li:nth-child(1) > a`
	if err := a.c.Click(ctx, firstTocLinkSelector); err != nil {
		log.Fatalf("Failed to click first TOC link using selector '%s': %v", firstTocLinkSelector, err)
		return
	}
	log.Println("Clicked the first link in TOC.")
	if err := a.c.WaitLoaded(ctx); err != nil {
		log.Fatalf("Failed to wait for new page to load after click: %v", err)
		return
	}
	log.Println("New page loaded.")

	// 4. Получить URL текущей страницы после полной загрузки.
	secURL, err := a.c.GetString(ctx, `window.location.href`)
	if err != nil {
		log.Fatalf("Failed to get current URL: %v", err)
	}

	// 5. Возвращаем JSON.
	if secURL != "" {
		result := Result{
			RootURL: a.rootUrl,
			DocName: docName,
			SecURL:  secURL,
		}
		jsonOutput, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal result to JSON: %v", err)
		}
		fmt.Println(string(jsonOutput)) // Выводим JSON в stdout
	} else {
		log.Println("Could not determine the secondary URL.")
	}
}

