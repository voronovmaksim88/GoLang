package main

import (
	"database/sql" // Пакет для работы с базами данных SQL
	"fmt"          // Пакет для форматированного ввода-вывода
	"github.com/joho/godotenv"
	"io"
	"os"            // Пакет для работы с операционной системой (файлы, переменные окружения)
	"path/filepath" // Пакет для работы с путями файловой системы
	"sort"          // Пакет для сортировки данных
	"strconv"       // Пакет для преобразования строк в числа и обратно
	"strings"       // Пакет для работы со строками

	"github.com/fatih/color" // Пакет для цветного вывода в консоль
	_ "github.com/lib/pq"    // Драйвер Postgres
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

	// Загружаем переменные окружения из файла .env (например, данные для подключения к БД)
	err := godotenv.Load()
	if err != nil {
		color.Red("Ошибка загрузки .env файла: %v", err)
		waitForEnter()
		os.Exit(1)
	}

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

	// Выводим статистику заказов и получаем количество отсутствующих заказов
	missingCount := printOrderStats(orders, dbOrderDict)
	// Если есть отсутствующие заказы, предлагаем создать для них папки
	if missingCount > 0 {
		createMissingOrderFolders(orders, dbOrderDict, clientDict)
	}

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
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("USER"),
		os.Getenv("PASS"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("NAME"),
	)

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %v", err)
	}

	// Проверяем доступность базы данных
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %v", err)
	}

	color.Green("Успешное подключение к базе данных Postgres!")
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

