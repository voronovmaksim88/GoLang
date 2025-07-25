package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/simonvetter/modbus"
	"log"
	"os"
	"time"
)

func runModbus() (err error) {
	// Отложенное восстановление после паники
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("паника: %v", r)
		}
	}()

	cfg := &modbus.ClientConfiguration{
		URL:      "rtu://COM7", // для windows
		Speed:    9600,
		DataBits: 8,
		Parity:   modbus.PARITY_EVEN,
		StopBits: 2,
		Timeout:  5 * time.Second,
	}

	client, err := modbus.NewClient(cfg)
	if err != nil {
		return err
	}

	// Отложенное закрытие соединения с обработкой ошибки
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Printf("Ошибка при закрытии соединения: %v", closeErr)
			}
		}
	}()

	// Открываем соединение
	err = client.Open()
	if err != nil {
		return err
	}

	err = client.SetUnitId(1)
	if err != nil {
		return err
	}

	regs, err := client.ReadRegisters(0x0109, 6, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}
	log.Printf("Результат: %v", regs)
	return nil
}

func main() {
	log.Println("Привет друг. Сейчас будем читать регистры из комбайна\n" +
		"BUILD_Y\t0x0109\tuint16\tтолько чтение\tгод\n" +
		"BUILD_M\t0x010A\tuint16\tтолько чтение\tмесяц\n" +
		"BUILD_D\t0x010B\tuint16\tтолько чтение\tдень\n" +
		"BUILD_H\t0x010C\tuint16\tтолько чтение\tчас\n" +
		"BUILD_m\t0x010D\tuint16\tтолько чтение\tминута\n" +
		"BUILD_S\t0x010E\tuint16\tтолько чтение\tсекунда")

	// Отложенный вызов ожидания Enter с обработкой ошибки
	defer func() {
		log.Println("Нажмите Enter для выхода...")
		_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Printf("Ошибка при чтении ввода: %v", err)
		}
	}()

	// Выполнение Modbus-запроса
	if err := runModbus(); err != nil {
		// Вывод ошибки красным цветом
		red := color.New(color.FgRed).SprintFunc()
		log.Printf("%s: %v", red("Ошибка"), err)
	}
}
