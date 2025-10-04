package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Реализация 1 функции

	// 1) Разделяем строку на слайс строк через запятую
	packagePart := strings.Split(data, ",")

	// 2) Проверяем длину слайса (должна равняться 2)
	if len(packagePart) != 2 {
		return 0, 0, errors.New("некорректные (неполные) данные")
	}

	// 3) Преобразовываем первый элемент слайса (количество шагов) в int, с обработкой на ошибки
	stepsStr := packagePart[0]
	if strings.ContainsAny(stepsStr, " \t\n\r") {
		return 0, 0, errors.New("некорректно указано количество шагов")
	}

	// Обработка знака "+"
	if after, ok := strings.CutPrefix(stepsStr, "+"); ok {
		stepsStr = after
	}

	if stepsStr == "" {
		return 0, 0, errors.New("некорректно указано количество шагов")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("некорректно указано количество шагов: %v", err)
	}

	// 4) Проверка количества шагов на положительное число
	if steps <= 0 {
		return 0, 0, errors.New("указано отрицательное количество шагов")
	}

	// 5) Преобразовываем второй элемент слайса (время) в time.Duration с обработкой ошибок
	durationStr := packagePart[1]
	if strings.ContainsAny(durationStr, " \t\n\r") {
		return 0, 0, errors.New("некорректно указана продолжительность")
	}

	if durationStr == "" {
		return 0, 0, errors.New("некорректно указана продолжительность")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("некорректно указана продолжительность: %v", err)
	}

	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}
	// 6) Возврат значений и nil(для ошибки)
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Реализация 2 функции

	// 1) Получаем данные о количестве шагов и продолжительности, с учетом ошибок
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка получения данных:", err)
		return ""
	}

	// 3) Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4) Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// 5) Вычисляем количество потраченных калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка расчета калорий:", err)
		return ""
	}

	// 6) Сформировываем строку для возврата
	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

	return result
}
