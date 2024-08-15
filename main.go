package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var modeRome = false
var modeArab = false

// Карта для перевода римских цифр в арабские
var romeToArab = map[string]int{
	"I": 1, "II": 2, "III": 3, "IV": 4, "V": 5, "VI": 6, "VII": 7, "VIII": 8, "IX": 9,
	"X": 10,
}

// Карта для арабских цифр в римские срез наполненный структурами
var arabToRome = []struct {
	dec    int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Введите значение (Калькулятор умеет выполнять операции (+, -, *, /) с двумя числами)")
		text, _ := reader.ReadString('\n')         //Ждём ввода данных в формате строки
		text = strings.Replace(text, "\n", "", -1) //Очищаем все пустоты пробелы табуляции
		arrValues := strings.Split(text, " ")
		if len(arrValues) == 3 {
			calculations(arrValues[0], arrValues[1], arrValues[2])
		} else {
			sentPanic("Не корректный ввод!")
		}

	}
}

func calculations(val1 string, sign string, val2 string) {
	var next = make([]bool, 2)
	var parsVal1, parsVal2 int
	var err error
	if romeToArab[val1] != 0 {
		modeRome = true
		parsVal1 = romeToArab[val1]
		next[0] = validationNumber(romeToArab[val1], val1, nil)
	} else {
		modeArab = true
		parsVal1, err = strconv.Atoi(val1) //Преобразуем строку в число
		next[0] = validationNumber(parsVal1, val1, err)
	}
	if romeToArab[val2] != 0 {
		modeRome = true
		parsVal2 = romeToArab[val2]
		next[1] = validationNumber(romeToArab[val2], val2, nil)
	} else {
		modeArab = true
		parsVal2, err = strconv.Atoi(val2) //Преобразуем строку в число
		next[1] = validationNumber(parsVal2, val2, err)
	}
	if modeArab && modeRome {
		sentPanic("Калькулятор умеет работать только с арабскими или римскими цифрами одновременно")
	}
	if next[0] && next[1] {
		switch sign {
		case "+":
			output(parsVal1 + parsVal2)
		case "-":
			output(parsVal1 - parsVal2)
		case "*":
			output(parsVal1 * parsVal2)
		case "/":
			output(parsVal1 / parsVal2)
		default:
			sentPanic("Арифметическая операция не опознана!")
		}
	}

}

func output(val int) {
	if modeRome {
		fmt.Println(transformToRome(val))
	} else {
		if val > 0 {
			fmt.Println(val)
		} else {
			sentPanic("Результатом работы калькулятора с арабскими числами могут быть отрицательные числа и ноль. Результатом работы калькулятора с римскими числами могут быть только положительные числа")
		}
	}
}

func validationNumber(num int, original string, err error) bool {
	var result bool = true
	if num%1 != 0 {
		result = false
		sendError(original, "%w Калькулятор умеет работать только с целыми числами. Вы ввели (%s)\n")
	} else if err != nil {
		result = false
		sentPanic("Не подходящее число!")
	} else if num > 10 {
		result = false
		sendError(original, "%w допустимое максимальное число 10 вы ввели (%s)\n")
	}
	return result
}

func sendError(val string, message string) {
	if message == "" {
		message = "%w Вы ввели некорректное значение (%s)\n"
	}
	err := errors.New("Возникла ошибка:")
	err = fmt.Errorf(message, err, val)
	fmt.Print(err)
}

func sentPanic(message string) {
	if message == "" {
		message = "Допустим ввод только Римских или Арабских значений."
	}
	err := errors.New(message)
	panic(err)
}
func transformToRome(num int) string {
	result := ""
	for _, item := range arabToRome {
		for num >= item.dec {
			result += item.symbol
			num -= item.dec
		}
	}
	return result
}
