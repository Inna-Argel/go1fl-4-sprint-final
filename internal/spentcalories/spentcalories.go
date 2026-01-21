package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверяем, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных, ожидается 'шаги,вид активности,продолжительность'")
	}

	// Парсим количество шагов
	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга количества шагов: %v", err)
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	// Получаем вид активности
	activityType := parts[1]

	// Парсим продолжительность
	durationStr := parts[2]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	// Проверяем, что продолжительность > 0
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient

	// Рассчитываем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	distanceKm := distance(steps, height)

	// Переводим продолжительность в часы
	durationInHours := duration.Hours()

	// Вычисляем среднюю скорость
	return distanceKm / durationInHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные параметры на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// Рассчитываем среднюю скорость
	meanSpeed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитываем количество калорий
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// Рассчитываем среднюю скорость
	meanSpeed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитываем количество калорий с коэффициентом для ходьбы
	calories := (weight * meanSpeed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Парсим данные тренировки
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Вычисляем дистанцию и скорость (общие для всех типов тренировок)
	distanceKm := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)

	// В зависимости от типа активности вычисляем калории
	var calories float64
	var errCal error

	switch strings.ToLower(activityType) {
	case "бег", "running":
		calories, errCal = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба", "walking":
		calories, errCal = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if errCal != nil {
		log.Println(errCal)
		return "", errCal
	}

	// Формируем строку результата
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		distanceKm,
		meanSpeed,
		calories)

	return result, nil
}
