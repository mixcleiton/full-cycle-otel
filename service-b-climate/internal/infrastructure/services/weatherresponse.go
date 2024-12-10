package weatherapi

type WeatherCurrentResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Current struct {
	LastUpdatedEpoch int64   `json:"last_updated_epoch"`
	LastUpdated      string  `json:"last_updated"`
	TempC            float64 `json:"temp_c"`
	TempF            float64 `json:"temp_f"`
	IsDay            int     `json:"is_day"`
	Condition        `json:"condition"`
	WindMPH          float64 `json:"wind_mph"`
	WindKph          float64 `json:"wind_kph"`
	WindDegree       int     `json:"wind_degree"`
	WindDir          string  `json:"wind_dir"`
	PressureMB       float64 `json:"pressure_mb"`
	PressureIn       float64 `json:"pressure_in"`
	PrecipMM         float64 `json:"precip_mm"`
	PrecipIn         float64 `json:"precip_in"`
	Humidity         int     `json:"humidity"`
	Cloud            int     `json:"cloud"`
	FeelslikeC       float64 `json:"feelslike_c"`
	FeelslikeF       float64 `json:"feelslike_f"`

	WindchillC float64 `json:"windchill_c"`
	WindchillF float64 `json:"windchill_f"`
	HeatindexC float64 `json:"heatindex_c"`
	HeatindexF float64 `json:"heatindex_f"`
	DewpointC  float64 `json:"dewpoint_c"`
	DewpointF  float64 `json:"dewpoint_f"`

	VisKM    float64 `json:"vis_km"`
	VisMiles float64 `json:"vis_miles"`
	Uv       float64 `json:"uv"`
	GustMPH  float64 `json:"gust_mph"`
	GustKph  float64 `json:"gust_kph"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocalTimeEpoch int64   `json:"localtime_epoch"`
	LocalTime      string  `json:"localtime"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}
