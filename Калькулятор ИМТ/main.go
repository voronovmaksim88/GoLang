package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Функция расчета индекса массы тела
func calculateBMI(weight float64, height float64) float64 {
	return weight / (height * height)
}

// Главная функция программы
func main() {
	reader := bufio.NewReader(os.Stdin)

	var weight, height float64
	var err error

	for {
		fmt.Print("Введите ваш вес (кг): ")
		inputWeight, _ := reader.ReadString('\n')
		weight, err = strconv.ParseFloat(strings.TrimSpace(inputWeight), 64)
		if err == nil {
			break
		}
		fmt.Println("Некорректный ввод веса. Повторите попытку.")
	}

	for {
		fmt.Print("Введите ваш рост (метры): ")
		inputHeight, _ := reader.ReadString('\n')
		height, err = strconv.ParseFloat(strings.TrimSpace(inputHeight), 64)
		if err == nil && height > 0 {
			break
		}
		fmt.Println("Некорректный ввод роста. Рост должен быть положительным числом.")
	}

	// Вычисляем ИМТ
	bmi := calculateBMI(weight, height)

	// Интерпретация результата
	switch {
	case bmi < 18.5:
		fmt.Printf("Ваш ИМТ равен %.2f. У вас недостаточный вес.\n", bmi)
	case bmi >= 18.5 && bmi <= 24.9:
		fmt.Printf("Ваш ИМТ равен %.2f. Ваш вес нормальный.\n", bmi)
	case bmi > 24.9 && bmi <= 29.9:
		fmt.Printf("Ваш ИМТ равен %.2f. У вас избыточный вес.\n", bmi)
	default:
		fmt.Printf("Ваш ИМТ равен %.2f. У вас ожирение.\n", bmi)
	}

	// Ожидание нажатия Enter перед закрытием окна
	fmt.Println("\nНажмите Enter для выхода...")
	_, _ = reader.ReadString('\n')
}
