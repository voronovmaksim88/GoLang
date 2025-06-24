package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

func printError(message string) {
	color.Red(message)
}

func printSuccess(message string) {
	color.Green(message)
}

func waitForEnter() {
	fmt.Print("\nНажмите Enter для выхода...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getDefaultInput(prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func tryConnect(host, user, password string) (*ssh.Client, error) {
	addr := host + ":22"
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", addr, config)
}

func copyFile(session *ssh.Session, localPath string, remotePath string) error {
	// Читаем локальный файл
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла %s: %v", localPath, err)
	}

	// Создаем команду для записи файла
	cmd := fmt.Sprintf("cat > %s", remotePath)
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания stdin pipe: %v", err)
	}

	// Запускаем команду
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("ошибка запуска команды: %v", err)
	}

	// Отправляем содержимое файла
	if _, err := stdin.Write(content); err != nil {
		return fmt.Errorf("ошибка записи в stdin: %v", err)
	}

	// Закрываем stdin и ждем завершения команды
	if err := stdin.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия stdin: %v", err)
	}
	if err := session.Wait(); err != nil {
		return fmt.Errorf("ошибка ожидания завершения команды: %v", err)
	}

	return nil
}

func copyFileRemote(conn *ssh.Client, localPath, remotePath string) error {
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("не удалось создать сессию: %v", err)
	}

	// Используем анонимную функцию для обработки ошибки закрытия сессии
	defer closeSession(session)

	// Читаем локальный файл
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла %s: %v", localPath, err)
	}

	// Создаем команду для записи файла
	cmd := fmt.Sprintf("cat > %s", remotePath)
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания stdin pipe: %v", err)
	}

	// Запускаем команду
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("ошибка запуска команды: %v", err)
	}

	// Отправляем содержимое файла
	if _, err := stdin.Write(content); err != nil {
		return fmt.Errorf("ошибка записи в stdin: %v", err)
	}

	// Закрываем stdin и ждем завершения команды
	if err := stdin.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия stdin: %v", err)
	}
	if err := session.Wait(); err != nil {
		return fmt.Errorf("ошибка ожидания завершения команды: %v", err)
	}

	return nil
}

// Новая функция для закрытия сессии
func closeSession(session *ssh.Session) {
	if session != nil {
		if err := session.Close(); err != nil {
			color.Yellow("Предупреждение: ошибка при закрытии сессии: %v", err)
		}
	}
}

func closeConnection(conn *ssh.Client) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			color.Yellow("Предупреждение: ошибка при закрытии соединения: %v", err)
		}
	}
}

