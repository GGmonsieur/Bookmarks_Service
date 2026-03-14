package logicfuncs

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Функция для извлечения Title из URL
func FetchTitle(url string) (string, error) {

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Проверяем что это точно HTML, а не картинка
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("not an html page")
	}

	limitReader := io.LimitReader(resp.Body, 1024*100)

	// Парсим HTML
	z := html.NewTokenizer(limitReader)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return "", fmt.Errorf("title not found")
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" {
				z.Next()
				return strings.TrimSpace(z.Token().Data), nil
			}
		}
	}
}
