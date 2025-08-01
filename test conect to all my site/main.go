package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

// Website Структура для хранения информации о сайте
type Website struct {
	URL         string
	DisplayName string
}

// CheckResult Результат проверки сайта
type CheckResult struct {
	Available  bool
	StatusCode int
	Error      error
	Duration   time.Duration
}

func checkWebsite(url string) CheckResult {
	startTime := time.Now()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	duration := time.Since(startTime)

	if err != nil {
		return CheckResult{
			Available: false,
			Error:     err,
			Duration:  duration,
		}
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("Ошибка при закрытии тела ответа: %v\n", cerr)
		}
	}()

	return CheckResult{
		Available:  resp.StatusCode >= 200 && resp.StatusCode < 400,
		StatusCode: resp.StatusCode,
		Error:      nil,
		Duration:   duration,
	}
}

func main() {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Список сайтов для проверки
	websites := []Website{
		{URL: "https://www.sibplc.ru/", DisplayName: "sibplc.ru"},
		{URL: "http://ovz1.9138995941.me2jm.vps.myjino.ru:49264/", DisplayName: "КИС1"},
		{URL: "https://kis2test.sibplc.ru/", DisplayName: "КИС2 тестовая"},
		{URL: "https://kis2work.sibplc.ru/", DisplayName: "КИС2 рабочая"},
		{URL: "https://test.sibplc-kis3.ru/", DisplayName: "КИС3 тестовая"},
		{URL: "https://sibplc-kis3.ru/", DisplayName: "КИС3 рабочая"},
		{URL: "https://sibtask.ru/", DisplayName: "sibtask"},
		{URL: "https://lk.sinlab.ru:3080/comp/", DisplayName: "ecomp"},
	}

	// Проверка каждого сайта
	for _, site := range websites {
		result := checkWebsite(site.URL)

		if result.Error != nil {
			fmt.Printf("%s: %v (время: %s)\n",
				red(fmt.Sprintf("Ошибка %s", site.DisplayName)),
				result.Error,
				yellow(result.Duration.Round(time.Millisecond)))
		} else if result.Available {
			fmt.Printf("%s (время: %s)\n",
				green(fmt.Sprintf("%s доступен", site.DisplayName)),
				yellow(result.Duration.Round(time.Millisecond)))
		} else {
			fmt.Printf("%s (статус код: %d, время: %s)\n",
				red(fmt.Sprintf("%s недоступен", site.DisplayName)),
				result.StatusCode,
				yellow(result.Duration.Round(time.Millisecond)))
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nНажмите Enter для выхода...")
	_, _ = reader.ReadString('\n')
}
