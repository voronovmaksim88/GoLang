package main

import (
	"database/sql"  // Пакет для работы с базами данных SQL
	"fmt"           // Пакет для форматированного ввода-вывода
	"os"            // Пакет для работы с операционной системой (файлы, переменные окружения)
	"path/filepath" // Пакет для работы с путями файловой системы
	"sort"          // Пакет для сортировки данных
	"strconv"       // Пакет для преобразования строк в числа и обратно
	"strings"       // Пакет для работы со строками

	"github.com/fatih/color"           // Пакет для цветного вывода в консоль
	_ "github.com/go-sql-driver/mysql" // Драйвер MySQL для работы с MariaDB
	"github.com/joho/godotenv"         // Пакет для загрузки переменных окружения из .env файла
)

// Константы определены в глобальной области видимости
const (
	orderNumberLength = 11 // Длина номера заказа (XXX-MM-YYYY)
	orderNumberParts  = 3  // Количество частей в номере заказа
	orderYearLength   = 4  // Длина года (YYYY)
	orderMonthLength  = 2  // Длина месяца (MM)
	ordersPerLine     = 5  // Количество заказов в строке для вывода
	errDBQuery        = "ошибка выполнения запроса: %v"
	errScanRow        = "ошибка чтения строки: %v"
	errRowsProcessing = "ошибка при обработке результатов: %v"
	errCloseRows      = "ошибка закрытия rows: %v\n"
	errFprintf        = "ошибка записи в stderr: %v\n"
)

