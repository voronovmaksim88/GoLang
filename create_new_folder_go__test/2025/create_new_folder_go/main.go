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

	// Подключаемся к базе данных
	db, err := connectToDB()
	if err != nil {
		color.Red("Ошибка подключения к базе данных: %v", err)
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

	// Получаем текущий год из структуры папок
	currentYear, err := getCurrentYear()
	if err != nil {
		color.Red("Ошибка определения года: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	color.Green("Будут проанализированы заказы за %d год", currentYear)
	fmt.Printf("Год сохранен в переменную: %d\n", currentYear)

	waitForEnter()
}

// connectToDB устанавливает соединение с базой данных
func connectToDB() (*sql.DB, error) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("NAME"),
	)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %v", err)
	}

	color.Green("Успешное подключение к базе данных MariaDB!")
	return db, nil
}

// getCurrentYear определяет год на основе структуры папок
func getCurrentYear() (int, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения текущей директории: %v", err)
	}

	parentDir := filepath.Dir(currentDir)
	folderName := filepath.Base(parentDir)

	year, err := strconv.Atoi(folderName)
	if err != nil || year < 2000 || year > 9999 {
		return 0, fmt.Errorf("перенесите папку со скриптом в папку с заказами")
	}

	return year, nil
}

// waitForEnter ожидает нажатия Enter с обработкой возможной ошибки
func waitForEnter() {
	fmt.Println("\nНажмите Enter для выхода...")
	_, err := fmt.Scanln()
	if err != nil && err.Error() != "unexpected newline" {
		color.Red("Ошибка при ожидании ввода: %v", err)
	}
}