func main() {
	// Обработка паники и аварийного завершения
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\nПроизошла непредвиденная ошибка: %v\n", r)
			waitForEnter()
		}
	}()

	// 1. УСТАНОВКА СОЕДИНЕНИЯ
	fmt.Println("=== Установка SSH соединения ===")

	var conn *ssh.Client
	var err error

	// Цикл подключения с повторными попытками
	for {
		// Запрос учетных данных
		host := getDefaultInput("Введите IP-адрес устройства", "192.168.88.60")
		user := getDefaultInput("Введите имя пользователя", "root")
		password := getDefaultInput("Введите пароль", "segnetics")

		fmt.Printf("\nПопытка подключения к %s@%s...\n", user, host)

		// Установка соединения
		conn, err = tryConnect(host, user, password)
		if err == nil {
			printSuccess("SSH соединение установлено успешно!")
			break
		}

		printError(fmt.Sprintf("Ошибка подключения: %v", err))
		retry := getUserInput("Повторить попытку? (y/n): ")
		if strings.ToLower(retry) != "y" {
			waitForEnter()
			return
		}
	}
	// Гарантированное закрытие соединения при выходе
	defer closeConnection(conn)

	// 2. ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ДЛЯ РАБОТЫ С СЕССИЯМИ

	// executeCommand - выполняет команду через SSH и возвращает результат
	executeCommand := func(cmd string) (string, error) {
		session, err := conn.NewSession()
		if err != nil {
			return "", fmt.Errorf("ошибка создания сессии: %v", err)
		}
		defer closeSession(session)

		output, err := session.CombinedOutput(cmd)
		if err != nil {
			return "", fmt.Errorf("ошибка выполнения '%s': %v", cmd, err)
		}
		return string(output), nil
	}

	// checkPathExists - проверяет существование пути на удаленной машине
	checkPathExists := func(path string, isFile bool) (bool, error) {
		var cmd string
		if isFile {
			cmd = fmt.Sprintf("test -f %s && echo 'exists' || echo 'not exists'", path)
		} else {
			cmd = fmt.Sprintf("test -d %s && echo 'exists' || echo 'not exists'", path)
		}

		output, err := executeCommand(cmd)
		if err != nil {
			return false, err
		}
		return strings.Contains(output, "exists"), nil
	}

	// 3. ТЕСТИРОВАНИЕ СОЕДИНЕНИЯ
	fmt.Println("\n=== Тестирование соединения ===")

	if output, err := executeCommand("echo 'Тестовое соединение'"); err != nil {
		printError(err.Error())
		waitForEnter()
		return
	} else {
		printSuccess("Соединение работает корректно")
		fmt.Printf("Ответ: %s\n", strings.TrimSpace(output))
	}

	// 4. ПРОВЕРКА СИСТЕМНОЙ ИНФОРМАЦИИ
	fmt.Println("\n=== Системная информация ===")

	// Получение информации о версии ядра
	if output, err := executeCommand("cat /proc/version"); err != nil {
		printError(fmt.Sprintf("Ошибка получения версии ядра: %v", err))
	} else {
		fmt.Printf("Версия ядра:\n%s\n", output)
	}

	// 5. РАБОТА С ДИРЕКТОРИЕЙ /etc/opt
	fmt.Println("\n=== Проверка рабочей директории ===")

	exists, err := checkPathExists("/etc/opt", false)
	if err != nil {
		printError(fmt.Sprintf("Ошибка проверки директории: %v", err))
		waitForEnter()
		return
	}

	if !exists {
		printError("Критическая ошибка: директория /etc/opt не существует!")
		waitForEnter()
		return
	}
	printSuccess("Директория /etc/opt доступна")

	// 6. КОПИРОВАНИЕ ФАЙЛОВ
	fmt.Println("\n=== Копирование файлов ===")

	filesToCopy := []string{"sequencer_v1r6.php", "segnetics.php"}
	for _, file := range filesToCopy {
		localPath := filepath.Join(".", file)

		// Проверка локального файла
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			printError(fmt.Sprintf("Локальный файл %s не найден", file))
			continue
		}

		// Копирование на удаленный сервер
		remotePath := "/etc/opt/" + file
		fmt.Printf("Копирование %s -> %s...\n", localPath, remotePath)

		if err := copyFileRemote(conn, localPath, remotePath); err != nil {
			printError(fmt.Sprintf("Ошибка копирования: %v", err))
		} else {
			printSuccess("Файл успешно скопирован")
		}
	}

	// 7. РАБОТА С ПРОЕКТАМИ
	fmt.Println("\n=== Проверка директории projects ===")

	exists, err = checkPathExists("/projects", false)
	if err != nil {
		printError(fmt.Sprintf("Ошибка проверки директории: %v", err))
		waitForEnter()
		return
	}

	if !exists {
		printError("Директория projects не существует!")
		waitForEnter()
		return
	}
	printSuccess("Директория projects доступна")

	// 8. РАБОТА С ФАЙЛОМ start.after
	fmt.Println("\n=== Работа с start.after ===")

	startAfterPath := "/projects/start.after"
	exists, err = checkPathExists(startAfterPath, true)
	if err != nil {
		printError(fmt.Sprintf("Ошибка проверки файла: %v", err))
	} else if exists {
		printSuccess("Файл start.after существует")

		// Запрос на удаление файла
		if getUserInput("Удалить файл start.after? (y/n): ") == "y" {
			if _, err := executeCommand("rm -f " + startAfterPath); err != nil {
				printError(fmt.Sprintf("Ошибка удаления: %v", err))
			} else {
				printSuccess("Файл успешно удален")
			}
		}
	} else {
		printError("Файл start.after не найден")
	}

	// 9. КОПИРОВАНИЕ НОВОГО ФАЙЛА start.after
	fmt.Println("\n=== Обновление start.after ===")

	localStartAfter := "start.after"
	if _, err := os.Stat(localStartAfter); err == nil {
		fmt.Printf("Копирование %s -> %s\n", localStartAfter, startAfterPath)

		session, err := conn.NewSession()
		if err != nil {
			printError(fmt.Sprintf("Ошибка создания сессии: %v", err))
		} else {
			defer closeSession(session)

			if err := copyFile(session, localStartAfter, startAfterPath); err != nil {
				printError(fmt.Sprintf("Ошибка копирования: %v", err))
			} else {
				printSuccess("Файл успешно обновлен")

				// Проверка содержимого
				if output, err := executeCommand("head -n 5 " + startAfterPath); err != nil {
					printError(fmt.Sprintf("Ошибка проверки: %v", err))
				} else {
					fmt.Printf("Начало файла:\n%s\n", output)
				}
			}
		}
	} else {
		printError("Локальный файл start.after не найден")
	}

	// Завершение работы
	fmt.Println("\n=== Работа завершена ===")
	waitForEnter()
}
