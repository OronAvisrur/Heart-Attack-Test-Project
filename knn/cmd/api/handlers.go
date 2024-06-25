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

// This function execute KNN algorithm on given data to predict the json message result
func (app *Config) KNN(write http.ResponseWriter, read *http.Request) {
	var requests_payload requestsPayload

	// Write the json to requestsPayload struct
	possible_error := app.readJSON(write, read, &requests_payload)

	if possible_error != nil {
		app.errorJSON(write, possible_error, http.StatusBadRequest)
		return
	}

	// Load the csv into slice [][]int object and seperate X, y by the last column
	X, y := load_dataset("heart.csv")

	// Set the payload as []int slice
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

	// Do scalling for X and X_to_predict
	X_scaled := minmax_scale_fit_transform(X)
	X_scaled_to_predict := minmax_to_predict_scale_fit_transform(X_to_predict)

	// Try to predict the result
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

	// Return answer to the broker
	app.writeJSON(write, http.StatusAccepted, pay_load)
}

// This function get csv file and return it as two slices of type int
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
	_, possible_error := reader.Read() // Skips header
	if possible_error != nil {
		return X, y
	}

	for {
		row, possible_error := reader.Read()
		if possible_error == io.EOF {
			break
		}
		if possible_error != nil {
			log.Fatal(possible_error)
		}

		//convert the string to int and ignore the error
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

		//Create a vector to add to X
		row_to_add := []int{inR0, inR1, inR2, inR3, inR4, inR5, inR6, inR7, inR8, inR9, inR10, inR11, inR12}

		//Add the vector to X and the last column value to the y
		X = append(X, row_to_add)
		y = append(y, inR13)
	}

	return X, y
}

// This function will do scalling to the X_to_predict in order to ensure the
// varibles will be between 0 to 1 so the prediction will be more acurate since all the varible on the same scale
func minmax_to_predict_scale_fit_transform(X_to_predict []int) []float32 {

	// Create the return slice
	X_scaled := make([]float32, len(X_to_predict))

	// Find min and max values in the vector
	min_value := slices.Min(X_to_predict)
	max_value := slices.Max(X_to_predict)
	range_value := max_value - min_value

	// Set new value to each varible on the vector
	for index, value := range X_to_predict {
		new_value := float32(value-min_value) / float32(range_value)
		X_scaled[index] = new_value
	}

	return X_scaled
}

// This function will do scalling to the X in order to ensure the
// varibles will be between 0 to 1 so the prediction will be more acurate since all the varible on the same scale
func minmax_scale_fit_transform(X [][]int) [][]float32 {

	// Create the return slice
	X_scaled := make([][]float32, len(X))

	// Loop through all the vectors
	for index, element := range X {

		// Find min and max values in the vector
		min_value := slices.Min(element)
		max_value := slices.Max(element)
		range_value := max_value - min_value

		// Set new value to each varible on the vector
		for _, value := range element {
			new_value := float32(value-min_value) / float32(range_value)
			X_scaled[index] = append(X_scaled[index], new_value)
		}
	}

	return X_scaled
}

// This function calculate distance between X_to_predict vector and all X vectors
// and return the distances
func calc_distance(X_to_predict []float32, X [][]float32) []float32 {

	distances := make([]float32, len(X))

	// Loop through all the vectors and calculate the distance with euclidean distance formula
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

// This function will predict the result of X_to_predict based on the results of X
func predict(X_to_predict []float32, X [][]float32, y []int, k int) []int {
	// Calculate distance between X_to_predict vector and all X vectors
	distances_array := calc_distance(X_to_predict, X)

	y_predicted := make([]int, len(y))

	group_A := 0
	group_B := 0

	// Run K times and find the K neighbors of X_to_predict vector
	for i := 0; i < k; i++ {
		closest_index := slices.Index(distances_array, slices.Min(distances_array))

		if y[closest_index] == 1 {
			group_A += 1
		} else {
			group_B += 1
		}

		distances_array = append(distances_array[:closest_index], distances_array[closest_index+1:]...)
	}

	// If more neighbors are from group_A X_to_predict is also belong to group_A
	// otherwise X_to_predict belong to group_B
	if group_A > group_B {
		y_predicted = append(y_predicted, 1)
	} else {
		y_predicted = append(y_predicted, 0)
	}

	return y_predicted
}
