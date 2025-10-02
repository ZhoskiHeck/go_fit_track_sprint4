package daysteps

import (
	"errors"
	"fmt"
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
	steps, err := strconv.Atoi(packagePart[0])
	if err != nil {
		return 0, 0, fmt.Errorf("некорректно указано количество шагов: %v", err)
	}

	// 4) Проверка количества шагов на положительное число
	if steps <= 0 {
		return 0, 0, errors.New("указано отрицательное количество шагов")
	}

	// 5) Преобразовываем второй элемент слайса (время) в time.Duration с обработкой ошибок
	duration, err := time.ParseDuration(packagePart[1])
	if err != nil {
		return 0, 0, fmt.Errorf("некорректно указана продолжительность: %v", err)
	}

	// 6) Возврат значений и nil(для ошибки)
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Реализация 2 функции

	// 1) Получаем данные о количестве шагов и продолжительности, с учетом ошибок
	steps, duration, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("Ошибка получения данных: %v", err)
	}

	// 2) Проверяем чтобы количество шагов было больше 0
	if steps <= 0 {
		return "Количество шагов должно быть больше 0"
	}

	// 3) Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4) Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// 5) Вычисляем количество потраченных калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return fmt.Sprintf("Ошибка расчета калорий: %v", err)
	}

	// 6) Сформировываем строку для возврата
	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

	return result
}
