package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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
		fmt.Println(red(fmt.Sprintf("Невозможно прочитать приватный ключ: %v", err)))
		fmt.Println("Нажмите Enter для завершения...")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return
	}

	// Сначала пробуем без пароля
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		// Если не получилось, запрашиваем пароль
		fmt.Print("Введите пароль для SSH ключа: ")
		reader := bufio.NewReader(os.Stdin)
		passphrase, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(red(fmt.Sprintf("Ошибка чтения пароля: %v", err)))
			fmt.Println("Нажмите Enter для завершения...")
			_, _ = reader.ReadString('\n')
			return
		}
		passphrase = strings.TrimSpace(passphrase)

		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
		if err != nil {
			fmt.Println(red(fmt.Sprintf("Невозможно распарсить приватный ключ: %v", err)))
			fmt.Println("Нажмите Enter для завершения...")
			_, _ = reader.ReadString('\n')
			return
		}
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
		fmt.Println(red(fmt.Sprintf("Невозможно подключиться к SSH-серверу: %v", err)))
		fmt.Println("Нажмите Enter для завершения...")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Ошибка при закрытии SSH-клиента: %v", err)
		}
	}()

	fmt.Println("Успешное подключение к SSH-серверу")

	// Выполняем последовательно все необходимые команды

	// 1. Обновляем репозиторий через git pull
	fmt.Println(yellow("\n=== Выполнение git pull ==="))
	if err := executeCommand(client, "cd SibTask6/ && git pull"); err != nil {
		fmt.Println(red(fmt.Sprintf("Ошибка выполнения git pull: %v", err)))
	}

	// 2. Собираем контейнеры с помощью docker-compose build
	fmt.Println(yellow("\n=== Выполнение docker-compose build ==="))
	if err := executeCommand(client, "cd SibTask6/ && docker-compose build"); err != nil {
		fmt.Println(red(fmt.Sprintf("Ошибка выполнения docker-compose build: %v", err)))
		// Продолжаем выполнение даже при ошибке
	}

	// 3. Запускаем контейнеры с помощью docker-compose up -d
	fmt.Println(yellow("\n=== Выполнение docker-compose up -d ==="))
	if err := executeCommand(client, "cd SibTask6/ && docker-compose up -d"); err != nil {
		fmt.Println(red(fmt.Sprintf("Ошибка выполнения docker-compose up -d: %v", err)))
	}

	// Создаем читатель для ввода с клавиатуры
	reader := bufio.NewReader(os.Stdin)

	// Ожидаем нажатия Enter перед завершением
	fmt.Println(green("\nВсе команды выполнены. Нажмите Enter для выхода..."))
	_, _ = reader.ReadString('\n')
}

// Вспомогательная функция для выполнения SSH команд
func executeCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("невозможно создать SSH-сессию: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("Ошибка при закрытии SSH-сессии: %v", err)
		}
	}()

	// Настройка стандартного вывода и ошибок
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session.Run(command)
}
