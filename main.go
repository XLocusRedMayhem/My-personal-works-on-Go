package main

import "fmt"

// 2026-10-03
// Напишите программу, которая в последовательности чисел находит сумму
// двузначных чисел, кратных 8. Программа в первой строке получает на
// вход число n - количество чисел в последовательности, во второй
// строке -- n чисел, входящих в данную последовательность.
func main() {
	for i := 0; i == 0; {
		var choice int
		fmt.Print("0. Выход.\n1. Вывод квадратов натуральных чисел от 1 до 10.\n2. Вывод второй цифры справа у числа.\n3. Различие и сходство цифр трёхзначного числа.\n4. Сумма чисел кратных 8.\nВыберите вариант: ")
		fmt.Scan(&choice)
		fmt.Print("\n")
		switch choice {
		case 0:
			i = 1
		case 1:
			Task_1()
		case 2:
			fmt.Print("Вводимое число должно быть от 10 до 10000.\n")
			Task_2()
		case 3:
			Task_3()
		case 4:
			Task_4()
		case 5:
			Task_5()
		}
	}
}

func Task_1() {
	for i := 1; i <= 10; i++ {
		fmt.Print(i, "*", i, "=", i*i, "\n")
	}
	fmt.Print("\n")
}

func Task_2() {
	var num int
	fmt.Print("Введите число: ")
	fmt.Scan(&num)
	if num > 10000 {
		fmt.Print("Число больше 10000\n\n")
	}
	if num < 0 {
		fmt.Print("Число отрицательное\n\n")
	}
	if num <= 10000 && num >= 0 {
		fmt.Print("Ответ: ", num%100/10, "\n\n")
	}
}

func Task_3() {
	var num int
	fmt.Print("Введите число: ")
	fmt.Scan(&num)
	var (
		num0 = num / 100
		num1 = num % 100 / 10
		num2 = num % 10
	)
	fmt.Print("Ответ: ")
	if num0 == num1 || num0 == num2 || num1 == num2 {
		fmt.Print("NO\n\n")
	} else {
		fmt.Print("YES\n\n")
	}
}

func Task_4() {
	var amountNumber int
	var number int
	var sumNambers int = 0
	fmt.Print("Введите количество чисел: ")
	fmt.Scan(&amountNumber)
	if amountNumber <= 0 {
		fmt.Println("Количетсво чисел не может быть 0 или отрицательным.")
	} else {
		fmt.Print("Введите числа через пробел: ")
		for i := 0; i < amountNumber; i++ {
			fmt.Scan(&number)
			if number < 10 && number > 99 {
				i--
				number = 0
			} else {
				if number%8 == 0 {
					sumNambers += number
				}
			}
		}
		fmt.Print("Сумма чисел: ", sumNambers, "\n\n")
	}
}

func Task_5() {

}
