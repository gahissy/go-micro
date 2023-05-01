package h

import (
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/o1egl/paseto"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func IsNil(v interface{}) bool {
	return v == nil
}

func IsNotNil(v interface{}) bool {
	return v != nil
}

func IsStrEmpty(v string) bool {
	return len(strings.TrimSpace(v)) == 0
}

func IsPointer(arg interface{}) bool {
	argType := reflect.TypeOf(arg)
	return argType.Kind() == reflect.Ptr
}

func ToInt(input string) int {
	if input == "" {
		return 0
	}
	res, _ := strconv.Atoi(input)
	return res
}

func ToFloat(input string) float32 {
	if input == "" {
		return 0
	}
	res, err := strconv.ParseFloat(input, 32)
	if err != nil {
		log.Error(err)
		return 0
	}
	return float32(res)
}

func Contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func CreatePasetoToken(audience string, subject string) (string, error) {

	issuer := os.Getenv("ENCRYPTION_DOMAIN")
	secret := []byte(os.Getenv("ENCRYPTION_KEY"))
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   audience,
		Issuer:     issuer,
		Subject:    subject,
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	return paseto.NewV2().Encrypt(secret, &jsonToken, []byte(""))

}

func CreateJwtToken(audience string, subject string) (string, error) {
	//issuer := os.Getenv("ENCRYPTION_DOMAIN")
	secret := []byte(os.Getenv("ENCRYPTION_KEY"))
	claims := jwt.MapClaims{}
	claims["sub"] = subject
	claims["role"] = audience
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)

}

func ToDate(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}
	res, err := time.Parse("2006-01-02", input)
	return &res, err
}

func ToJsonString(input interface{}) (string, error) {
	if jsonBytes, err := json.Marshal(input); err != nil {
		log.Error("Error marshaling JSON:", err)
		return "", err
	} else {
		return string(jsonBytes), nil
	}
}

func GetUpdatedFields[T interface{}](p1, p2 *T) map[string]interface{} {

	updatedFields := make(map[string]interface{})
	v1 := reflect.ValueOf(p1).Elem()
	v2 := reflect.ValueOf(p2).Elem()
	for i := 0; i < v1.NumField(); i++ {
		field1 := v1.Field(i)
		field2 := v2.Field(i)
		if !reflect.DeepEqual(field1.Interface(), field2.Interface()) {
			updatedFields[v1.Type().Field(i).Name] = field2.Interface()
		}
	}
	return updatedFields
}

func Now() time.Time {
	return time.Now().UTC()
}

func EncoreBase64(message string) string {
	messageBytes := []byte(message)
	encoded := base64.StdEncoding.EncodeToString(messageBytes)
	return encoded
}
