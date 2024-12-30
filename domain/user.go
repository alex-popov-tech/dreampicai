package domain

const (
	AccessTokenCookieKey  = "at"
	RefreshTokenCookieKey = "rt"
	AccountIdCookieKey    = "a"
)

const (
	AccountContextKey = "account"
)

type UserAuth struct {
	ID           string
	Email        string
	Provider     string
	AccessToken  string
	RefreshToken string
}

type Account struct {
	ID    int32
	Email string
	// IsLoggedIn bool
	UserAuth UserAuth
}
