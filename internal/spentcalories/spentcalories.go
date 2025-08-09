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

func validateInputs(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return errors.New("number of steps is less than or equal to 0")
	}

	if weight <= 0 {
		return errors.New("weight must be greater than zero")
	}

	if height <= 0 {
		return errors.New("height must be greater than zero")
	}

	if duration <= 0 {
		return errors.New("duration must be greater than zero")
	}

	return nil
}

func formatTrainingInfo(activity string, steps int, height float64, duration time.Duration, calories float64) string {
	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), distance, meanSpeed, calories)
}

func parseTraining(data string) (int, string, time.Duration, error) {
	values := strings.Split(data, ",")

	if len(values) != 3 {
		return 0, "", 0, errors.New("length is less than three")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(values[0]))

	if err != nil {
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(values[2])

	if err != nil {
		return 0, "", 0, err
	}

	activity := strings.TrimSpace(values[1])

	if err := validateInputs(steps, 1, 1, duration); err != nil {
		return 0, "", 0, err
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepsLength := stepLengthCoefficient * height

	return float64(steps) * stepsLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return formatTrainingInfo(activity, steps, height, duration, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInputs(steps, weight, height, duration); err != nil {
		return 0, err
	}

	meanSpeed := meanSpeed(steps, height, duration)

	return (weight * meanSpeed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInputs(steps, weight, height, duration); err != nil {
		return 0, err
	}

	meanSpeed := meanSpeed(steps, height, duration)
	calories := (weight * meanSpeed * duration.Minutes()) / minInH

	return calories * walkingCaloriesCoefficient, nil
}
