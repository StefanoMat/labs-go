package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func GetWeather(city string) (WeatherResponse, error) {
	apiKey := "d4dd060277dd4721b5b23031242612"
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("can not return weather")
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return WeatherResponse{}, err
	}

	tempF := result.Current.TempC*1.8 + 32
	tempK := result.Current.TempC + 273.15
	return WeatherResponse{TempC: result.Current.TempC, TempF: tempF, TempK: tempK}, nil
}

func ZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !isValidZipCode(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	locationURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(locationURL)
	if err != nil {
		fmt.Printf("Error getting location: %v\n", err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	var location struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		fmt.Printf("Error decoding location: %v\n", err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	if location.Erro {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weather, err := GetWeather(location.Localidade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func isValidZipCode(zip string) bool {
	re := regexp.MustCompile(`^[0-9]{8}$`)
	return re.MatchString(zip)
}
