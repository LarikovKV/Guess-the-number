package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type GameResult struct { // Структура для хранения результатов игры
	Date       string `json:"date"`
	Outcome    string `json:"outcome"`
	Attempts   int    `json:"attempts"`
	Difficulty string `json:"difficulty"`
}

func generateSecretNumber(maxNumber int) int { // Функция для генерации случайного числа в диапазоне от 1 до maxNumber
	return rand.Intn(maxNumber) + 1
}

func calculateDifference(number int, secretNumber int) int { // Функция для вычисления разницы между введенным числом и загаданным числом
	difference := number - secretNumber

	if difference < 0 { // Если разница отрицательная, то делаем ее положительной
		difference = -difference
	}

	return difference
}

func showTemperatureHint(difference int) { // Функция для вывода подсказки о "температуре" в зависимости от разницы между введенным числом и загаданным
	if difference == 0 {
		return // Если разница равна нулю, то игрок угадал число, и подсказка не нужна
	} else if difference <= 5 {
		fmt.Println("Горячо🔥")
	} else if difference <= 15 {
		fmt.Println("Тепло")
	} else {
		fmt.Println("Холодно❄️")
	}
}

func compareNumbers(number int, secretNumber int) int { // Функция для сравнения введенного числа с загаданным числом
	if number < secretNumber {
		return -1
	} else if number > secretNumber { // Если введенное число больше загаданного, возвращаем 1
		return 1
	}

	return 0
}

func validateNumber(input string, maxNumber int) (int, bool) { // Функция для проверки введенного числа на корректность
	number, err := strconv.Atoi(input)

	if err != nil { // Проверяем на ошибки преобразования
		fmt.Println("Введите целое число.")
		return 0, false
	}

	if number < 1 || number > maxNumber { // Проверка что наше число в нужном диапазоне а не минус 5 или 150
		fmt.Printf("Число должно быть от 1 до %d.\n", maxNumber)
		return 0, false
	}

	return number, true
}

func saveResult(result GameResult) error {
	results := []GameResult{}

	data, err := os.ReadFile("results.json")

	if err == nil {
		err = json.Unmarshal(data, &results)

		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	results = append(results, result)

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("results.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() { // Функция main — точка входа в программу
	scanner := bufio.NewScanner(os.Stdin) // Создаем сканер для чтения ввода пользователя из стандартного потока ввода (stdin)
	fmt.Println("Добро пожаловать в игру Угадай число!")

	for {
		fmt.Println("Выберите уровень сложности:")
		fmt.Println("1 — Easy (1–50, 15 попыток)")
		fmt.Println("2 — Medium (1–100, 10 попыток)")
		fmt.Println("3 — Hard (1–200, 5 попыток)")

		fmt.Print("Ваш выбор (1, 2 или 3): ")

		if !scanner.Scan() {
			fmt.Println("Ввод завершён")
			return
		}

		choice := scanner.Text()

		maxSecretNumber := 100
		maxAttempts := 10 // Max попыток

		switch choice { // Выбор уровня сложности
		case "1":
			maxSecretNumber = 50
			maxAttempts = 15
			fmt.Println("Вы выбрали Easy")
		case "2":
			maxSecretNumber = 100
			maxAttempts = 10
			fmt.Println("Вы выбрали Medium")
		case "3":
			maxSecretNumber = 200
			maxAttempts = 5
			fmt.Println("Вы выбрали Hard")
		default: // Если введено неверное значение, выводим сообщение и завершаем программу
			fmt.Println("Неверный выбор сложности")
			return // Завершаем программу, если введено неверное значение
		}

		secretNumber := generateSecretNumber(maxSecretNumber) // Генерируем случайное число в зависимости от выбранного уровня сложности
		fmt.Printf("Я загадал число от 1 до %d. У вас %d попыток!\n", maxSecretNumber, maxAttempts)

		win := false       // Игрок пока не угадал число
		attemptsUsed := 0  // Количество использованных попыток
		history := []int{} // Добавил слайс для хранения истории введенных чисел

		for attempt := 1; attempt <= maxAttempts; attempt++ { // Цикл до maxAttempts попыток
			fmt.Printf("\033[33mПопытка #%d из %d\033[0m — Введите число: ", attempt, maxAttempts) // Выводим номер текущей попытки и общее количество попыток с цветом ANSI код

			if !scanner.Scan() { //проверка с завершением
				fmt.Println("Ввод завершён. Выходим из игры.")
				return
			}

			input := scanner.Text()

			number, valid := validateNumber(input, maxSecretNumber)

			if !valid {
				attempt--
				continue
			}

			fmt.Println("Вы ввели число:", number)
			history = append(history, number)     // Добавляем число расширяем историю попыток
			attemptsUsed = len(history)           // Количество потраченных попыток
			fmt.Println("Ваши попытки:", history) // Показываем историю

			difference := calculateDifference(number, secretNumber) // Вычисляем разницу между введенным числом и загаданным
			showTemperatureHint(difference)                         // Выводим подсказку о "температуре" в зависимости от разницы между введенным числом и загаданным

			comparison := compareNumbers(number, secretNumber)

			switch comparison {
			case -1:
				fmt.Println("Секретное число больше👆")
			case 1:
				fmt.Println("Секретное число меньше👇")
			case 0:
				win = true
				fmt.Println("\033[32mВы угадали! Поздравляю!🎉\033[0m") // Добавляем цвет ANSI код

			}

			if win {
				break
			}

		}
		if !win { // Если игрок исчерпал все попытки и не угадал число
			fmt.Println("\033[31mВы не угадали. Попробуйте еще раз.\033[0m")

		}
		fmt.Println("Игра закончена!")
		fmt.Println("Секретное число было:", secretNumber)

		outcome := "loss" // В случае проигрыша

		if win { // Если игрок угадал число, то меняем исход на "win"
			outcome = "win"
		}

		result := GameResult{ // Создаем структуру с результатами игры
			Date:       time.Now().Format(time.RFC3339), // Записываем текущую дату и время
			Outcome:    outcome,
			Attempts:   attemptsUsed,
			Difficulty: choice,
		}

		fmt.Printf("Результат игры: %+v\n", result) // Для проверки вывожу результат игры в консоль

		err := saveResult(result) // Сохраняем результат игры в файл results.json

		if err != nil {
			fmt.Println("Ошибка сохранения результата:", err)
			return
		}

		fmt.Println("Результат сохранён в results.json")

		fmt.Print("Сыграть ещё раз? (да/нет): ") // Спрашиваем будем еще продолжать или нет

		if !scanner.Scan() {
			fmt.Println("Ввод завершён")
			return // Завершаем программу, если ввод завершён
		}

		answer := scanner.Text() // Запоминаем ответ пользователя
		fmt.Println("Ваш ответ:", answer)

		if answer == "нет" { // Если пользователь ввел "нет", то выводим сообщение и завершаем программу
			fmt.Println("Спасибо за игру")
			return
		}

		if answer == "да" { // Если пользователь ввел "да", то выводим сообщение и начинаем новую игру
			fmt.Println("Начинаем новую игру")
		}
	}
}
