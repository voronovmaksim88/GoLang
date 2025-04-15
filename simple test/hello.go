package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// Создаем читатель для ввода с клавиатуры
	reader := bufio.NewReader(os.Stdin)

	// Запрашиваем имя пользователя
	fmt.Print("Введите ваше имя: ")
	name, _ := reader.ReadString('\n') // Игнорируем ошибку здесь, так как она маловероятна при вводе с клавиатуры
	name = strings.TrimSpace(name)     // Убираем лишние пробелы и символы новой строки

	// Выводим приветствие
	fmt.Printf("Привет, %s! Добро пожаловать в мир Go!\n", name)

	// Дополнительно: простая проверка длины имени
	if len(name) > 0 {
		// Создаем зеленый принтер
		green := color.New(color.FgGreen).SprintFunc()
		// Выводим сообщение зеленым цветом
		fmt.Printf("%s\n", green(fmt.Sprintf("Ваше имя состоит из: '%d' символов:", len([]rune(name)))))
	} else {
		fmt.Println("Вы не ввели имя")
	}

	// Ожидаем нажатия Enter перед завершением
	fmt.Println("\nНажмите Enter для выхода...")
	_, _ = reader.ReadString('\n') // Явно игнорируем оба возвращаемых значение
}