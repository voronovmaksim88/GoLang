package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		color.Red("Ошибка загрузки .env файла: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Формируем строку подключения для MariaDB
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("NAME"),
	)

	// Пытаемся подключиться к базе данных
	db, err := sql.Open("mysql", connString)
	if err != nil {
		color.Red("Ошибка открытия соединения с базой данных: %v", err)
		waitForEnter()
		os.Exit(1)
	}
	defer func() {
		if db != nil {
			if err := db.Close(); err != nil {
				color.Red("Ошибка закрытия соединения с БД: %v", err)
			}
		}
	}()

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		color.Red("Ошибка ping базы данных: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Выводим сообщение об успехе
	color.Green("Успешное подключение к базе данных MariaDB!")

	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		color.Red("Ошибка получения текущей директории: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Получаем родительскую директорию
	parentDir := filepath.Dir(currentDir)
	folderName := filepath.Base(parentDir)

	// Проверяем, является ли имя родительской папки числом (годом)
	year, err := strconv.Atoi(folderName)
	if err != nil || year < 2000 || year > 9999 {
		color.Red("Перенесите папку со скриптом в папку с заказами")
		waitForEnter()
		os.Exit(1)
	}

	// Выводим сообщение и сохраняем год
	color.Green("Будут проанализированы заказы за %d год", year)

	// Сохраняем год в переменную для дальнейшего использования
	currentYear := year

	// Для примера выводим сохраненный год
	fmt.Printf("Год сохранен в переменную: %d\n", currentYear)

	waitForEnter()
}

// waitForEnter ожидает нажатия Enter с обработкой возможной ошибки
func waitForEnter() {
	fmt.Println("\nНажмите Enter для выхода...")
	_, err := fmt.Scanln()
	if err != nil && err.Error() != "unexpected newline" {
		color.Red("Ошибка при ожидании ввода: %v", err)
	}
}
