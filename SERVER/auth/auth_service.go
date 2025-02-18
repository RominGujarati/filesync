package auth

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("9cd4b6bb1975b118e64e80d4804bf9283d828c35a1344f1fe25910ee0c8f37c73bff118e985cae63d0ef77e43bfd7461c1e908b081bd8bd178218e4416d51784ae1d60269167c89146ea0b99ef0c6e7594514cd33dfac027d652a3e6856ffba66e76849c96b0e895bc9149c73d24b58124308cfcc8d64eac07ae35fe3aba10062c7e826c338edb9bc6880fd0d7b73214ea31197bbf1e300b83733cd67d13dfaf6de8fbaa505d07692e3b126b690782ca963ef9406189ff029e0670071e0cc92701881f4dd8509f4863e0950e819b5eca627448df6ce1cbc3cbe2866b0e617b9bba89301b91f0d7d699fbf2c1c9c1953230f9f2ef68be020e560e781a2731095e")

// Hash the user password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check password validity
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT token
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
