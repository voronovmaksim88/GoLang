package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

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

	// Получаем множество заказов на основе имён папок
	orders, err := getOrderNumbersInFolder()
	if err != nil {
		color.Red("Ошибка получения списка заказов: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Проверяем, есть ли заказы
	if len(orders) == 0 {
		color.Yellow("\nНе найдены папки с заказами, продолжить? (y/n)")
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil || (answer != "y" && answer != "Y") {
			color.Yellow("Работа программы прервана.")
			waitForEnter()
			os.Exit(0)
		}
	} else {
		// Выводим номера заказов
		printOrders(orders)
	}

	// Создаем словарь заказов из БД
	dbOrderDict, err := createMariaDBOrderDict(db, currentYear)
	if err != nil {
		color.Red("Ошибка создания словаря заказов из БД: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	if len(dbOrderDict) == 0 {
		color.Yellow("В базе данных нет заказов за %d год", currentYear)
	} else {
		color.Green("\nУспешно загружено %d заказов из базы данных", len(dbOrderDict))
	}

	// Создаем словарь клиентов из БД
	clientDict, err := createMariaDBClientDict(db)
	if err != nil {
		color.Red("Ошибка создания словаря клиентов из БД: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	if len(clientDict) == 0 {
		color.Yellow("В базе данных нет клиентов")
	} else {
		color.Green("Успешно загружено %d клиентов из базы данных", len(clientDict))
	}

	printOrderStats(orders, dbOrderDict)

	waitForEnter()
}

// printOrders выводит номера заказов отсортированными по 5 в строке
func printOrders(orders map[string]bool) {
	color.Cyan("\nНайденные номера заказов по именам папок:")
	const ordersPerLine = 5
	// Преобразуем map в slice для сортировки
	orderList := make([]string, 0, len(orders))
	for order := range orders {
		orderList = append(orderList, order)
	}

	// Сортируем по первым 3 цифрам
	sort.Slice(orderList, func(i, j int) bool {
		return orderList[i][:3] < orderList[j][:3]
	})

	// Выводим по 5 заказов в строке
	for i := 0; i < len(orderList); i += ordersPerLine {
		end := i + ordersPerLine
		if end > len(orderList) {
			end = len(orderList)
		}

		// Создаем строку с 5 заказами
		line := ""
		for _, order := range orderList[i:end] {
			line += fmt.Sprintf("%-12s", order) // Выравниваем по ширине 12 символов
		}
		fmt.Println(line)
	}
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

// getOrderNumbersInFolder возвращает множество номеров заказов в родительской папке
func getOrderNumbersInFolder() (map[string]bool, error) {
	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения текущей директории: %v", err)
	}

	// Переходим на уровень выше (родительская папка)
	parentDir := filepath.Dir(currentDir)

	// Создаем множество для хранения номеров заказов
	orderNumbers := make(map[string]bool)

	// Читаем содержимое родительской папки
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения директории: %v", err)
	}

	// Перебираем все элементы в директории
	for _, entry := range entries {
		// Проверяем, что это директория
		if entry.IsDir() {
			dirName := entry.Name()

			// Проверяем, что имя начинается с 3 цифр и имеет достаточную длину
			if len(dirName) >= 11 && isDigit(dirName[:3]) {
				// Берем первые 11 символов как номер заказа
				orderNumber := dirName[:11]
				orderNumbers[orderNumber] = true
			}
		}
	}

	return orderNumbers, nil
}

// isDigit проверяет, состоит ли строка только из цифр
func isDigit(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func createMariaDBOrderDict(db *sql.DB, year int) (map[string]string, error) {
	orderDict := make(map[string]string)
	yearStr := strconv.Itoa(year)

	rows, err := db.Query("SELECT serial, client FROM task WHERE serial LIKE ?", "%-%-"+yearStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			_, fprintfErr := fmt.Fprintf(os.Stderr, "ошибка закрытия rows: %v\n", closeErr)
			if fprintfErr != nil {
				fmt.Printf("ошибка записи в stderr: %v\n", fprintfErr)
			}
		}
	}()

	for rows.Next() {
		var serial, client string
		if err := rows.Scan(&serial, &client); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		parts := strings.Split(serial, "-")
		if len(parts) == 3 && len(parts[2]) == 4 && len(parts[1]) == 2 {
			orderDict[serial] = client
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return orderDict, nil
}

// createMariaDBClientDict создает словарь клиентов из MariaDB (id клиента: имя клиента)
func createMariaDBClientDict(db *sql.DB) (map[string]string, error) {
	// Создаем словарь для хранения результатов
	clientDict := make(map[string]string)

	// Выполняем SQL-запрос
	rows, err := db.Query("SELECT id, name FROM client")
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			_, fprintfErr := fmt.Fprintf(os.Stderr, "ошибка закрытия rows: %v\n", closeErr)
			if fprintfErr != nil {
				fmt.Printf("ошибка записи в stderr: %v\n", fprintfErr)
			}
		}
	}()

	// Обрабатываем результаты
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		clientDict[id] = name
	}

	// Проверяем ошибки после итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return clientDict, nil
}

// printOrderStats выводит статистику заказов
func printOrderStats(folderOrders map[string]bool, dbOrders map[string]string) {
	color.Cyan("\nСтатистика заказов:")
	color.Blue("Найдено в папках: %d", len(folderOrders))
	color.Blue("Найдено в базе данных: %d", len(dbOrders))

	// Проверяем, какие заказы из базы данных отсутствуют в папках
	missingCount := 0
	for order := range dbOrders {
		if _, exists := folderOrders[order]; !exists {
			missingCount++
		}
	}

	if missingCount > 0 {
		color.Yellow("Заказов в БД, но отсутствующих в папках: %d", missingCount)
	} else {
		color.Green("Все заказы из БД присутствуют в папках")
	}
}
