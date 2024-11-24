package token

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/ecdsafile"
	db "talkliketv.click/tltv/db/sqlc"
)

const (
	KeyID            = "fake-key-id"
	Issuer           = "tltv"
	FakeAudience     = "example-users"
	PermissionsClaim = "perm"
	Expiration       = "exp"
	IssuedAt         = "iat"
	UserContextKey   = "user"
	UserIdContextKey = "userid"
)

type FakeAuthenticator struct {
	PrivateKey *ecdsa.PrivateKey
	KeySet     jwk.Set
	// duration of JWS key in hours
	duration *time.Duration
}

var _ JWSValidator = (*FakeAuthenticator)(nil)

// NewFakeAuthenticator creates an authenticator example which uses a hard coded
// ECDSA key to validate JWT's that it has signed itself.
func NewFakeAuthenticator(d *time.Duration, k []byte) (*FakeAuthenticator, error) {
	privKey, err := ecdsafile.LoadEcdsaPrivateKey(k)
	if err != nil {
		return nil, fmt.Errorf("loading PEM private key: %w", err)
	}

	set := jwk.NewSet()
	pubKey := jwk.NewECDSAPublicKey()

	err = pubKey.FromRaw(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("parsing jwk key: %w", err)
	}

	err = pubKey.Set(jwk.AlgorithmKey, jwa.ES256)
	if err != nil {
		return nil, fmt.Errorf("setting key algorithm: %w", err)
	}

	err = pubKey.Set(jwk.KeyIDKey, KeyID)
	if err != nil {
		return nil, fmt.Errorf("setting key ID: %w", err)
	}

	set.Add(pubKey)

	return &FakeAuthenticator{PrivateKey: privKey, KeySet: set, duration: d}, nil
}

// ValidateJWS ensures that the critical JWT claims needed to ensure that we
// trust the JWT are present and with the correct values.
func (f *FakeAuthenticator) ValidateJWS(jwsString string) (jwt.Token, error) {
	return jwt.Parse([]byte(jwsString), jwt.WithKeySet(f.KeySet),
		jwt.WithAudience(FakeAudience), jwt.WithIssuer(Issuer))
}

// SignToken takes a JWT and signs it with our private key, returning a JWS.
func (f *FakeAuthenticator) SignToken(t jwt.Token) ([]byte, error) {
	hdr := jws.NewHeaders()
	if err := hdr.Set(jws.AlgorithmKey, jwa.ES256); err != nil {
		return nil, fmt.Errorf("setting algorithm: %w", err)
	}
	if err := hdr.Set(jws.TypeKey, "JWT"); err != nil {
		return nil, fmt.Errorf("setting type: %w", err)
	}
	if err := hdr.Set(jws.KeyIDKey, KeyID); err != nil {
		return nil, fmt.Errorf("setting Key ID: %w", err)
	}
	return jwt.Sign(t, jwa.ES256, f.PrivateKey, jwt.WithHeaders(hdr))
}

// CreateJWSWithClaims is a helper function to create JWT's with the specified
// claims.
func (f *FakeAuthenticator) CreateJWSWithClaims(claims []string, user db.User) ([]byte, error) {
	t := jwt.New()
	err := t.Set(jwt.IssuerKey, Issuer)
	if err != nil {
		return nil, fmt.Errorf("setting issuer: %w", err)
	}
	err = t.Set(jwt.AudienceKey, FakeAudience)
	if err != nil {
		return nil, fmt.Errorf("setting audience: %w", err)
	}
	err = t.Set(PermissionsClaim, claims)
	if err != nil {
		return nil, fmt.Errorf("setting permissions: %w", err)
	}
	err = t.Set(Expiration, time.Now().Add(time.Hour**f.duration))
	if err != nil {
		return nil, fmt.Errorf("setting expiration: %w", err)
	}
	err = t.Set(IssuedAt, time.Now())
	if err != nil {
		return nil, fmt.Errorf("setting issued at: %w", err)
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("marshalling user: %w", err)
	}
	err = t.Set(jwt.SubjectKey, string(jsonUser))
	if err != nil {
		return nil, fmt.Errorf("setting subject key: %w", err)
	}
	return f.SignToken(t)
}
