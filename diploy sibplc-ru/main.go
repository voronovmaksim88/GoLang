package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

func main() {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Определение адреса и учетных данных SSH-сервера
	host := "176.124.213.202"
	user := "root"

	// Настройка аутентификации по SSH
	key, err := os.ReadFile("C:\\Users\\Maksim\\.ssh\\id_rsa")
	if err != nil {
		fatalError(red(fmt.Sprintf("Невозможно прочитать приватный ключ: %v", err)))
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fatalError(red(fmt.Sprintf("Невозможно распарсить приватный ключ: %v", err)))
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Подключение к серверу SSH
	addr := fmt.Sprintf("%s:22", host)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fatalError(red(fmt.Sprintf("Невозможно подключиться к SSH-серверу: %v", err)))
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Ошибка при закрытии SSH-клиента: %v", err)
		}
	}()

	fmt.Println(green("Успешное подключение к SSH-серверу"))

	// Выполняем последовательно все необходимые команды

	// 1. Обновляем репозиторий через git pull
	fmt.Println(yellow("\n=== Выполнение git pull ==="))
	if err := executeCommand(client, "cd sibplc_ru/sibplc.ru/ && git pull"); err != nil {
		fmt.Println(red(fmt.Sprintf("Ошибка выполнения git pull: %v", err)))
	}

	// 2. Собираем контейнеры с помощью docker-compose build
	fmt.Println(yellow("\n=== Выполнение docker-compose build ==="))
	if err := executeCommand(client, "cd sibplc_ru/sibplc.ru/ && docker-compose build"); err != nil {
		fmt.Println(red(fmt.Sprintf("Ошибка выполнения docker-compose build: %v", err)))
	}

	// 3. Запускаем контейнеры с помощью docker-compose up -d
	fmt.Println(yellow("\n=== Выполнение docker-compose up -d ==="))
	if err := executeCommand(client, "cd sibplc_ru/sibplc.ru/ && docker-compose up -d"); err != nil {
		fmt.Println(red(fmt.Sprintf("Ошибка выполнения docker-compose up -d: %v", err)))
	}

	// Ожидаем нажатия Enter перед завершением (только при успешном выполнении)
	fmt.Println(green("\nВсе команды выполнены. Нажмите Enter для выхода..."))
	waitForEnter()
}

// Функция для вывода ошибки и ожидания Enter перед выходом
func fatalError(message string) {
	fmt.Println(message)
	waitForEnter()
	os.Exit(1)
}

// Ожидание нажатия Enter
func waitForEnter() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

// Вспомогательная функция для выполнения SSH команд
func executeCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("невозможно создать SSH-сессию: %v", err)
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil && closeErr != io.EOF {
			log.Printf("Ошибка при закрытии SSH-сессии: %v", closeErr)
		}
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session.Run(command)
}
