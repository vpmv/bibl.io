package dto

type Author struct {
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	DateOfBirth string  `json:"date_of_birth" example:"1952-01-18"`
	DateOfDeath string  `json:"date_of_death" example:"2020-07-25"`
	Nationality string  `json:"nationality"`
	Description string  `json:"description"`
	Books       []*Book `json:"books,omitempty"`
}

type Book struct {
	Key      string    `json:"key"`
	ISBN     []string  `json:"isbn"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Language string    `json:"language"`
	Authors  []*Author `json:"authors,omitempty"`
}
