package ebayapi
import (
	"sync"
	"time"
)

// thread safe, immutable structure that abstracts
// away refreshing tokens, auto refreshes if trying to
// use token after expiry
type TokenManager struct {
	b64secret string
	token     string
	expiresAt time.Time // time that token expires
	mu        sync.Mutex
}

// note time is expiresIn (duration till expiration) not when it expires
func NewTokenManager(b64secret string, token string, expiresIn time.Duration) *TokenManager {
	return &TokenManager{
		b64secret: b64secret,
		token: token,
		expiresAt: time.Now().Add(expiresIn),
		mu: sync.Mutex{},
	}
}

// check if token is expired if so refresh token
// then return new value otherwise just return token
func (T *TokenManager) getToken() string {
	T.mu.Lock()
	defer T.mu.Unlock()

	if time.Now().After(T.expiresAt) {
		// refresh token
		response, err := RefreshToken(T.b64secret)
		
		if err != nil {
			panic("could not refresh token")
		}

		// update times
		T.token = response.AccessToken
		T.expiresAt = time.Now().Add(time.Duration(response.ExpiresIn))
	}

	return T.token
}