// createPostgresOrderDict создает словарь заказов из таблицы order (номер заказа: ID клиента)
func createMariaDBOrderDict(db *sql.DB, year int) (map[string]string, error) {
	// Создаем словарь для хранения заказов
	orderDict := make(map[string]string)
	yearStr := strconv.Itoa(year)

	// Выполняем запрос к таблице order, фильтруя заказы по году
	rows, err := db.Query("SELECT serial, customer_id FROM orders WHERE serial LIKE $1", "%-%-"+yearStr)
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

// createMariaDBClientDict создает словарь клиентов из таблицы counterparty (ID клиента: имя клиента)
func createMariaDBClientDict(db *sql.DB) (map[string]string, error) {
	// Создаем словарь для хранения клиентов
	clientDict := make(map[string]string)

	// Выполняем запрос к таблице client
	rows, err := db.Query("SELECT id, name FROM counterparty")
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

// printOrderStats выводит статистику заказов, сравнивая данные из папок и базы данных, и возвращает количество отсутствующих заказов
func printOrderStats(folderOrders map[string]bool, dbOrders map[string]string) (missingCount int) {
	color.Cyan("\nСтатистика заказов:")
	// Выводим количество заказов, найденных в папках
	color.Blue("Найдено в папках: %d", len(folderOrders))
	// Выводим количество заказов, найденных в базе данных
	color.Blue("Найдено в базе данных: %d", len(dbOrders))

	// Подсчитываем заказы, которые есть в базе данных, но отсутствуют в папках
	missingCount = 0
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

	return missingCount
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

// createOrderFolder создаёт папку для заказа и необходимые подпапки, копируя шаблоны ТЗ и КП
func createOrderFolder(folderName string, clientDict map[string]string, dbOrders map[string]string) error {
	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("ошибка получения текущей директории: %v", err)
	}

	// Переходим на уровень выше (родительская папка)
	parentDir := filepath.Dir(currentDir)

	// Получаем имя клиента по ID из dbOrders, если оно есть
	clientName := "Unknown"
	if clientID, exists := dbOrders[folderName]; exists {
		if name, found := clientDict[clientID]; found {
			// Заменяем недопустимые символы в имени клиента на '_'
			invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
			clientName = name
			for _, char := range invalidChars {
				clientName = strings.ReplaceAll(clientName, char, "_")
			}
		}
	}

	// Формируем имя папки в формате НомерЗаказа_ИмяКлиента
	fullFolderName := fmt.Sprintf("%s_%s", folderName, clientName)
	// Формируем путь к основной папке заказа
	orderDir := filepath.Join(parentDir, fullFolderName)

	// Создаём основную папку заказа, если она ещё не существует
	if err := os.MkdirAll(orderDir, 0755); err != nil {
		color.Red("Ошибка при создании папки %s: %v", fullFolderName, err)
		return fmt.Errorf("ошибка создания папки %s: %v", fullFolderName, err)
	}
	color.Green("\nОсновная папка для заказа %s успешно создана", fullFolderName)

	// Список стандартных подпапок
	subfolders := []string{
		"Чеклисты", "Фото и видео", "ТЗ", "Счета входящие", "Схема",
		"ПО", "Паспорт", "КП", "Документы",
	}

	// Создаём каждую подпапку
	for _, subfolder := range subfolders {
		subfolderPath := filepath.Join(orderDir, subfolder)
		if err := os.MkdirAll(subfolderPath, 0755); err != nil {
			color.Red("Ошибка при создании подпапки %s: %v", subfolder, err)
			return fmt.Errorf("ошибка создания подпапки %s: %v", subfolder, err)
		}
	}

	// Копируем шаблон ТЗ
	tzTemplate := filepath.Join(currentDir, "ТЗ.odt")
	tzDestination := filepath.Join(orderDir, "ТЗ", fmt.Sprintf("%s_ТЗ_в1р1.odt", folderName))
	if _, err := os.Stat(tzTemplate); err == nil {
		if err := copyFile(tzTemplate, tzDestination); err != nil {
			color.Red("Ошибка копирования шаблона ТЗ для %s: %v", fullFolderName, err)
			return fmt.Errorf("ошибка копирования шаблона ТЗ: %v", err)
		}
		color.Green("Шаблон ТЗ скопирован для заказа %s", fullFolderName)
	} else {
		color.Yellow("Внимание: Файл шаблона ТЗ не найден по пути: %s", tzTemplate)
	}

	// Копируем шаблон КП
	kpTemplate := filepath.Join(currentDir, "КП.xls")
	kpDestination := filepath.Join(orderDir, "КП", fmt.Sprintf("%s_КП_в1р1.xls", folderName))
	if _, err := os.Stat(kpTemplate); err == nil {
		if err := copyFile(kpTemplate, kpDestination); err != nil {
			color.Red("Ошибка копирования шаблона КП для %s: %v", fullFolderName, err)
			return fmt.Errorf("ошибка копирования шаблона КП: %v", err)
		}
		color.Green("Шаблон КП скопирован для заказа %s", fullFolderName)
	} else {
		color.Yellow("Внимание: Файл шаблона КП не найден по пути: %s", kpTemplate)
	}

	return nil
}

// copyFile копирует файл из src в dst
func copyFile(src, dst string) error {
	// Открываем исходный файл для чтения
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("ошибка открытия исходного файла %s: %v", src, err)
	}
	// Закрываем исходный файл и логируем возможные ошибки
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil {
			color.Red("Ошибка закрытия исходного файла %s: %v", src, closeErr)
		}
	}()

	// Создаём целевой файл
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("ошибка создания целевого файла %s: %v", dst, err)
	}
	// Закрываем целевой файл и логируем возможные ошибки
	defer func() {
		if closeErr := destFile.Close(); closeErr != nil {
			color.Red("Ошибка закрытия целевого файла %s: %v", dst, closeErr)
		}
	}()

	// Копируем содержимое файла
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("ошибка копирования файла %s в %s: %v", src, dst, err)
	}

	// Устанавливаем права доступа для целевого файла
	if err := destFile.Chmod(0644); err != nil {
		return fmt.Errorf("ошибка установки прав для файла %s: %v", dst, err)
	}

	return nil
}

// createMissingOrderFolders запрашивает у пользователя создание папок для отсутствующих заказов и завершает программу при ошибке
func createMissingOrderFolders(folderOrders map[string]bool, dbOrders map[string]string, clientDict map[string]string) {
	// Собираем заказы, которые есть в базе данных, но отсутствуют в папках
	var missingOrders []string
	for order := range dbOrders {
		if _, exists := folderOrders[order]; !exists {
			missingOrders = append(missingOrders, order)
		}
	}

	// Запрашиваем подтверждение на создание папок
	color.Cyan("Создать папки для отсутствующих заказов? (y/n)")
	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil || (answer != "y" && answer != "Y") {
		// Если пользователь отказался или произошла ошибка ввода, выводим сообщение об отмене
		color.Yellow("Создание папок отменено пользователем")
		return
	}

	// Создаём папку для каждого отсутствующего заказа
	for _, order := range missingOrders {
		if err := createOrderFolder(order, clientDict, dbOrders); err != nil {
			// При ошибке выводим сообщение и завершаем программу
			color.Red("Не удалось создать папку для заказа %s: %v", order, err)
			waitForEnter()
			os.Exit(1)
		}
	}
}
