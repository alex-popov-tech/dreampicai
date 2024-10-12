package model

const UserContextKey = "user"

type User struct {
	ID         string
	Email      string
	IsLoggedIn bool
}
