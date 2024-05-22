package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type Data struct {
	Main Main `json:"main"`
	Wind Wind `json:"wind"`
	AQI  bigquery.Value
}

type Main struct {
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	Temp_min float64 `json:"temp_min"`
	Temp_max float64 `json:"temp_max"`
	Pressure float64 `json:"pressure"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

func (m *Main) kelvinToCelsius() {
	m.Temp = m.Temp - 273.15
	m.Temp_min = m.Temp_min - 273.15
	m.Temp_max = m.Temp_max - 273.15
}

func (m *Main) celsiusToFarhenheit() {
	m.Temp = (m.Temp * 9 / 5) + 32
	m.Temp_max = (m.Temp_max * 9 / 5) + 32
	m.Temp_min = (m.Temp_min * 9 / 5) + 32
}

func (w *Wind) mpsToMph() {
	w.Speed = w.Speed * 2.23694
}

func (d *Data) TruncateDecimals() {
	d.Main.Temp = math.Trunc(d.Main.Temp)
	d.Main.Temp_min = math.Trunc(d.Main.Temp_min)
	d.Main.Temp_max = math.Trunc(d.Main.Temp_max)
	d.Wind.Speed = math.Trunc(d.Wind.Speed)
}

func fetchLiveWeatherData() Data {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=40.65&lon=-111.85&appid=25a3874cd58d5c0d253a5a7fc33f9ebe"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	data.Main.kelvinToCelsius()

	return data
}

func (d *Data) predict() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := bigquery.NewClient(ctx, "continual-modem-424017-r2")
	if err != nil {
		log.Printf("Failed to create BigQuery client: %v", err)
	}
	defer client.Close()

	predictionQuery := fmt.Sprintf(
		`
		SELECT *
		FROM ML.PREDICT(MODEL %s,
		  (
		  SELECT
		    TIMESTAMP("%v") AS date_timestamp,
		    %v AS temperature,
		    %v AS temp_min,
		    %v AS temp_max,
		    %v AS wind,
		    %v AS humidity,
			%v AS pressure
		  )
		)
		`,
		"`continual-modem-424017-r2.data.prediction_model`",
		time.Now().Format("2006-01-02"),
		d.Main.Temp,
		d.Main.Temp_min,
		d.Main.Temp_max,
		d.Wind.Speed,
		d.Main.Humidity,
		d.Main.Pressure)

	query := client.Query(predictionQuery)
	it, err := query.Read(ctx)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
	}

	var results []map[string]bigquery.Value
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to read query results: %v", err)
		}
		results = append(results, row)
	}

	d.AQI = results[0]["predicted_AQI"]
}
