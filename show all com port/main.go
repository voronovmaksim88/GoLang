package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/tarm/serial"
	"log"
	"os"
)

func main() {
	// Отложенный вызов ожидания Enter с обработкой ошибки
	defer func() {
		log.Println("Нажмите Enter для выхода...")
		_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Printf("Ошибка при чтении ввода: %v", err)
		}
	}()

	// Создаем функции для вывода текста зеленым и красным цветом
	green := color.New(color.FgGreen).PrintfFunc()
	red := color.New(color.FgRed).PrintfFunc()

	// Список возможных COM-портов для проверки
	ports := []string{
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8",
		"COM9", "COM10", "COM11", "COM12", "COM13", "COM14", "COM15",
	}

	fmt.Println("Доступные COM-порты:")

	// Счетчик доступных портов
	availablePorts := 0

	for _, port := range ports {
		// Пытаемся открыть порт
		config := &serial.Config{Name: port, Baud: 9600}
		s, err := serial.OpenPort(config)
		if err == nil {
			// Если порт открыт успешно, выводим его имя зеленым цветом
			green("%s\n", port)
			availablePorts++
			// Закрываем порт и обрабатываем ошибку
			if err := s.Close(); err != nil {
				fmt.Printf("Ошибка при закрытии порта %s: %v\n", port, err)
			}
		}
	}

	// Если нет доступных портов, выводим сообщение красным цветом
	if availablePorts == 0 {
		red("Нет доступных COM-портов.\n")
	}
}
