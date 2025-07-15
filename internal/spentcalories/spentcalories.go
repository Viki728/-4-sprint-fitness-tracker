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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	//Разделяем строку на слайс строк
	sliceStepActionTime := strings.Split(data, ",")

	//Проверяем, чтобы длина слайса была равна 3, иначе  - выдаем ошибку
	if len(sliceStepActionTime) < 3 || len(sliceStepActionTime) > 3 {
		return 0, "", 0, errors.New("slice length is less than 3")
	}

	//Преобразуем первый и третий элементы слайса в int и time.Duration соответсвенно
	//и возвращаем возможные ошибки
	stepParseTraining, err := strconv.Atoi(sliceStepActionTime[0])
	if err != nil {
		return 0, "", 0, errors.New("failed to convert first element of slice")
	}

	if stepParseTraining <= 0 {
		return 0, "", 0, errors.New("number of steps <= 0")
	}

	timeParseTraining, err := time.ParseDuration(sliceStepActionTime[2])
	if err != nil {
		return 0, "", 0, errors.New("failed to convert three element of slice")
	}

	if timeParseTraining <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}
	return stepParseTraining, sliceStepActionTime[1], timeParseTraining, nil
}

func distance(steps int, height float64) float64 {
	//Определяем пройденную дистанцию в км
	valueDistance := (height * stepLengthCoefficient * float64(steps)) / mInKm
	return valueDistance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	//Проверка продолжительности активности
	if duration <= 0 {
		return 0
	}

	//Вычисляем дистанцию при помощи функции distance()
	distMeanSpeed := distance(steps, height)

	//Вычисляем и возвращаем среднюю скорость
	vMeanSpeed := distMeanSpeed / duration.Hours()
	return vMeanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	//Получить значения из строки данных с помощью функции parseTraining()
	stepsTrainingInfo, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if duration <= 0 {
		log.Println("duration of the workout is 0")
		return "", err
	}
	//Проверяем какой вид тренировки был передан в строке
	switch activity {
	case "Бег":
		distanceTrainingInfo := distance(stepsTrainingInfo, height)
		v := meanSpeed(stepsTrainingInfo, height, duration)
		calor, err := RunningSpentCalories(stepsTrainingInfo, weight, height, duration)
		if err != nil {
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distanceTrainingInfo, v, calor)
		return result, nil

	case "Ходьба":
		distanceTrainingInfo := distance(stepsTrainingInfo, height)
		v := meanSpeed(stepsTrainingInfo, height, duration)
		calor, err := WalkingSpentCalories(stepsTrainingInfo, weight, height, duration)
		if err != nil {
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distanceTrainingInfo, v, calor)
		return result, nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//Проверка входных параметров на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("incorrect input parameters")
	}

	//Рассчитываем среднюю скорость с помощью meanSpeed()
	vRunningSpentCalories := meanSpeed(steps, height, duration)

	//Рассчитываем и возвращаем кол-во калорий
	calorRunningSpentCalories := (duration.Minutes() * weight * vRunningSpentCalories) / minInH
	return calorRunningSpentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//Проверка входных параметров на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("incorrect input parameters")
	}

	//Рассчитываем среднюю скорость с помощью meanSpeed()
	vWalkingSpentCalories := meanSpeed(steps, height, duration)

	//Рассчитываем и возвращаем кол-во калорий (с учетом корректирующего коэффициента)
	calorWalkingSpentCalories := ((duration.Minutes() * weight * vWalkingSpentCalories) / minInH) * walkingCaloriesCoefficient
	return calorWalkingSpentCalories, nil
}
