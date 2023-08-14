package models

import "time"

type Movie struct{
	ID 			int 		`json:"id"`
	Title 		string 		`json:"title"`
	ReleaseDate time.Time 	`json:"release_date"`
	RunTime 	int 		`json:"runtime"`
	MPAARating 	string 		`json:"mpaa_rating"`
	Description string 		`json:"description"`
	Image 		string 		`json:"image"`
	CreatedAT 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}