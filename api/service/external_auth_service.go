package service

import (
	"api/model"
	"api/repository"
	"api/utils/constants"
	"api/utils/errs"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"math/big"
	"net/http"
)

type IExternalAuthService interface {
	GoogleAuth(token string) (*model.User, error)
	AppleAuth(token string) (*model.User, error)
}

type ExternalAuthService struct {
	userRepo repository.IUserRepo
	codeRepo repository.ICodeRepo
}

func NewExternalAuthService(userRepo repository.IUserRepo, codeRepo repository.ICodeRepo) *ExternalAuthService {
	return &ExternalAuthService{
		userRepo: userRepo,
		codeRepo: codeRepo,
	}
}

type AppleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}
type AppleKeyResponse struct {
	Keys []AppleKey `json:"keys"`
}
type AppleIDTokenPayload struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func (s *ExternalAuthService) GoogleAuth(token string) (*model.User, error) {
	tokenInfo, err := validateGoogleToken(token)
	if err != nil {
		return nil, err
	}

	return s.findOrCreateGoogleUser(tokenInfo)
}

func (s *ExternalAuthService) AppleAuth(token string) (*model.User, error) {
	clientID := "sundayti.olympguide"
	tokenPayload, err := validateAppleIDToken(clientID, token)
	if err != nil {
		return nil, err
	}
	return s.findOrCreateAppleUser(tokenPayload)
}

func (s *ExternalAuthService) findOrCreateGoogleUser(tokenInfo *oauth2.Tokeninfo) (*model.User, error) {
	if user, err := s.userRepo.GetUserByGoogleID(tokenInfo.UserId); err == nil {
		return user, nil
	}

	if user, err := s.userRepo.GetUserByEmail(tokenInfo.Email); err == nil {
		user.GoogleID = &tokenInfo.UserId
		if err = s.userRepo.UpdateUser(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	password, err := generatePassword(constants.PwdGenerateSize)
	if err != nil {
		return nil, err
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:        tokenInfo.Email,
		GoogleID:     &tokenInfo.UserId,
		PasswordHash: hashPassword,
	}

	err = s.codeRepo.PublishGeneratedPassword(context.Background(), tokenInfo.Email, password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *ExternalAuthService) findOrCreateAppleUser(tokenPayload *AppleIDTokenPayload) (*model.User, error) {
	if user, err := s.userRepo.GetUserByAppleID(tokenPayload.Sub); err == nil {
		return user, nil
	}

	if user, err := s.userRepo.GetUserByEmail(tokenPayload.Email); err == nil {
		user.AppleID = &tokenPayload.Sub
		if err = s.userRepo.UpdateUser(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	password, err := generatePassword(constants.PwdGenerateSize)
	if err != nil {
		return nil, err
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:        tokenPayload.Email,
		GoogleID:     &tokenPayload.Sub,
		PasswordHash: hashPassword,
	}

	err = s.codeRepo.PublishGeneratedPassword(context.Background(), tokenPayload.Email, password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getApplePublicKey(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keyResp AppleKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&keyResp); err != nil {
		return nil, err
	}

	for _, key := range keyResp.Keys {
		if key.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}

			n := new(big.Int).SetBytes(nBytes)

			var e int
			if len(eBytes) < 8 {
				eBytes = append(make([]byte, 8-len(eBytes)), eBytes...)
			}
			e = int(binary.BigEndian.Uint64(eBytes))

			pubKey := &rsa.PublicKey{
				N: n,
				E: e,
			}
			return pubKey, nil
		}
	}

	return nil, errs.InvalidAppleToken
}

func validateAppleIDToken(idToken string, clientID string) (*AppleIDTokenPayload, error) {
	parser := jwt.NewParser()
	token, parts, err := parser.ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if len(parts) < 2 {
		return nil, errs.InvalidAppleToken
	}

	header, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errs.InvalidAppleToken
	}

	publicKey, err := getApplePublicKey(header)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok = token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errs.InvalidAppleToken
		}
		return publicKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, errs.InvalidAppleToken
	}

	if claims["iss"] != "https://appleid.apple.com" || claims["aud"] != clientID {
		return nil, errs.InvalidAppleToken
	}

	payload := &AppleIDTokenPayload{
		Sub:   claims["sub"].(string),
		Email: claims["email"].(string),
	}
	return payload, nil
}

func validateGoogleToken(token string) (*oauth2.Tokeninfo, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(token).Do()
	if err != nil {
		return nil, errs.InvalidGoogleToken
	}

	if tokenInfo.Email == "" || !tokenInfo.VerifiedEmail {
		return nil, errs.InvalidGoogleToken
	}
	return tokenInfo, nil
}

func generatePassword(length int) (string, error) {
	const lower = "abcdefghijklmnopqrstuvwxyz"
	const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const digits = "0123456789"
	const all = lower + upper + digits

	if length < 3 {
		return "", nil
	}

	password := make([]byte, length)

	sets := []string{lower, upper, digits}
	for i := 0; i < 3; i++ {
		char, err := randomCharFromSet(sets[i])
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	for i := 3; i < length; i++ {
		char, err := randomCharFromSet(all)
		if err != nil {
			return "", err
		}
		password[i] = char
	}
	shuffle(password)
	return string(password), nil
}

func randomCharFromSet(set string) (byte, error) {
	ma := big.NewInt(int64(len(set)))
	n, err := rand.Int(rand.Reader, ma)
	if err != nil {
		return 0, err
	}
	return set[n.Int64()], nil
}

func shuffle(data []byte) {
	for i := len(data) - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		data[i], data[j] = data[j], data[i]
	}
}
