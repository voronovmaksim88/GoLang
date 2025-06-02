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
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\nПроизошла непредвиденная ошибка: %v\n", r)
			waitForEnter()
		}
	}()

	var conn *ssh.Client
	var err error
	var session *ssh.Session

	fmt.Println("Начинаем подключение...")

	for {
		// Получаем данные от пользователя
		host := getDefaultInput("Введите IP-адрес устройства", "192.168.88.32")
		user := getDefaultInput("Введите имя пользователя", "root")
		password := getDefaultInput("Введите пароль", "segnetics")

		fmt.Printf("Пытаемся подключиться к %s как %s...\n", host, user)

		// Пытаемся подключиться
		conn, err = tryConnect(host, user, password)
		if err == nil {
			printSuccess("Подключение успешно установлено!")
			break
		}

		printError(fmt.Sprintf("Ошибка подключения: %v", err))
		retry := getUserInput("Хотите попробовать снова? (y/n): ")
		if strings.ToLower(retry) != "y" {
			waitForEnter()
			return
		}
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	fmt.Println("Выполняем тестовую команду...")

	// Выполняем простую команду для проверки
	cmd := "echo 'test connection'"
	fmt.Printf("Отправляем команду: %s\n", cmd)

	// Создаем новую сессию для команды
	session, err = conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Ошибка создания новой сессии: %v", err))
		waitForEnter()
		return
	}
	defer session.Close()

	fmt.Println("Сессия создана, пытаемся выполнить команду...")

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка выполнения команды: %v", err))
		fmt.Println("Проверяем соединение...")

		// Проверяем, живо ли соединение
		if _, err := conn.NewSession(); err != nil {
			printError("Соединение потеряно. Попытка переподключения...")
			waitForEnter()
			return
		}

		fmt.Println("Соединение активно, но команда не выполняется")
		waitForEnter()
		return
	}

	// Выводим результат
	printSuccess("Команда выполнена успешно. ")
	fmt.Printf("Результат:\n%s\n", output)

	fmt.Println("Теперь пробуем выполнить cat /proc/version...")

	// Создаем новую сессию для следующей команды
	session, err = conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Ошибка создания сессии для cat: %v", err))
		waitForEnter()
		return
	}
	defer session.Close()

	cmd = "cat /proc/version"
	fmt.Printf("Отправляем команду: %s\n", cmd)

	output, err = session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка выполнения команды cat: %v", err))
		waitForEnter()
		return
	}

	printSuccess("Команда cat выполнена успешно.")
	fmt.Printf("Результат:\n%s\n", output)

	fmt.Println("Проверяем существование директории /etc/opt...")

	// Создаем новую сессию для проверки директории
	session, err = conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Ошибка создания сессии для проверки директории: %v", err))
		waitForEnter()
		return
	}
	defer session.Close()

	// Проверяем существование директории
	cmd = "test -d /etc/opt && echo 'Directory exists' || echo 'Directory does not exist'"
	output, err = session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка проверки директории: %v", err))
		waitForEnter()
		return
	}

	if strings.Contains(string(output), "does not exist") {
		printError("Ошибка: директория /etc/opt не существует!")
		waitForEnter()
		return
	}

	printSuccess("Директория /etc/opt существует, продолжаем...")

	// Копируем файлы
	files := []string{"sequencer_v1r6.php", "segnetics.php"}
	for _, file := range files {
		fmt.Printf("\nНачинаем копирование файла %s...\n", file)

		// Проверяем существование локального файла
		if _, err := os.Stat(file); os.IsNotExist(err) {
			printError(fmt.Sprintf("Локальный файл %s не найден", file))
			continue
		}

		// Создаем новую сессию для каждого файла
		session, err = conn.NewSession()
		if err != nil {
			printError(fmt.Sprintf("Не удалось создать сессию для копирования %s: %v", file, err))
			continue
		}
		defer session.Close()

		// Получаем абсолютный путь к локальному файлу
		localPath := filepath.Join(".", file)
		remotePath := "/etc/opt/" + file // Используем Unix-стиль пути

		fmt.Printf("Копируем %s в %s...\n", localPath, remotePath)

		// Копируем файл
		if err := copyFile(session, localPath, remotePath); err != nil {
			printError(fmt.Sprintf("Ошибка копирования файла %s: %v", file, err))
		} else {
			printSuccess(fmt.Sprintf("Файл %s успешно скопирован", file))
		}
	}

	fmt.Println("\nВыводим список файлов в /etc/opt:")

	// Создаем новую сессию для вывода списка файлов
	session, err = conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Ошибка создания сессии для вывода списка файлов: %v", err))
		waitForEnter()
		return
	}
	defer session.Close()

	// Выполняем команду ls
	cmd = "ls -l /etc/opt"
	output, err = session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка получения списка файлов: %v", err))
		waitForEnter()
		return
	}

	fmt.Printf("Содержимое директории /etc/opt:\n%s\n", output)

	fmt.Println("Проверяем существование директории projects...")

	// Создаем новую сессию для проверки директории projects
	session, err = conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Ошибка создания сессии для проверки директории projects: %v", err))
		waitForEnter()
		return
	}
	defer session.Close()

	// Проверяем существование директории projects
	cmd = "test -d /projects && echo 'Directory projects exists' || echo 'Directory projects does not exist'"
	output, err = session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка проверки директории projects: %v", err))
		waitForEnter()
		return
	}

	if strings.Contains(string(output), "does not exist") {
		printError("Ошибка: директория projects не существует!")
		waitForEnter()
		return
	}

	printSuccess("Директория projects существует, продолжаем...")

	fmt.Println("Проверяем существование файла start.after в директории projects...")

	// Создаем новую сессию для проверки файла start.after
	session, err = conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Ошибка создания сессии для проверки файла start.after: %v", err))
		waitForEnter()
		return
	}
	defer session.Close()

	// Проверяем существование файла start.after
	cmd = "test -f /projects/start.after && echo 'File start.after exists' || echo 'File start.after does not exist'"
	output, err = session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка проверки файла start.after: %v", err))
		waitForEnter()
		return
	}

	if strings.Contains(string(output), "does not exist") {
		printError("Файл start.after не существует в директории projects")
	} else {
		printSuccess("Файл start.after существует в директории projects")
	}

	// Предлагаем пользователю удалить файл
	removeFile := getUserInput("Хотите удалить файл start.after? (y/n): ")
	if strings.ToLower(removeFile) == "y" {
		fmt.Println("Пытаемся удалить файл start.after...")

		// Создаем новую сессию для удаления файла
		session, err = conn.NewSession()
		if err != nil {
			printError(fmt.Sprintf("Ошибка создания сессии для удаления файла: %v", err))
			waitForEnter()
			return
		}
		defer session.Close()

		// Удаляем файл
		cmd = "rm -f /projects/start.after"
		_, err = session.CombinedOutput(cmd)
		if err != nil {
			printError(fmt.Sprintf("Ошибка удаления файла start.after: %v", err))
		} else {
			printSuccess("Файл start.after успешно удален")

			// Проверяем, что файл действительно удален
			session, err = conn.NewSession()
			if err != nil {
				printError(fmt.Sprintf("Ошибка создания сессии для проверки удаления: %v", err))
				waitForEnter()
				return
			}
			defer session.Close()

			cmd = "test -f /projects/start.after && echo 'File still exists' || echo 'File successfully removed'"
			output, err = session.CombinedOutput(cmd)
			if err != nil {
				printError(fmt.Sprintf("Ошибка проверки удаления файла: %v", err))
			} else {
				fmt.Printf("Результат проверки: %s", output)
			}
		}
	} else {
		fmt.Println("Удаление файла start.after отменено")
	}

	// Копирование файла start.after из локальной директории на сервер
	fmt.Println("\nПроверяем наличие локального файла start.after...")
	localStartAfter := "start.after"

	if _, err := os.Stat(localStartAfter); os.IsNotExist(err) {
		printError("Локальный файл start.after не найден в директории с программой")
	} else {
		fmt.Println("Локальный файл start.after найден, начинаем копирование...")

		// Создаем новую сессию для копирования
		session, err = conn.NewSession()
		if err != nil {
			printError(fmt.Sprintf("Ошибка создания сессии для копирования: %v", err))
			waitForEnter()
			return
		}
		defer session.Close()

		remotePath := "/projects/start.after"
		fmt.Printf("Копируем %s в %s...\n", localStartAfter, remotePath)

		// Используем нашу функцию copyFile
		if err := copyFile(session, localStartAfter, remotePath); err != nil {
			printError(fmt.Sprintf("Ошибка копирования файла: %v", err))
		} else {
			printSuccess("Файл start.after успешно скопирован на сервер")

			// Проверяем, что файл появился на сервере
			session, err = conn.NewSession()
			if err != nil {
				printError(fmt.Sprintf("Ошибка создания сессии для проверки: %v", err))
				waitForEnter()
				return
			}
			defer session.Close()

			cmd = fmt.Sprintf("ls -la %s", remotePath)
			output, err = session.CombinedOutput(cmd)
			if err != nil {
				printError(fmt.Sprintf("Ошибка проверки файла на сервере: %v", err))
			} else {
				fmt.Printf("Файл на сервере:\n%s\n", output)

				// Выводим первые 10 строк файла для проверки
				session, err = conn.NewSession()
				if err != nil {
					printError(fmt.Sprintf("Ошибка создания сессии для проверки содержимого: %v", err))
					waitForEnter()
					return
				}
				defer session.Close()

				cmd = fmt.Sprintf("head -n 10 %s", remotePath)
				output, err = session.CombinedOutput(cmd)
				if err != nil {
					printError(fmt.Sprintf("Ошибка чтения начала файла: %v", err))
				} else {
					fmt.Printf("Первые 10 строк файла:\n%s\n", output)
				}
			}
		}
	}

	waitForEnter()
}
