package models

// Cast represents actors, actresses, and supporting roles in a movie
type Cast struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Crew represents directors, producers, and other film crew members
type Crew struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Movie represents a movie and its cast/crew
type Movie struct {
	URL   string `json:"url"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Cast  []Cast `json:"cast"`
	Crew  []Crew `json:"crew"`
}

// MovieRef represents a movie reference in a person's filmography
type MovieRef struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Person represents an actor or crew member
type Person struct {
	URL    string     `json:"url"`
	Type   string     `json:"type"`
	Name   string     `json:"name"`
	Movies []MovieRef `json:"movies"`
}