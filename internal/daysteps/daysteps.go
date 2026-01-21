package daysteps

import (
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
	// Разделяем строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных, ожидается 'шаги, продолжительность'")
	}

	// Парсим количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга количества шагов: %v", err)
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем количество потраченных калорий
	calories, err := WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Формируем строку
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories)
}