func main() {

	// Определяем текущий год на основе имени родительской папки
	currentYear, err := getCurrentYear()
	if err != nil {
		color.Red("Ошибка определения года: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Выводим информацию о годе анализа
	color.White("\nБудут проанализированы заказы за %d год", currentYear)

	// Получаем список заказов из имен папок в родительской директории
	orders, err := getOrderNumbersInFolder()
	if err != nil {
		color.Red("Ошибка получения списка заказов: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Проверяем, найдены ли заказы в папках
	if len(orders) == 0 {
		// Если заказы не найдены, запрашиваем у пользователя продолжение
		color.Yellow("\nНе найдены папки с заказами, продолжить? (y/n)")
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil || (answer != "y" && answer != "Y") {
			color.Yellow("Работа программы прервана.")
			waitForEnter()
			os.Exit(0)
		}
	} else {
		// Выводим список заказов из папок
		printOrders(orders)
	}

	// Загружаем переменные окружения из файла .env (например, данные для подключения к БД)
	err = godotenv.Load()
	if err != nil {
		color.Red("Ошибка загрузки .env файла: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Устанавливаем соединение с базой данных MariaDB
	db, err := connectToDB()
	if err != nil {
		color.Red("Ошибка подключения к базе данных: %v", err)
		waitForEnter()
		os.Exit(1)
	}
	// Закрываем соединение с БД при завершении функции main
	defer func() {
		if db != nil {
			if err := db.Close(); err != nil {
				color.Red("Ошибка закрытия соединения с БД: %v", err)
			}
		}
	}()

	// Создаем словарь заказов из базы данных (номер заказа: ID клиента)
	dbOrderDict, err := createMariaDBOrderDict(db, currentYear)
	if err != nil {
		color.Red("Ошибка создания словаря заказов из БД: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Выводим информацию о количестве загруженных заказов
	if len(dbOrderDict) == 0 {
		color.Yellow("В базе данных нет заказов за %d год", currentYear)
	} else {
		color.Green("Успешно загружено %d заказов из базы данных", len(dbOrderDict))
	}

	// Создаем словарь клиентов из базы данных (ID клиента: имя клиента)
	clientDict, err := createMariaDBClientDict(db)
	if err != nil {
		color.Red("Ошибка создания словаря клиентов из БД: %v", err)
		waitForEnter()
		os.Exit(1)
	}

	// Выводим информацию о количестве загруженных клиентов
	if len(clientDict) == 0 {
		color.Yellow("В базе данных нет клиентов")
	} else {
		color.Green("Успешно загружено %d клиентов из базы данных", len(clientDict))
	}

	// Выводим статистику заказов, сравнивая данные из папок и базы данных
	printOrderStats(orders, dbOrderDict)

	// Ожидаем нажатия Enter перед завершением программы
	waitForEnter()
}

// printOrders выводит номера заказов, отсортированные по первым трем цифрам, по 5 в строке
func printOrders(orders map[string]bool) {
	color.Cyan("\nНайденные номера заказов по именам папок:")
	// Преобразуем map в slice для сортировки
	orderList := make([]string, 0, len(orders))
	for order := range orders {
		orderList = append(orderList, order)
	}

	// Сортируем заказы по первым трем цифрам (XXX в XXX-MM-YYYY)
	sort.Slice(orderList, func(i, j int) bool {
		return orderList[i][:3] < orderList[j][:3]
	})

	// Выводим заказы по ordersPerLine (5) в строке
	for i := 0; i < len(orderList); i += ordersPerLine {
		end := i + ordersPerLine
		if end > len(orderList) {
			end = len(orderList)
		}

		// Формируем строку с выравниванием для текущей группы заказов
		line := ""
		for _, order := range orderList[i:end] {
			line += fmt.Sprintf("%-12s", order) // Выравниваем по ширине 12 символов
		}
		fmt.Println(line)
	}
	fmt.Println("")
}

// connectToDB устанавливает соединение с базой данных MariaDB
func connectToDB() (*sql.DB, error) {
	// Формируем строку подключения из переменных окружения
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("NAME"),
	)

	// Открываем соединение с базой данных
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %v", err)
	}

	// Проверяем доступность базы данных
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %v", err)
	}

	color.Green("Успешное подключение к базе данных MariaDB!")
	return db, nil
}

// getCurrentYear определяет текущий год на основе имени родительской папки
func getCurrentYear() (int, error) {
	// Получаем текущую рабочую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения текущей директории: %v", err)
	}

	// Получаем родительскую директорию и её имя
	parentDir := filepath.Dir(currentDir)
	folderName := filepath.Base(parentDir)

	// Преобразуем имя папки в число (год)
	year, err := strconv.Atoi(folderName)
	if err != nil || year < 2000 || year > 9999 {
		return 0, fmt.Errorf("перенесите папку со скриптом в папку с заказами")
	}

	return year, nil
}

// waitForEnter ожидает нажатия клавиши Enter для завершения программы
func waitForEnter() {
	fmt.Println("\nНажмите Enter для выхода...")
	_, err := fmt.Scanln()
	if err != nil && err.Error() != "unexpected newline" {
		color.Red("Ошибка при ожидании ввода: %v", err)
	}
}

// getOrderNumbersInFolder возвращает множество номеров заказов из имен папок в родительской директории
func getOrderNumbersInFolder() (map[string]bool, error) {
	// Получаем текущую рабочую директорию
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
		// Проверяем, что элемент является директорией
		if entry.IsDir() {
			dirName := entry.Name()

			// Проверяем, что имя папки соответствует формату заказа (начинается с 3 цифр и достаточно длинное)
			if len(dirName) >= orderNumberLength && isDigit(dirName[:3]) {
				orderNumber := dirName[:orderNumberLength]
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

// createMariaDBOrderDict создает словарь заказов из таблицы task (номер заказа: ID клиента)
func createMariaDBOrderDict(db *sql.DB, year int) (map[string]string, error) {
	// Создаем словарь для хранения заказов
	orderDict := make(map[string]string)
	yearStr := strconv.Itoa(year)

	// Выполняем запрос к таблице task, фильтруя заказы по году
	rows, err := db.Query("SELECT serial, client FROM task WHERE serial LIKE ?", "%-%-"+yearStr)
	if err != nil {
		return nil, fmt.Errorf(errDBQuery, err)
	}
	// Закрываем rows при завершении функции
	defer closeRows(rows)

	// Обрабатываем каждую строку результата
	for rows.Next() {
		var serial, client string
		if err := rows.Scan(&serial, &client); err != nil {
			return nil, fmt.Errorf(errScanRow, err)
		}
		// Проверяем формат номера заказа (XXX-MM-YYYY)
		parts := strings.Split(serial, "-")
		if len(parts) == orderNumberParts && len(parts[2]) == orderYearLength && len(parts[1]) == orderMonthLength {
			orderDict[serial] = client
		}
	}

	// Проверяем наличие ошибок после обработки всех строк
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(errRowsProcessing, err)
	}

	return orderDict, nil
}

// createMariaDBClientDict создает словарь клиентов из таблицы client (ID клиента: имя клиента)
func createMariaDBClientDict(db *sql.DB) (map[string]string, error) {
	// Создаем словарь для хранения клиентов
	clientDict := make(map[string]string)

	// Выполняем запрос к таблице client
	rows, err := db.Query("SELECT id, name FROM client")
	if err != nil {
		return nil, fmt.Errorf(errDBQuery, err)
	}
	// Закрываем rows при завершении функции
	defer closeRows(rows)

	// Обрабатываем каждую строку результата
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf(errScanRow, err)
		}
		clientDict[id] = name
	}

	// Проверяем наличие ошибок после обработки всех строк
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(errRowsProcessing, err)
	}

	return clientDict, nil
}

// printOrderStats выводит статистику заказов, сравнивая данные из папок и базы данных
func printOrderStats(folderOrders map[string]bool, dbOrders map[string]string) {
	color.Cyan("\nСтатистика заказов:")
	// Выводим количество заказов, найденных в папках
	color.Blue("Найдено в папках: %d", len(folderOrders))
	// Выводим количество заказов, найденных в базе данных
	color.Blue("Найдено в базе данных: %d", len(dbOrders))

	// Подсчитываем заказы, которые есть в базе данных, но отсутствуют в папках
	missingCount := 0
	for order := range dbOrders {
		if _, exists := folderOrders[order]; !exists {
			missingCount++
		}
	}

	// Выводим результат сравнения
	if missingCount > 0 {
		color.Yellow("Заказов в БД, но отсутствующих в папках: %d", missingCount)
	} else {
		color.Green("Все заказы из БД присутствуют в папках")
	}
}

// closeRows закрывает объект rows и логирует возможные ошибки
func closeRows(rows *sql.Rows) {
	if closeErr := rows.Close(); closeErr != nil {
		_, fprintfErr := fmt.Fprintf(os.Stderr, errCloseRows, closeErr)
		if fprintfErr != nil {
			fmt.Printf(errFprintf, fprintfErr)
		}
	}
}
