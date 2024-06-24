package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
)

type requestsPayload struct {
	Age                                int     `json:"age"`
	Gender                             int     `json:"gender"`
	ChestPain                          int     `json:"chest_pain"`
	RestingBloodPressure               int     `json:"resting_blood_pressure"`
	CholestoralInMg                    int     `json:"cholestoral_in_mg"`
	FastingBloodSugar                  int     `json:"fasting_blood_sugar"`
	RestingElectrocardiographicResults int     `json:"resting_electrocardiographic_results"`
	MaximumHeartRateAchieved           int     `json:"maximum_heart_rate_achieved"`
	ExerciseInducedAngina              int     `json:"exercise_induced_angina"`
	PreviousPeak                       float64 `json:"previous_peak"`
	SlopeOfThePeakExercise             int     `json:"slope_of_the_peak_exercise"`
	NumberOfMajorVessels               int     `json:"number_of_major_vessels"`
	Thalassemia                        int     `json:"thalassemia"`
}

func (app *Config) KNN(write http.ResponseWriter, read *http.Request) {
	var requests_payload requestsPayload

	possible_error := app.readJSON(write, read, &requests_payload)

	if possible_error != nil {
		app.errorJSON(write, possible_error, http.StatusBadRequest)
		return
	}

	X, y := load_dataset("heart.csv")

	X_to_predict := []int{
		requests_payload.Age,
		requests_payload.Gender,
		requests_payload.ChestPain,
		requests_payload.RestingBloodPressure,
		requests_payload.CholestoralInMg,
		requests_payload.FastingBloodSugar,
		requests_payload.RestingElectrocardiographicResults,
		requests_payload.MaximumHeartRateAchieved,
		requests_payload.ExerciseInducedAngina,
		int(requests_payload.PreviousPeak),
		requests_payload.SlopeOfThePeakExercise,
		requests_payload.NumberOfMajorVessels,
		requests_payload.Thalassemia,
	}

	X_scaled := minmax_scale_fit_transform(X)
	X_scaled_to_predict := minmax_to_predict_scale_fit_transform(X_to_predict)

	y_predicted := predict(X_scaled_to_predict, X_scaled, y, 3)

	var result string

	if y_predicted[0] == 0 {
		result = "Yes"
	} else {
		result = "No"
	}

	pay_load := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("The result is: %s", result),
	}

	app.writeJSON(write, http.StatusAccepted, pay_load)
}

func load_dataset(file_name string) ([][]int, []int) {
	var X [][]int
	var y []int

	// Try to open the csv file in read-write mode.
	csvFile, csvFileError := os.OpenFile(file_name, os.O_RDWR, os.ModePerm)
	if csvFileError != nil {
		panic(csvFileError)
	}
	// Ensure the file is closed once the function returns
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	_, err := reader.Read() // Skips header
	if err != nil {
		return X, y
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		inR0, _ := strconv.Atoi(row[0])
		inR1, _ := strconv.Atoi(row[1])
		inR2, _ := strconv.Atoi(row[2])
		inR3, _ := strconv.Atoi(row[3])
		inR4, _ := strconv.Atoi(row[4])
		inR5, _ := strconv.Atoi(row[5])
		inR6, _ := strconv.Atoi(row[6])
		inR7, _ := strconv.Atoi(row[7])
		inR8, _ := strconv.Atoi(row[8])
		inR9, _ := strconv.Atoi(row[9])
		inR10, _ := strconv.Atoi(row[10])
		inR11, _ := strconv.Atoi(row[7])
		inR12, _ := strconv.Atoi(row[8])
		inR13, _ := strconv.Atoi(row[9])

		row_to_add := []int{inR0, inR1, inR2, inR3, inR4, inR5, inR6, inR7, inR8, inR9, inR10, inR11, inR12}

		X = append(X, row_to_add)
		y = append(y, inR13)
	}

	return X, y
}

func minmax_to_predict_scale_fit_transform(X_to_predict []int) []float32 {

	X_scaled := make([]float32, len(X_to_predict))

	min_value := slices.Min(X_to_predict)
	max_value := slices.Max(X_to_predict)
	range_value := max_value - min_value

	for index, value := range X_to_predict {
		new_value := float32(value-min_value) / float32(range_value)
		X_scaled[index] = new_value
	}

	return X_scaled
}

func minmax_scale_fit_transform(X [][]int) [][]float32 {

	X_scaled := make([][]float32, len(X))

	for index, element := range X {
		min_value := slices.Min(element)
		max_value := slices.Max(element)
		range_value := max_value - min_value

		for _, value := range element {
			new_value := float32(value-min_value) / float32(range_value)
			X_scaled[index] = append(X_scaled[index], new_value)
		}
	}

	return X_scaled
}

func calc_distance(X_to_predict []float32, X [][]float32) []float32 {

	distances := make([]float32, len(X))

	for i := 0; i < len(X); i++ {
		euclidean_distance := 0.0

		x := X[i]
		y := X_to_predict

		for r := 0; r < len(x); r++ {
			difference := x[r] - y[r]
			euclidean_distance += math.Pow(float64(difference), 2)
		}

		euclidean_distance = math.Sqrt(float64(euclidean_distance))
		distances[i] = float32(euclidean_distance)
	}

	return distances
}

func predict(X_to_predict []float32, X [][]float32, y []int, k int) []int {
	distances_array := calc_distance(X_to_predict, X)

	y_predicted := make([]int, len(y))

	group_A := 0
	group_B := 0

	for i := 0; i < k; i++ {
		closest_index := slices.Index(distances_array, slices.Min(distances_array))

		if y[closest_index] == 1 {
			group_A += 1
		} else {
			group_B += 1
		}

		distances_array = append(distances_array[:closest_index], distances_array[closest_index+1:]...)
	}

	if group_A > group_B {
		y_predicted = append(y_predicted, 1)
	} else {
		y_predicted = append(y_predicted, 0)
	}

	return y_predicted
}
