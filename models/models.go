package models

type Service struct {
	Id          string
	Name        string
	Description string
	Versions    map[string]string
}

type ServiceVersion struct {
	Name string
}
