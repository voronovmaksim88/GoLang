package main

import (
	"fmt"
)

func calculateBMI(weight float64, height float64) float64 {
	return weight / (height * height)
}

func main() {
	var weight float64 = 92   // вес в килограммах
	var height float64 = 1.78 // рост в метрах

	bmi := calculateBMI(weight, height)
	fmt.Printf("Ваш индекс массы тела равен %.2f\n", bmi)
}
