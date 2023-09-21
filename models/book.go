package models

type Book struct {
	ID     string `json:"id" bson:"_id"`
	Author string `json:"author" bson:"author"`
	Title  string `json:"title" bson:"title"`
	Year   string `json:"year" bson:"year"`
}
