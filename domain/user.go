package user

type User struct {
	name    string `json:"name"`
	age     int64  `json:"age"`
	address string `json:"address"`
}


type CreateUser struct {
	name    string `json:"name"`
	age     int64  `json:"age"`
	address string `json:"address"`
}
