package googlesearchconsole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	types "github.com/Leapforce-nl/go_types"
	"github.com/getsentry/sentry-go"
)

var tokenMutex sync.Mutex

type Token struct {
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Expiry       time.Time
}

type ApiError struct {
	Error       string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

func LockToken() {
	tokenMutex.Lock()
}

func UnlockToken() {
	tokenMutex.Unlock()
}

func (t *Token) Useable() bool {
	if t == nil {
		return false
	}
	if t.AccessToken == "" || t.RefreshToken == "" {
		return false
	}
	return true
}

func (t *Token) Refreshable() bool {
	if t == nil {
		return false
	}
	if t.RefreshToken == "" {
		return false
	}
	return true
}

func (t *Token) IsExpired() (bool, error) {
	if !t.Useable() {
		return true, &types.ErrorString{"Token is not valid."}
	}
	if t.Expiry.Add(-60 * time.Second).Before(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (gsc *GoogleSearchConsole) GetToken(url string, hasRefreshToken bool) error {
	guid := types.NewGUID()
	fmt.Println("GetTokenGUID:", guid)
	fmt.Println(url)

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		return err
	}

	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("GetTokenGUID:", guid)
		fmt.Println("AccessToken:", gsc.Token.AccessToken)
		fmt.Println("Refresh:", gsc.Token.RefreshToken)
		fmt.Println("Expiry:", gsc.Token.Expiry)
		fmt.Println("Now:", time.Now())

		eoError := ApiError{}

		err = json.Unmarshal(b, &eoError)
		if err != nil {
			return err
		}

		message := fmt.Sprintln("Error:", eoError.Error, ", ", eoError.Description)
		fmt.Println(message)

		if res.StatusCode == 401 {
			if gsc.IsLive {
				sentry.CaptureMessage("GoogleSearchConsole refreshtoken not valid, login needed to retrieve a new one. Error: " + message)
			}
			gsc.InitToken()
		}

		return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v, url: %s", res.StatusCode, req.URL)}
	}

	token := Token{}

	err = json.Unmarshal(b, &token)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(token)

	/*if gsc.Token != nil {
		fmt.Println("old token:")
		fmt.Println(gsc.Token.AccessToken)
		fmt.Println("old refresh token:")
		fmt.Println(gsc.Token.RefreshToken)
		fmt.Println("old expiry:")
		fmt.Println(gsc.Token.Expiry)
	}*/

	token.Expiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	if gsc.Token == nil {
		gsc.Token = &Token{}
	}

	gsc.Token.Expiry = token.Expiry
	gsc.Token.AccessToken = token.AccessToken

	if hasRefreshToken {
		gsc.Token.RefreshToken = token.RefreshToken

		err = gsc.SaveTokenToBigQuery()
		if err != nil {
			return err
		}
	}

	fmt.Println("new token:")
	fmt.Println(gsc.Token.AccessToken)
	fmt.Println("new refresh token:")
	fmt.Println(gsc.Token.RefreshToken)
	fmt.Println("new expiry:")
	fmt.Println(gsc.Token.Expiry)
	fmt.Println("GetTokenGUID:", guid)

	return nil
}

func (gsc *GoogleSearchConsole) GetTokenFromCode(code string) error {
	//fmt.Println("GetTokenFromCode")
	url := fmt.Sprintf("%s?code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&scope=&grant_type=authorization_code", gsc.TokenURL, code, gsc.RedirectURL, gsc.ClientID, gsc.ClientSecret)

	return gsc.GetToken(url, true)
}

func (gsc *GoogleSearchConsole) GetTokenFromRefreshToken() error {
	fmt.Println("***GetTokenFromRefreshToken***")

	//always get refresh token from BQ prior to using it
	gsc.GetTokenFromBigQuery()

	if !gsc.Token.Refreshable() {
		return gsc.InitToken()
	}

	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&refresh_token=%s&grant_type=refresh_token&access_type=offline&prompt=consent", gsc.TokenURL, gsc.ClientID, gsc.ClientSecret, gsc.Token.RefreshToken)

	return gsc.GetToken(url, false)
}

func (gsc *GoogleSearchConsole) ValidateToken() error {
	LockToken()
	defer UnlockToken()

	if !gsc.Token.Useable() {

		err := gsc.GetTokenFromRefreshToken()
		if err != nil {
			return err
		}

		if !gsc.Token.Useable() {
			if gsc.IsLive {
				sentry.CaptureMessage("GoogleSearchConsole refreshtoken not found or empty, login needed to retrieve a new one.")
			}
			err := gsc.InitToken()
			if err != nil {
				return err
			}
			//return &types.ErrorString{""}
		}
	}

	isExpired, err := gsc.Token.IsExpired()
	if err != nil {
		return err
	}
	if isExpired {
		//fmt.Println(time.Now(), "[token expired]")
		err = gsc.GetTokenFromRefreshToken()
		if err != nil {
			return err
		}
	}

	return nil
}

func (gsc *GoogleSearchConsole) InitToken() error {

	if gsc == nil {
		return &types.ErrorString{"GoogleSearchConsole variable not initialized"}
	}

	url := fmt.Sprintf("%s?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&access_type=offline&prompt=consent", gsc.AuthURL, gsc.ClientID, gsc.RedirectURL, "https://www.googleapis.com/auth/webmasters")

	fmt.Println("Go to this url to get new access token:\n")
	fmt.Println(url + "\n")

	// Create a new redirect route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		//
		// get authorization code
		//
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		fmt.Println(code)

		err = gsc.GetTokenFromCode(code)
		if err != nil {
			fmt.Println(err)
		}

		w.WriteHeader(http.StatusFound)

		return
	})

	http.ListenAndServe(":8080", nil)

	return nil
}
