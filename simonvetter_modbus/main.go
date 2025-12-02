package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/simonvetter/modbus"
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
		"BUILD_Y\t 0x0109\t uint16\t только_чтение\t год\n" +
		"BUILD_M\t 0x010A\t uint16\t только_чтение\t месяц\n" +
		"BUILD_D\t 0x010B\t uint16\t только_чтение\t день\n" +
		"BUILD_H\t 0x010C\t uint16\t только_чтение\t час\n" +
		"BUILD_m\t 0x010D\t uint16\t только_чтение\t минута\n" +
		"BUILD_S\t 0x010E\t uint16\t только_чтение\t секунда")

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
