package airline

// Airline represents the structure of an airline.
type Airline struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl"`
}
