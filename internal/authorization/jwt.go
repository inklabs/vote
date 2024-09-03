package authorization

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/inklabs/cqrs"
	"go.opentelemetry.io/otel/attribute"
)

const Expiration = 24 * time.Hour

type JWTClaims struct {
	jwt.RegisteredClaims
	Email   string
	UserID  string
	IsAdmin bool
}

type jwtClaimsContext struct {
	claims *JWTClaims
	ctx    context.Context
}

func (j jwtClaimsContext) Context() context.Context {
	return j.ctx
}

func (a jwtClaimsContext) Email() string {
	return a.claims.Email
}

func (a jwtClaimsContext) UserID() string {
	return a.claims.UserID
}

func (a jwtClaimsContext) IsAdmin() bool {
	return a.claims.IsAdmin
}

type jwtAuthorization struct {
	signingKey []byte
}

func NewJWTAuthorization(signingKey []byte) *jwtAuthorization {
	return &jwtAuthorization{
		signingKey: signingKey,
	}
}

func (a *jwtAuthorization) VerifyCommand(ctx context.Context, handler cqrs.CommandHandler, command cqrs.Command) error {
	_, span := tracer.Start(ctx, "jwt-auth.verify-command")
	defer span.End()

	claimsContext, err := a.getContext(ctx)
	if err != nil {
		return err
	}

	if verifier, ok := handler.(CommandVerifier); ok {
		err = verifier.VerifyAuthorization(claimsContext, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *jwtAuthorization) VerifyAsyncCommand(ctx context.Context, handler cqrs.AsyncCommandHandler, command cqrs.AsyncCommand) error {
	_, span := tracer.Start(ctx, "jwt-auth.verify-async-command")
	defer span.End()

	claimsContext, err := a.getContext(ctx)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.String(UserIDKey, claimsContext.UserID()))

	if verifier, ok := handler.(AsyncCommandVerifier); ok {
		err = verifier.VerifyAuthorization(claimsContext, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *jwtAuthorization) VerifyQuery(ctx context.Context, handler cqrs.QueryHandler, query cqrs.Query) error {
	_, span := tracer.Start(ctx, "jwt-auth.verify-query")
	defer span.End()

	claimsContext, err := a.getContext(ctx)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.String(UserIDKey, claimsContext.UserID()))

	if verifier, ok := handler.(QueryVerifier); ok {
		err = verifier.VerifyAuthorization(claimsContext, query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *jwtAuthorization) VerifyRequest(ctx context.Context) error {
	_, span := tracer.Start(ctx, "jwt-auth.verify-request")
	defer span.End()

	claimsContext, err := a.getContext(ctx)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.String(UserIDKey, claimsContext.UserID()))

	return nil
}

func (a *jwtAuthorization) getContext(ctx context.Context) (*jwtClaimsContext, error) {
	if authorizationToken, ok := ctx.Value("authorization").(string); ok {
		splitToken := strings.Split(authorizationToken, "Bearer ")
		if len(splitToken) > 1 {
			claims, err := a.getClaims(splitToken[1])
			if err != nil {
				return nil, err
			}

			return &jwtClaimsContext{
				claims: claims,
				ctx:    ctx,
			}, nil
		}
	}

	return nil, cqrs.ErrAccessDenied
}

func (a *jwtAuthorization) getClaims(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("unable to get claims")
}

func NewSignedToken(claims JWTClaims, signingKey []byte) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(Expiration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "Vote",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
