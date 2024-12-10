package climateapi

type CurrentClimateResponse struct {
	City       string  `json:"city"`
	Fahrenheit float64 `json:"temp_F"`
	Celsius    float64 `json:"temp_C"`
	Kelvin     float64 `json:"temp_K"`
}
