package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
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

// Структура для перевода арабских в римские
var arabToRome = []struct {
	dec    int
	symbol string
}{
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	regex := regexp.MustCompile(`^[0-9AIXCLV]{1,4}[*+-/][0-9AIXCLV]{1,4}$`)
	for {
		//Сбрасываем режимы на каждом вводе
		modeRome = false
		modeArab = false

		fmt.Println("Введите значение (Калькулятор умеет выполнять операции (+, -, *, /) с двумя числами)")
		text, _ := reader.ReadString('\n') //Ждём ввода данных в формате строки
		text = strings.TrimRight(text, "\n\r")
		text = strings.Replace(text, " ", "", -1)
		var separator, success = searchSeparator(text)

		if regex.MatchString(text) && success {
			arrValues := strings.Split(text, separator)
			calculations(arrValues[0], separator, arrValues[1])
		} else {
			success = false
		}
		if !success {
			sentPanic("Выдача паники, так как формат математической операции не удовлетворяет заданию — два операнда и один оператор (+, -, /, *).")
		}

	}
}

func searchSeparator(data string) (string, bool) {
	if strings.Contains(data, "*") {
		return "*", true
	} else if strings.Contains(data, "/") {
		return "/", true
	} else if strings.Contains(data, "+") {
		return "+", true
	} else if strings.Contains(data, "-") {
		return "-", true
	} else {
		return "", false
	}
}

func calculations(val1 string, sign string, val2 string) {
	var next = make([]bool, 2)
	regexRome := regexp.MustCompile(`^[AIXCLV]{1,4}$`)
	var parsVal1, parsVal2 int
	var err error
	if regexRome.MatchString(val1) && romeToArab[val1] != 0 {
		modeRome = true
		parsVal1 = romeToArab[val1]
		next[0] = validationNumber(romeToArab[val1], val1, nil)
	} else {
		modeArab = true
		parsVal1, err = strconv.Atoi(val1) //Преобразуем строку в число
		next[0] = validationNumber(parsVal1, val1, err)
	}
	if regexRome.MatchString(val2) && romeToArab[val2] != 0 {
		modeRome = true
		parsVal2 = romeToArab[val2]
		next[1] = validationNumber(romeToArab[val2], val2, nil)
	} else {
		modeArab = true
		parsVal2, err = strconv.Atoi(val2) //Преобразуем строку в число
		next[1] = validationNumber(parsVal2, val2, err)
	}
	if modeArab && modeRome {
		sentPanic("Выдача паники, так как используются одновременно разные системы счисления.")
	}
	if next[0] && next[1] {
		switch sign {
		case "+":
			output(parsVal1 + parsVal2)
		case "-":
			val := parsVal1 - parsVal2
			if modeRome {
				if val > 0 {
					output(val)
				} else {
					sentPanic("Выдача паники, так как в римской системе нет отрицательных чисел.")
				}
			} else {
				output(val)
			}
		case "*":
			output(parsVal1 * parsVal2)
		case "/":
			output(parsVal1 / parsVal2)
		default:
			sentPanic("Выдача паники, так как строка не является математической операцией.")
		}
	}

}

func output(val int) {
	if modeRome {
		fmt.Println(transformToRome(val))
	} else {
		fmt.Println(val)
	}
}

func validationNumber(num int, original string, err error) bool {
	var result bool = true
	if num%1 != 0 {
		result = false
		sendError(original, "%w Калькулятор умеет работать только с целыми числами. Вы ввели (%s)\n")
	} else if err != nil {
		result = false
		sentPanic("Паника не подходящее число!")
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
	err := errors.New("возникла ошибка ")
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
