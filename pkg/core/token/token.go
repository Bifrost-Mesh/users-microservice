package token

type (
	Token = string

	TokenService interface {
		Issue(userID int32) (*Token, error)
		GetUserIDFromToken(token Token) (*int32, error)
	}
)
