package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Запрос директории у пользователя
	fmt.Print("Введите путь к директории для проверки: ")
	directory, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Ошибка при чтении ввода: %v\n", err)
		waitForEnter(reader)
		os.Exit(1)
	}

	// Удаляем символы новой строки из введённого пути
	directory = strings.TrimSpace(directory)

	// Проверка существования директории
	info, err := os.Stat(directory)
	if err != nil {
		fmt.Printf("Ошибка при доступе к директории %s: %v\n", directory, err)
		waitForEnter(reader)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Printf("%s не является директорией\n", directory)
		waitForEnter(reader)
		os.Exit(1)
	}

	fmt.Printf("Сканирование директории %s...\n", directory)

	// Счетчик найденных файлов с длинными путями
	longPathsCount := 0

	// Рекурсивный обход директории
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Ошибка доступа к пути %s: %v\n", path, err)
			return filepath.SkipDir
		}

		// Проверка длины пути
		if len(path) > 240 {
			fmt.Printf("Длинный путь (%d символов): %s\n", len(path), path)
			longPathsCount++
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка при обходе директории: %v\n", err)
		waitForEnter(reader)
		os.Exit(1)
	}

	fmt.Printf("\nСканирование завершено. Найдено файлов с длинными путями: %d\n", longPathsCount)

	// Ожидание нажатия Enter в конце программы
	waitForEnter(reader)
}

// Функция для ожидания нажатия Enter
func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nНажмите Enter для завершения...")
	_, _ = reader.ReadString('\n')
}
