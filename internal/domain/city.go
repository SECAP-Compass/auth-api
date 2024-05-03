package domain

type City struct {
	Name        string `json:"name"`
	Population  int    `json:"population"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	Maps struct {
		GoogleMaps    string `json:"googleMaps"`
		OpenStreetMap string `json:"openStreetMap"`
	} `json:"maps"`
}
