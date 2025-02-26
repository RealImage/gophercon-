package degrees

// common config for movie or person
type commonConfig struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

// config for getting actor data
type RespActorDto struct {
	Url    string         `json:"url"`
	Type   string         `json:"type"`
	Name   string         `json:"name"`
	Movies []commonConfig `json:"movies"`
}

// config for getting movie data
type RespMovieDto struct {
	Url  string         `json:"url"`
	Type string         `json:"type"`
	Name string         `json:"name"`
	Cast []commonConfig `json:"cast"`
	Crew []commonConfig `json:"crew"`
}

// for storing path taken
type Path struct {
	Nodes []Node
}

// storing info for path which will be required later
type Node struct {
	Person1 string
	Role1   string
	Person2 string
	Role2   string
	Url2    string
	Movie   string
}
