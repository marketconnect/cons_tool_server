package app

import (
	"context"
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
func (a *App) Process(ctx context.Context) (Result, error) {
	log.Printf("Starting processing for rootUrl: %s\n", a.rootUrl)

	// 1. Открыть страницу, переданную в rootUrl.
	if err := a.c.OpenURL(ctx, a.rootUrl); err != nil {
		return Result{}, fmt.Errorf("failed to open root URL: %w", err)
	}
	if err := a.c.WaitLoaded(ctx); err != nil {
		return Result{}, fmt.Errorf("failed to wait for root URL to load: %w", err)
	}

	// 2. Сохранить doc_name из текста.
	docNameSelector := `body > div.content.document-page > section > div.external-block > div.external-block__content > div.document-page__title`
	jsGetDocName := fmt.Sprintf(`document.querySelector('%s').innerText`, docNameSelector)

	docName, err := a.c.GetString(ctx, jsGetDocName)
	if err != nil {
		return Result{}, fmt.Errorf("failed to get document name using selector '%s': %w", docNameSelector, err)
	}
	docName = strings.TrimSpace(docName)
	log.Printf("Found document name: '%s'\n", docName)

	// 3. Перейти по ссылке из первого элемента оглавления.
	firstTocLinkSelector := `body > div.content.document-page > section > div.external-block > div.external-block__content > div.document-page__toc > ul > li:nth-child(1) > a`
	if err := a.c.Click(ctx, firstTocLinkSelector); err != nil {
		return Result{}, fmt.Errorf("failed to click first TOC link using selector '%s': %w", firstTocLinkSelector, err)
	}
	log.Println("Clicked the first link in TOC.")
	if err := a.c.WaitLoaded(ctx); err != nil {
		return Result{}, fmt.Errorf("failed to wait for new page to load after click: %w", err)
	}
	log.Println("New page loaded.")

	// 4. Получить URL текущей страницы после полной загрузки.
	secURL, err := a.c.GetString(ctx, `window.location.href`)
	if err != nil {
		return Result{}, fmt.Errorf("failed to get current URL: %w", err)
	}

	// 5. Возвращаем JSON.
	if secURL != "" {
		return Result{
			RootURL: a.rootUrl,
			DocName: docName,
			SecURL:  secURL,
		}, nil
	} else {
		return Result{}, fmt.Errorf("could not determine the secondary URL")
	}
}


