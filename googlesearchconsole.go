package googlesearchconsole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	types "github.com/Leapforce-nl/go_types"

	googleoauth2 "github.com/Leapforce-nl/go_googleoauth2"
)

const apiName string = "GoogleSearchConsole"

// GoogleSearchConsole stores GoogleSearchConsole configuration
//
type GoogleSearchConsole struct {
	SiteURL string
	BaseURL string
	oAuth2  *googleoauth2.GoogleOAuth2
}

// methods
//
func (gsc *GoogleSearchConsole) InitOAuth2(clientID string, clientSecret string, scopes []string, bigQuery *bigquerytools.BigQuery, isLive bool) error {
	_oAuth2 := new(googleoauth2.GoogleOAuth2)
	_oAuth2.ApiName = apiName
	_oAuth2.ClientID = clientID
	_oAuth2.ClientSecret = clientSecret
	_oAuth2.Scopes = scopes
	_oAuth2.BigQuery = bigQuery
	_oAuth2.IsLive = isLive

	gsc.oAuth2 = _oAuth2

	return nil
}

func (gsc *GoogleSearchConsole) Validate() error {
	if gsc.BaseURL == "" {
		return &types.ErrorString{fmt.Sprintf("%s BaseURL not provided", apiName)}
	}
	if gsc.SiteURL == "" {
		return &types.ErrorString{fmt.Sprintf("%s SiteURL not provided", apiName)}
	}

	if !strings.HasSuffix(gsc.BaseURL, "/") {
		gsc.BaseURL = gsc.BaseURL + "/"
	}

	if !strings.HasSuffix(gsc.SiteURL, "/") {
		gsc.SiteURL = gsc.SiteURL + "/"
	}

	return nil
}

func (gsc *GoogleSearchConsole) GetHttpClient() (*http.Client, error) {
	/*err := gsc.Wait()
	if err != nil {
		return nil, err
	}*/

	err := gsc.oAuth2.ValidateToken()
	if err != nil {
		return nil, err
	}

	return new(http.Client), nil
}

func (gsc *GoogleSearchConsole) Post(url string, values map[string]string, model interface{}) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return gsc.PostBuffer(url, buf, model)
}

func (gsc *GoogleSearchConsole) PostBytes(url string, b []byte, model interface{}) error {
	return gsc.PostBuffer(url, bytes.NewBuffer(b), model)
}

func (gsc *GoogleSearchConsole) PostBuffer(url string, buf *bytes.Buffer, model interface{}) error {
	client, errClient := gsc.GetHttpClient()
	if errClient != nil {
		return errClient
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		fmt.Println("errNewRequest")
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

func (gsc *GoogleSearchConsole) PrintError(res *http.Response) error {
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
