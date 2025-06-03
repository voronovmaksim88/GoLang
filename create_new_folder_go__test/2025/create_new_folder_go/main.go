package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/fatih/color"
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		color.Red("Ошибка загрузки .env файла: %v", err)
		fmt.Println("\nНажмите Enter для выхода...")
		fmt.Scanln()
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
		fmt.Println("\nНажмите Enter для выхода...")
		fmt.Scanln()
		os.Exit(1)
	}
	defer db.Close()

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		color.Red("Ошибка ping базы данных: %v", err)
		fmt.Println("\nНажмите Enter для выхода...")
		fmt.Scanln()
		os.Exit(1)
	}

	// Выводим сообщение об успехе
	color.Green("Успешное подключение к базе данных MariaDB!")
	fmt.Println("\nНажмите Enter для выхода...")
	fmt.Scanln()
}