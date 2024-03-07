package tglogin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	CookieName      = "X-telegram-data"
	CookieSeparator = "&"
)

type TelegramUserData struct {
	TGID      int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"` // UNIX timestamp
	Hash      string `json:"hash"`
}

// FromJSON TelegramUserData from json
// It doesn't check is data valid, all invalid fields will be ignored.
// If data is invalid, user will not pass authorisation check.
func FromJSON(body io.Reader) (data TelegramUserData) {
	json.NewDecoder(body).Decode(&data)
	return
}

// FromQuery TelegramUserData from query
// It doesn't check is data valid, all invalid fields will be ignored.
// If data is invalid, user will not pass authorisation check.
func FromQuery(query url.Values) (data TelegramUserData) {
	data.TGID, _ = strconv.ParseInt(query.Get("id"), 10, 64)
	data.FirstName = query.Get("first_name")
	data.LastName = query.Get("last_name")
	data.Username = query.Get("username")
	data.PhotoURL = query.Get("photo_url")
	data.AuthDate, _ = strconv.ParseInt(query.Get("auth_date"), 10, 64)
	data.Hash = query.Get("hash")
	return
}

func FromCookie(r *http.Request) (data TelegramUserData, err error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil || cookie.Value == "" {
		return TelegramUserData{}, err
	}
	s, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return TelegramUserData{}, err
	}
	splitted := strings.Split(s, CookieSeparator)
	for _, v := range splitted {
		key, value, _ := strings.Cut(v, "=")
		switch key {
		case "id":
			data.TGID, _ = strconv.ParseInt(value, 10, 64)
		case "first_name":
			data.FirstName = value
		case "last_name":
			data.LastName = value
		case "username":
			data.Username = value
		case "photo_url":
			data.PhotoURL = value
		case "auth_date":
			data.AuthDate, _ = strconv.ParseInt(value, 10, 64)
		case "hash":
			data.Hash = value
		}
	}
	return
}

func SetCookie(w http.ResponseWriter, user TelegramUserData) {
	values := pairs(user)
	values = append(values, fmt.Sprintf("hash=%s", user.Hash))
	value := strings.Join(values, CookieSeparator)
	cookie := &http.Cookie{
		Name:  CookieName,
		Value: url.QueryEscape(value),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  CookieName,
		Value: "",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

// IsTelegramAuthorization
// https://core.telegram.org/widgets/login#checking-authorization
func IsTelegramAuthorization(user TelegramUserData, token string) bool {
	// Concate field to check string
	checkString := createCheckString(user)

	// Create encoder based on bot token
	encoder := hmac_sha256Encoder(token)

	// Encode check string by encoder.
	encoder.Write([]byte(checkString))
	encodedCheckString := encoder.Sum(nil)

	// Compare encodedCheckString with hash
	return hex.EncodeToString(encodedCheckString) == user.Hash
}

// IsExpiredDate Check ttl of telegram data.
func IsExpiredData(authDate int64, ttl time.Duration) bool {
	loginTime := time.Unix(authDate, 0).UTC()
	sinceFromLogin := time.Since(loginTime)
	return sinceFromLogin < ttl
}

// createCheckString Concate fields to check string.
// Concate in alphabet sorted.
func createCheckString(user TelegramUserData) string {
	params := pairs(user)
	return strings.Join(params, "\n")
}

func pairs(user TelegramUserData) []string {
	params := make([]string, 0, 6)
	params = append(params, fmt.Sprintf("auth_date=%d", user.AuthDate))
	params = append(params, fmt.Sprintf("first_name=%s", user.FirstName))
	params = append(params, fmt.Sprintf("id=%d", user.TGID))

	if user.LastName != "" {
		params = append(params, fmt.Sprintf("last_name=%s", user.LastName))
	}
	if user.PhotoURL != "" {
		params = append(params, fmt.Sprintf("photo_url=%s", user.PhotoURL))
	}
	if user.Username != "" {
		params = append(params, fmt.Sprintf("username=%s", user.Username))
	}
	return params
}

// hmac_sha256Encoder generate new encoder based on token
func hmac_sha256Encoder(token string) hash.Hash {
	sha256Encoder := sha256.New()
	sha256Encoder.Write([]byte(token))
	secretKey := sha256Encoder.Sum(nil)
	return hmac.New(sha256.New, secretKey)
}
