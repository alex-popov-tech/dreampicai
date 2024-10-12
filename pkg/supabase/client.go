package supabase

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	gtt "github.com/supabase-community/gotrue-go/types"
	s "github.com/supabase-community/supabase-go"
)

var Client *s.Client

func CreateClient(url, key string) (*s.Client, error) {
	c, err := s.NewClient(url, key, &s.ClientOptions{})
	// init global supabase client for later usage across the app
	// if err != nil app won't start, so its safe to use global here
	Client = c
	return c, err
}

// real use case, sometimes client returns something like this:
// response status code 400: {"code":400,"error_code":"invalid_credentials","msg":"Invalid login credentials"}
// and this func grabs 'msg' if its there, or just return initial err.Error()
func TryGerSupabaseErrorMessage(supabaseClientError error) error {
	re := regexp.MustCompile(`"msg":"(.*?)"`)
	matches := re.FindStringSubmatch(supabaseClientError.Error())

	if len(matches) > 1 {
		return errors.New(matches[1])
	}

	return supabaseClientError
}

func GetSessionFromCallbackCookie(cookies [](*http.Cookie)) (gtt.Session, error) {
	var authCookie *http.Cookie
	for _, it := range cookies {
		fmt.Printf("cookie: %v", it)
		if strings.Contains(it.Name, "auth-token") {
			authCookie = it
			break
		}
	}
	if authCookie == nil {
		return gtt.Session{}, errors.New("'auth-token' cookie not found")
	}

	cookieValue := authCookie.Value                // base64-eyJhY2Nlc3NfdG9rZW4iOiJleUpoYkdjaU9pSkl...VekkxTmlJc0ltdHaQ0k2SW05RFpHUX19;
	encoded := cookieValue[7 : len(cookieValue)-1] // eyJhY2Nlc3NfdG9rZW4iOiJleUpoYkdjaU9pSkl...VekkxTmlJc0ltdHaQ0k2SW05RFpHUX19
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return gtt.Session{}, errors.New("Cannot decode auth cookie: " + err.Error())
	}

	var session gtt.Session
	err = json.Unmarshal(decodedBytes, &session)
	if err != nil {
		return gtt.Session{}, errors.New("Cannot parse auth cookie as json: " + err.Error())
	}

	return session, nil
}

func GetSessionFromQuery(values url.Values) (gtt.Session, error) {
	session := gtt.Session{}

	if val := values.Get("access_token"); val == "" {
		return session, fmt.Errorf("Cannot get 'access_token' from query: %v", values)
	} else {
		session.AccessToken = val
	}

	if val := values.Get("refresh_token"); val == "" {
		return session, fmt.Errorf("Cannot get 'refresh_token' from query: %v", values)
	} else {
		session.RefreshToken = val
	}

	if val := values.Get("token_type"); val == "" {
		return session, fmt.Errorf("Cannot get 'token_type' from query: %v", values)
	} else {
		session.TokenType = val
	}

	if str := values.Get("expires_at"); str == "" {
		return session, fmt.Errorf("Cannot get 'expires_at' from query: %v", values)
	} else {
		val, err := strconv.Atoi(str)
		if err != nil {
			return session, fmt.Errorf("Cannot get or parse 'expires_at' from query: %v", values)
		}
		session.ExpiresAt = int64(val)
	}

	if str := values.Get("expires_in"); str == "" {
		return session, fmt.Errorf("Cannot get 'expires_in' from query: %v", values)
	} else {
		val, err := strconv.Atoi(str)
		if err != nil {
			return session, fmt.Errorf("Cannot get or parse 'expires_in' from query: %v", values)
		}
		session.ExpiresIn = val
	}

	return session, nil
}
