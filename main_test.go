package main

import (
	"testing"
)

func TestKelvinToCelsius(t *testing.T) {
	data := Data{
		Main: Main{Temp: 300, Temp_min: 290, Temp_max: 310},
	}
	data.Main.kelvinToCelsius()
	data.TruncateDecimals()

	if data.Main.Temp != 27 {
		t.Errorf("Expected 26, but got %f", data.Main.Temp)
	}
	if data.Main.Temp_min != 16 {
		t.Errorf("Expected 16, but got %f", data.Main.Temp_min)
	}
	if data.Main.Temp_max != 36 {
		t.Errorf("Expected 36, but got %f", data.Main.Temp_max)
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	data := Data{
		Main: Main{Temp: 25, Temp_min: 20, Temp_max: 30},
	}
	data.Main.celsiusToFarhenheit()
	data.TruncateDecimals()

	if data.Main.Temp != 77 {
		t.Errorf("Expected 77, but got %f", data.Main.Temp)
	}
	if data.Main.Temp_min != 68 {
		t.Errorf("Expected 68, but got %f", data.Main.Temp_min)
	}
	if data.Main.Temp_max != 86 {
		t.Errorf("Expected 86, but got %f", data.Main.Temp_max)
	}
}

func TestMpsToMph(t *testing.T) {
	data := Data{
		Wind: Wind{Speed: 10},
	}
	data.Wind.mpsToMph()
	data.TruncateDecimals()

	if data.Wind.Speed != 22 {
		t.Errorf("Expected %v, but got %f", 22, data.Wind.Speed)
	}
}
