package domain

const AccountContextKey = "account"
const UserContextKey = "user"
const AuthContextKey = "auth"

type Auth struct {
	ID       string
	Email    string
	Provider string
	IsInit   bool
}

type User struct {
	ID         string
	Email      string
	IsLoggedIn bool
}
