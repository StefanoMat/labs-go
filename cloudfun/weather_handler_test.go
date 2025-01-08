package main_test

import (
	"encoding/json"
	"github.com/stefanomat/weather/weather"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func TestWeatherHandler(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", weather.ZipCodeHandler)

	t.Run("GET request success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/weather?cep=95650000", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Status code deve ser 200")

		var resp WeatherResponse
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err, "Deve conseguir deserializar JSON de retorno")
		assert.NotEmptyf(t, resp.TempC, "TempC não deve ser vazio")
		assert.NotEmptyf(t, resp.TempF, "TempF não deve ser vazio")
		assert.NotEmptyf(t, resp.TempK, "TempK não deve ser vazio")
	})

	t.Run("GET request invalid cep", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/weather?cep=012345670", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "Status code deve ser 422")
		assert.Contains(t, rr.Body.String(), "invalid zipcode", "Mensagem deve indicar CEP inválido")
	})

	t.Run("GET request can not find zipcode", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/weather?cep=99999999", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code, "Status code deve ser 404 (Not Found)")
		assert.Contains(t, rr.Body.String(), "can not find zipcode", "Mensagem deve indicar CEP não encontrado")
	})
}
