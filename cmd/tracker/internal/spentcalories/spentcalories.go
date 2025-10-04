package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Реализация 3 функции

	// 1) Разделяем строку на слайс строк по запятым
	packagePart := strings.Split(data, ",")

	// 2) Проверяем длину слайса (должна равняться 3)
	if len(packagePart) != 3 {
		return 0, "", 0, errors.New("некорректные (неполные) данные")
	}

	// 3) Преобразовываем первый элемент слайса (количество шагов) в int
	stepsStr := strings.TrimSpace(packagePart[0]) // удаляем пробелы по краям
	if stepsStr == "" {
		return 0, "", 0, errors.New("некорректно указано количество шагов")
	}

	// обработка знака "+" и его удаление
	if after, ok := strings.CutPrefix(stepsStr, "+"); ok {
		stepsStr = after
	}

	// преобразуем строку в целое число
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("некорректно указано количество шагов: %v", err)
	}

	// проверяем что количество шагов больше 0
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	// 4) Преобразовываем третий элемент слайса в time.Duration
	durationStr := strings.TrimSpace(packagePart[2])
	if durationStr == "" {
		return 0, "", 0, errors.New("некорректно указана продолжительность")
	}

	// забираем продолжительность используя time
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("некорректно указана продолжительность: %v", err)
	}

	// проверяем что продолжительность больше 0
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}

	// 5) Определить вид активности
	activityType := strings.TrimSpace(packagePart[1])

	// 6) Возврат значений
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// Реализация 4 функции

	// 1) Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient

	// 2) Вычисляем общую дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 3) Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Реализация 5 функции

	// 1) Проверяем что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// 2) Вычисляем дистанцию через функцию distance()
	dist := distance(steps, height)

	// 3) Вычисляем среднюю скорость
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	speed := dist / hours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Реализация 6 функции

	// 1) Получаем значения из строки данных через parseTraining() с выводом и обработкой ошибок log.Println(err)
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка получения данных тренировки:", err)
		return "", err
	}

	// 2) Проверяем, какой вид тренировки был передан в строке
	var calories float64
	var caloriesErr error

	// приводим тип активности к нижнему регистру для унификации сравнения
	switch strings.ToLower(activityType) {
	case "ходьба", "walking":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	case "бег", "running":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки: " + activityType)
	}

	// проверяем ошибку расчета калорий
	if caloriesErr != nil {
		log.Println("Ошибка расчета калорий:", caloriesErr)
		return "", caloriesErr
	}

	// 3) Рассчитываем дистанцию и среднюю скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// 4) Сформировываем строку с информацией о тренировке
	info := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		dist,
		speed,
		calories)
	return info, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Реализация 7 функции

	// 1) Проверка входных параметров на корректность
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// 2) Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed == 0 {
		return 0, errors.New("средняя скорость равна 0")
	}

	// 3) Рассчитываем количество калорий
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Реализация 8 функции

	// 1) Проверка входных параметров на корректность
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// 2) Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed == 0 {
		return 0, errors.New("средняя скорость равна 0")
	}

	// 3) Рассчитываем количество калорий
	durationInMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationInMinutes) / minInH

	// 4) Учитываем корректирующий коэффициент
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}
