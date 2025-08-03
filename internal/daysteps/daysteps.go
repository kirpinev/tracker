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
	values := strings.Split(data, ",")

	if len(values) != 2 {
		return 0, 0, errors.New("длина меньше двух")
	}

	steps, err := strconv.Atoi(values[0])

	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов меньше или равно 0")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(values[1]))

	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)

		return ""
	}

	distance := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Println(err)

		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
