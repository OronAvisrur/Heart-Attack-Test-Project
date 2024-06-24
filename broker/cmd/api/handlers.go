package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
	Knn    KnnPayload  `json:"knn,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type KnnPayload struct {
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

// Broker handler for the Config type
func (app *Config) Broker(write http.ResponseWriter, read *http.Request) {
	// Define the payload to be sent as a JSON response
	payload := jsonResponse{
		Error:   false,
		Message: "Clicked the broker",
	}

	// Write the JSON response with status OK (200)
	_ = app.writeJSON(write, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(write http.ResponseWriter, read *http.Request) {
	var request_payload RequestPayload

	possible_error := app.readJSON(write, read, &request_payload)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	switch request_payload.Action {
	case "auth":
		app.authenticate(write, request_payload.Auth)
	case "mail":
		app.sendMail(write, request_payload.Mail)
	case "knn":
		app.calculateKNN(write, request_payload.Knn)
	default:
		app.errorJSON(write, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(write http.ResponseWriter, authentic AuthPayload) {
	//create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(authentic, "", "\t")

	//call the service
	request, possible_error := http.NewRequest("POST", "http://authentication/authenticate", bytes.NewBuffer(jsonData))
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	client := &http.Client{}
	response, possible_error := client.Do(request)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(write, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(write, errors.New("error calling auth service"))
		return
	}

	//create a varible we'll read response.Body into
	var jsonFromService jsonResponse

	//decode the json from the auth service
	possible_error = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(write, possible_error, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(write, http.StatusAccepted, payload)
}

func (app *Config) sendMail(write http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	//call the mail service
	mailServiceURL := "http://mailer/send"

	//post to mail service
	request, possible_error := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, possible_error := client.Do(request)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}
	defer response.Body.Close()

	//make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(write, errors.New("error calling mail service"))
		return
	}

	//send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(write, http.StatusAccepted, payload)
}

func (app *Config) calculateKNN(write http.ResponseWriter, authentic KnnPayload) {
	//create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(authentic, "", "\t")

	//call the service
	request, possible_error := http.NewRequest("POST", "http://knn/knn", bytes.NewBuffer(jsonData))
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	client := &http.Client{}
	response, possible_error := client.Do(request)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}
	defer response.Body.Close()

	//make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(write, errors.New("error calling knn service"))
		return
	}

	//create a varible we'll read response.Body into
	var jsonFromService jsonResponse

	//decode the json from the auth service
	possible_error = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(write, possible_error, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(write, http.StatusAccepted, payload)
}
