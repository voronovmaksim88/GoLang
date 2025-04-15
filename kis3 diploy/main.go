package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

func main() {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	// Определение адреса и учетных данных SSH-сервера
	host := "176.124.213.202"
	user := "root"

	// Настройка аутентификации по SSH
	key, err := os.ReadFile("C:\\Users\\Maksim\\.ssh\\id_rsa")
	if err != nil {
		log.Fatalf("Невозможно прочитать приватный ключ: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Невозможно распарсить приватный ключ: %v", err)
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
		log.Fatalf("Невозможно подключиться к SSH-серверу: %v", err)
	}
	defer client.Close()

	fmt.Println("Успешное подключение к SSH-серверу")

	// Выполняем последовательно все необходимые команды

	// 1. Обновляем репозиторий через git pull
	fmt.Println(yellow("\n=== Выполнение git pull ==="))
	executeCommand(client, "cd kis3_v2r2/ && git pull")

	// 2. Собираем контейнеры с помощью docker-compose build
	fmt.Println(yellow("\n=== Выполнение docker-compose build ==="))
	executeCommand(client, "cd kis3_v2r2/ && docker-compose build")

	// 3. Запускаем контейнеры с помощью docker-compose up -d
	fmt.Println(yellow("\n=== Выполнение docker-compose up -d ==="))
	executeCommand(client, "cd kis3_v2r2/ && docker-compose up -d")

	// Создаем читатель для ввода с клавиатуры
	reader := bufio.NewReader(os.Stdin)

	// Ожидаем нажатия Enter перед завершением
	fmt.Println(green("\nВсе команды выполнены. Нажмите Enter для выхода..."))
	_, _ = reader.ReadString('\n')
}

// Вспомогательная функция для выполнения SSH команд
func executeCommand(client *ssh.Client, command string) {
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Невозможно создать SSH-сессию: %v", err)
	}
	defer session.Close()

	// Настройка стандартного вывода и ошибок
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(command)
	if err != nil {
		log.Fatalf("Ошибка выполнения команды '%s': %v", command, err)
	}
}
