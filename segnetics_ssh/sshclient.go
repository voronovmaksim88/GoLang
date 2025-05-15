package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

func printError(message string) {
	color.Red(message)
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

func main() {
	var conn *ssh.Client
	var err error

	for {
		// Получаем данные от пользователя
		host := getUserInput("Введите IP-адрес устройства (например, 192.168.111.2): ")
		user := getDefaultInput("Введите имя пользователя", "root")
		password := getDefaultInput("Введите пароль", "segnetics")

		// Пытаемся подключиться
		conn, err = tryConnect(host, user, password)
		if err == nil {
			break
		}

		printError(fmt.Sprintf("Ошибка подключения: %v", err))
		retry := getUserInput("Хотите попробовать снова? (y/n): ")
		if strings.ToLower(retry) != "y" {
			return
		}
	}
	defer conn.Close()

	// Создаём сессию
	session, err := conn.NewSession()
	if err != nil {
		printError(fmt.Sprintf("Не удалось создать сессию: %v", err))
		return
	}
	defer session.Close()

	// Выполняем команду
	cmd := "cat /proc/version"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		printError(fmt.Sprintf("Ошибка выполнения команды: %v", err))
		return
	}

	// Выводим результат
	fmt.Printf("Результат выполнения '%s':\n%s\n", cmd, output)

	// Ожидаем нажатия Enter
	fmt.Print("Нажмите Enter для выхода...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
