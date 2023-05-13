package h

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"github.com/o1egl/paseto"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
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

func CreatePasetoToken(secret []byte, issuer string, audience string, subject string) (string, error) {

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

func CreateJwtToken(secret string, audience string, subject string) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = subject
	claims["role"] = audience
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Errorf("Error signing token: %v", err)
	}
	return signed, err

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

func CopyFields(dst, src interface{}) error {
	return copier.Copy(dst, src)
}

func Diff(original, updated interface{}) map[string]interface{} {
	originalValue := reflect.ValueOf(original)
	updatedValue := reflect.ValueOf(updated)

	// Check if the original and update values are pointers, and dereference them if they are
	if originalValue.Kind() == reflect.Ptr {
		originalValue = originalValue.Elem()
	}
	if updatedValue.Kind() == reflect.Ptr {
		updatedValue = updatedValue.Elem()
	}

	diff := make(map[string]interface{})

	for i := 0; i < originalValue.NumField(); i++ {
		originalField := originalValue.Field(i)
		updatedField := updatedValue.FieldByName(originalValue.Type().Field(i).Name)

		// If the update field is present, non-nil, and its value is different from the original field's value, add the field to the diff map
		if updatedField.IsValid() {
			if updatedField.Kind() == reflect.Ptr {
				if !updatedField.IsNil() && !reflect.DeepEqual(originalField.Interface(), updatedField.Elem().Interface()) {
					diff[ToSnakeCase(originalValue.Type().Field(i).Name)] = updatedField.Elem().Interface()
				}
			} else if !reflect.DeepEqual(originalField.Interface(), updatedField.Interface()) {
				diff[ToSnakeCase(originalValue.Type().Field(i).Name)] = updatedField.Interface()
			}
		}
	}

	return diff
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

func FindFileInParents(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		filename := filepath.Join(dir, name)
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}

		// Move up to the parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached the root directory, stop searching
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find %s file in parent directories", name)
}

func ToSnakeCase(str string) string {
	var output []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			output = append(output, '_')
		}
		output = append(output, r)
	}
	return strings.ToLower(string(output))
}
