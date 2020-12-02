package googlesearchconsole

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"
	errortools "github.com/leapforce-libraries/go_errortools"

	go_oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName string = "GoogleSearchConsole"
	//apiURL          string = "https://www.googleapis.com/calendar/v3"
	authURL         string = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL        string = "https://oauth2.googleapis.com/token"
	tokenHTTPMethod string = http.MethodPost
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// GoogleSearchConsole stores GoogleSearchConsole configuration
//
type GoogleSearchConsole struct {
	SiteURL string
	baseURL string
	oAuth2  *go_oauth2.OAuth2
}

// methods
//
func NewGoogleSearchConsole(baseURL string, clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*GoogleSearchConsole, error) {
	gsc := GoogleSearchConsole{}
	gsc.baseURL = baseURL

	maxRetries := uint(3)
	config := go_oauth2.OAuth2Config{
		ApiName:         apiName,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Scope:           scope,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHTTPMethod,
		MaxRetries:      &maxRetries,
	}
	gsc.oAuth2 = go_oauth2.NewOAuth(config, bigQuery, isLive)

	return &gsc, nil
}

func (gsc *GoogleSearchConsole) Validate() *errortools.Error {
	if gsc.baseURL == "" {
		return errortools.ErrorMessage(fmt.Sprintf("%s baseURL not provided", apiName))
	}
	if gsc.SiteURL == "" {
		return errortools.ErrorMessage(fmt.Sprintf("%s SiteURL not provided", apiName))
	}

	if !strings.HasSuffix(gsc.baseURL, "/") {
		gsc.baseURL = gsc.baseURL + "/"
	}

	if !strings.HasSuffix(gsc.SiteURL, "/") {
		gsc.SiteURL = gsc.SiteURL + "/"
	}

	return nil
}

func (gsc *GoogleSearchConsole) post(url string, buf *bytes.Buffer, model interface{}) (*http.Request, *http.Response, *errortools.Error) {
	err := GoogleSearchControlError{}
	request, response, e := gsc.oAuth2.Post(url, buf, model, &err)
	if e != nil {
		if err.Err.Message != "" {
			e.SetMessage(err.Err.Message)
		}

		return request, response, e
	}

	return request, response, nil
}

/*
func (gsc *GoogleSearchConsole) GetHttpClient() (*http.Client, *errortools.Error) {

	_, e := gsc.oAuth2.ValidateToken()
	if e != nil {
		return nil, e
	}

	return new(http.Client), nil
}


func (gsc *GoogleSearchConsole) Post(url string, values map[string]string, model interface{}) *errortools.Error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return gsc.PostBuffer(url, buf, model)
}

func (gsc *GoogleSearchConsole) PostBytes(url string, b []byte, model interface{}) *errortools.Error {
	return gsc.PostBuffer(url, bytes.NewBuffer(b), model)
}

func (gsc *GoogleSearchConsole) PostBuffer(url string, buf *bytes.Buffer, model interface{}) *errortools.Error {
	client, errClient := gsc.GetHttpClient()
	if errClient != nil {
		return errClient
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}

	gsc.oAuth2.LockToken()

	// Add authorization token to header
	var bearer = "Bearer " + gsc.oAuth2.Token.AccessToken
	req.Header.Add("authorization", bearer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	gsc.oAuth2.UnlockToken()
	if err != nil {
		fmt.Println("errDo")
		return err
	}

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("ERROR in Post:", url)
		//fmt.Println(url)
		fmt.Println("StatusCode", res.StatusCode)
		//fmt.Println(gsc.Token.AccessToken)
		return gsc.PrintError(res)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &model)
	if err != nil {
		fmt.Println("errUnmarshal1")
		return gsc.PrintError(res)
	}

	return nil
}

func (gsc *GoogleSearchConsole) PrintError(res *http.Response) *errortools.Error {
	fmt.Println("Status", res.Status)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("errUnmarshal1")
		return err
	}

	ee := GoogleSearchControlError{}

	err = json.Unmarshal(b, &ee)
	if err != nil {
		fmt.Println("errUnmarshal1")
		return err
	}

	message := fmt.Sprintf("Server returned statuscode %v, error:%s", res.StatusCode, ee.Err.Message)
	return &types.ErrorString{message}
}
*/
