package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	// Получаем данные от пользователя
	host := getUserInput("Введите IP-адрес устройства (например, 192.168.1.1): ")
	user := getUserInput("Введите имя пользователя: ")
	password := getUserInput("Введите пароль: ")

	// Формируем адрес подключения
	addr := host + ":22"

	// Конфигурация SSH клиента
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Не проверять известные хосты
	}

	// Подключаемся
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
		return
	}
	defer conn.Close()

	// Создаём сессию
	session, err := conn.NewSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось создать сессию: %v\n", err)
		return
	}
	defer session.Close()

	// Выполняем команду
	cmd := "cat /proc/version"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения команды: %v\n", err)
		return
	}

	// Выводим результат
	fmt.Printf("Результат выполнения '%s':\n%s\n", cmd, output)

	// Ожидаем нажатия Enter
	fmt.Print("Нажмите Enter для выхода...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
