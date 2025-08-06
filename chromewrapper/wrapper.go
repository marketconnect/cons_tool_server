package chromewrapper

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

// Chrome - это обертка для управления браузером.
type Chrome struct {
	timeout time.Duration
}

// NewChromeWrapper создает новый экземпляр обертки.
func NewChromeWrapper() *Chrome {
	return &Chrome{
		timeout: 60 * time.Second, // Таймаут по умолчанию
	}
}

// SetTimeout устанавливает таймаут для отдельных операций.
func (c *Chrome) SetTimeout(seconds int) {
	c.timeout = time.Duration(seconds) * time.Second
}

// Init создает главный контекст браузера.
func Init() (context.Context, context.CancelFunc) {
	// Для отладки можно установить headless в false
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	// Функция для корректного закрытия всех контекстов
	cancelAll := func() {
		cancelCtx()
		cancelAlloc()
	}

	return ctx, cancelAll
}

// NavigateAndWait переходит по URL и ожидает, пока элемент станет видимым.
func (c *Chrome) NavigateAndWait(ctx context.Context, url string, waitSelector string) error {
	taskCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(waitSelector, chromedp.ByQuery),
	)
}

// GetString выполняет JavaScript и возвращает результат в виде строки.
func (c *Chrome) GetString(ctx context.Context, jsExpression string) (string, error) {
	var result string
	taskCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	err := chromedp.Run(taskCtx, chromedp.Evaluate(jsExpression, &result))
	return result, err
}

// GetStringSlice выполняет JavaScript и возвращает результат в виде среза строк.
func (c *Chrome) GetStringSlice(ctx context.Context, jsExpression string) ([]string, error) {
	var result []string
	taskCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	err := chromedp.Run(taskCtx, chromedp.Evaluate(jsExpression, &result))
	return result, err
}

// ClickAndWait кликает на элемент и ожидает появления другого элемента.
func (c *Chrome) ClickAndWait(ctx context.Context, clickSelector string, waitSelector string) error {
	taskCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return chromedp.Run(taskCtx,
		chromedp.Click(clickSelector, chromedp.NodeVisible),
		chromedp.WaitVisible(waitSelector, chromedp.ByQuery),
	)
}