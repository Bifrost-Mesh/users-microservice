package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	goJWT "github.com/golang-jwt/jwt/v5"

	"github.com/Bifrost-Mesh/users-microservice/pkg/constants"
	"github.com/Bifrost-Mesh/users-microservice/pkg/utils"
)

const JWT_VALIDITIY_PERIOD = 24 * time.Hour

type (
	JWTService struct {
		signingKey string
	}

	// JSON web tokens (JWTs) claims are pieces of information asserted about a subject.
	JWTClaims struct {
		// Registered claims are standard claims registered with the Internet Assigned Numbers
		// Authority (IANA) and defined by the JWT specification to ensure interoperability with
		// third-party, or external, applications.
		goJWT.RegisteredClaims

		/*
		  Custom claims consist of non-registered public or private claims.

		    (1) Public claims : You can create custom claims for public consumption, which might
		        contain generic information like name and email. If you create public claims, you must
		        either register them or use collision-resistant names through namespacing and take
		        reasonable precautions to make sure you are in control of the namespace you use.

		    (2) Private claims : You can create private custom claims to share information specific to
		        your application. For example, while a public claim might contain generic information
		        like name and email, private claims would be more specific, such as employee ID and
		        department name.
		*/
	}
)

func NewJWTService(signingKey string) *JWTService {
	return &JWTService{signingKey}
}

func (j *JWTService) Issue(userID int32) (*string, error) {
	jwtSigner := goJWT.NewWithClaims(goJWT.SigningMethodHS256, JWTClaims{
		//nolint:exhaustruct
		RegisteredClaims: goJWT.RegisteredClaims{
			Issuer:    "Instagram Clone",
			Subject:   strconv.Itoa(int(userID)),
			IssuedAt:  goJWT.NewNumericDate(time.Now()),
			ExpiresAt: goJWT.NewNumericDate(time.Now().Add(JWT_VALIDITIY_PERIOD)),
		},
	})

	jwt, err := jwtSigner.SignedString([]byte(j.signingKey))
	if err != nil {
		return nil, utils.WrapErrorWithPrefix("Failed generating JWT", err)
	}
	return &jwt, nil
}

func (j *JWTService) GetUserIDFromToken(jwt string) (*int32, error) {
	parsedJWT, err := goJWT.Parse(jwt,
		func(_ *goJWT.Token) (any, error) { return []byte(j.signingKey), nil },
		goJWT.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, goJWT.ErrTokenExpired) {
			return nil, constants.ErrExpiredJWT
		}

		return nil, constants.ErrInvalidJWT
	}

	subject, err := parsedJWT.Claims.GetSubject()
	if err != nil {
		return nil, constants.ErrInvalidJWT
	}

	userID, err := strconv.ParseInt(subject, 10, 32)
	if err != nil {
		return nil, constants.ErrInvalidJWT
	}
	return aws.Int32(int32(userID)), nil
}
