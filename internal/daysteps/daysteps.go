package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Реализация 1 функции

	// 1) Разделяем строки на слайс строк через запятую
	packagePart := strings.Split(data, ",")

	// 2) Проверяем длину слайса (должна равняться 2)
	if len(packagePart) != 2 {
		return 0, 0, errors.New("Некорректные (неполные) данные")
	}

	// 3) Преобразовываем первый элемент слайса (количество шагов) в int, с обработкой на ошибки
	steps, err := strconv.Atoi(packagePart[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Некорректно указано количество шагов", err)
	}

	// 4) Проверка количества шагов на положительное число
	if steps <= 0 {
		return 0, 0, errors.New("Указано отрицательное количество шагов")
	}

	// 5) Преобразовываем второй элемент слайса (время) в time.Duration с обработкой ошибок
	duration, err := time.ParseDuration(packagePart[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Некорректно указана продолжительность", err)
	}

	// 6) Возврат значений и nil(для ошибки)
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
}
