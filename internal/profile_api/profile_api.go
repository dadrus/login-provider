package profile_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"login-provider/internal/utils"
	"net/http"
	"time"
)

type Address struct {
	Street  string `json:"street" mapstructure:"street"`
	City    string `json:"city" mapstructure:"city"`
	Zip     string `json:"zip" mapstructure:"zip"`
	State   string `json:"state" mapstructure:"state"`
	Country string `json:"country" mapstructure:"country"`
}

type User struct {
	ID          int        `json:"id" mapstructure:"id"`
	FirstName   string     `json:"first_name" mapstructure:"first_name"`
	LastName    string     `json:"last_name" mapstructure:"last_name"`
	UserName    string     `json:"user_name" mapstructure:"user_name"`
	Gender      string     `json:"gender" mapstructure:"gender"`
	Birthday    *time.Time `time_format:"2006-01-02" json:"birthday" mapstructure:"birthday"`
	Address     *Address   `json:"address" mapstructure:"address"`
	Email       string     `json:"email" mapstructure:"email"`
	PhoneNumber string     `json:"phone" mapstructure:"phone"`
}

type AuthenticationRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	ProfileUrl string `json:"profile_url" mapstructure:"profile_url"`
	User       User   `json:"user" mapstructure:"user"`
}

func (ar AuthenticationResponse) CreateIdTokenClaims(grantedScopes []string) map[string]string {
	claims := make(map[string]string)

	if utils.Contains(grantedScopes, "profile") {
		claims["profile"] = ar.ProfileUrl
		claims["name"] = ar.User.UserName
		claims["family_name"] = ar.User.LastName
		claims["given_name"] = ar.User.FirstName
		claims["preferred_username"] = ar.User.UserName
		claims["gender"] = ar.User.Gender
		if ar.User.Birthday != nil {
			claims["birthdate"] = ar.User.Birthday.Format("2006-01-02")
		}
		claims["zoneinfo"] = time.Local.String()
		claims["locale"] = "de_DE"
		claims["updated_at"] = time.Now().String()
	}

	if utils.Contains(grantedScopes, "email") {
		claims["email"] = ar.User.Email
		claims["email_verified"] = "true"
	}

	if utils.Contains(grantedScopes, "address") {
		claims["street_address"] = ar.User.Address.Street
		claims["locality"] = ar.User.Address.City
		claims["region"] = ar.User.Address.State
		claims["postal_code"] = ar.User.Address.Zip
		claims["country"] = ar.User.Address.Country
	}

	if utils.Contains(grantedScopes, "phone") {
		claims["phone_number"] = ar.User.PhoneNumber
		claims["phone_number_verified"] = "false"
	}

	return claims
}

func (ar *AuthenticationResponse) Unmarshal(data interface{}) error {
	return mapstructure.Decode(data, ar)
}

func AuthenticateUser(url string, userName, password string) (*AuthenticationResponse, error) {
	jsonValue, _ := json.Marshal(AuthenticationRequest{UserName: userName, Password: password})

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("unexpected status code")
	}

	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResponse AuthenticationResponse
	if err := json.Unmarshal(rawData, &authResponse); err != nil {
		return nil, err
	}

	return &authResponse, nil
}