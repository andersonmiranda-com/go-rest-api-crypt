package main

// Book struct (Model)
type User struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Email       string `json:"email"`
	EmailHash   string `json:"emailHash"`
	Password    string `json:"password"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	Surnames    string `json:"surnames"`
	Mobile      string `json:"mobile"`
	CountryId   string `json:"countryId"`
	Zipcode     string `json:"zipcode"`
	Location    string `json:"location"`
	LanguageId  string `json:"languageId"`
	SexId       string `json:"sexId"`
	Birthdate   string `json:"birthdate"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
	Factor1     string `json:"factor1"`
	Factor2     string `json:"factor2"`
	Status      string `json:"status"`
}

type UserPrivateKey struct {
	N1 string `json:"n1" bson:"n1"`
	R1 string `json:"r1" bson:"r1"`
}

type UserResponse struct {
	Status        string  `json:"status"`
	Data          User    `json:"data"`
	ExecutionTime float64 `json:"executionTime"`
}

type ErrorResponse struct {
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	ExecutionTime float64 `json:"executionTime"`
}

type GenericResponse struct {
	Status        string            `json:"status"`
	Data          map[string]string `json:"data"`
	ExecutionTime float64           `json:"executionTime"`
}
