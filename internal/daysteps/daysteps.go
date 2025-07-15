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

	//Разделяем строку на слайс строк
	sliceStepAndTime := strings.Split(data, ",")

	//Проверяем, чтобы длина слайса была равна 2. В ином случае выводим ошибку
	if len(sliceStepAndTime) != 2 {
		return 0, 0, errors.New("slice length is less than 2")
	}

	//Преобразуем кол-во шагов в тип int
	step, err := strconv.Atoi(sliceStepAndTime[0])
	if err != nil {
		return 0, 0, errors.New("failed to convert first element of slice")
	}

	//Возвращаем ошибку, если количество шагов равно 0
	if step <= 0 {
		return 0, 0, errors.New("number of steps <= 0")
	}

	//Преобразуем второй элемент слайса в time.Duration
	t, err := time.ParseDuration(sliceStepAndTime[1])
	if err != nil {
		return 0, 0, errors.New("failed to convert second element of slice")
	}

	if t <= 0 {
		return 0, 0, errors.New("duration <= 0")
	}
	return step, t, nil
}

func DayActionInfo(data string, weight, height float64) string {
	//Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	//Проверяем, чтобы количество шагов было больше 0
	if steps == 0 {
		return ""
	}

	//Вычисляем дистанцию в км
	diststion := (float64(steps) * stepLength) / float64(mInKm)

	//Вычисляем кол-во калорий, потраченных на прогулке

	calor, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	//Вывод результатов
	result := fmt.Sprintf("Количество шагов: %d.\n", steps) + fmt.Sprintf("Дистанция составила %.2f км.\n", diststion) + fmt.Sprintf("Вы сожгли %.2f ккал.\n", calor)
	return result
}
