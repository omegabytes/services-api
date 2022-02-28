package models

type Service struct {
	Id          int
	Name        string
	Description string
	Versions    []ServiceVersion
}

type ServiceVersion struct {
	SemVer string `json:"semver"`
}
