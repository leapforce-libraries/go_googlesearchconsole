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
)

// GoogleSearchConsole stores GoogleSearchConsole configuration
//
type GoogleSearchConsole struct {
	// config
	ClientID     string
	ClientSecret string
	ClientURL    string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	BaseURL      string
	Token        *Token
	BigQuery     *bigquerytools.BigQuery
	IsLive       bool
}

// methods
//
func (gsc *GoogleSearchConsole) Init() error {
	if gsc.BaseURL == "" {
		return &types.ErrorString{"GoogleSearchConsole BaseURL not provided"}
	}
	if gsc.ClientURL == "" {
		return &types.ErrorString{"GoogleSearchConsole ClientURL not provided"}
	}

	if !strings.HasSuffix(gsc.BaseURL, "/") {
		gsc.BaseURL = gsc.BaseURL + "/"
	}

	if !strings.HasSuffix(gsc.ClientURL, "/") {
		gsc.ClientURL = gsc.ClientURL + "/"
	}

	return nil
}

func (gsc *GoogleSearchConsole) GetHttpClient() (*http.Client, error) {
	/*err := gsc.Wait()
	if err != nil {
		return nil, err
	}*/

	err := gsc.ValidateToken()
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

	LockToken()

	// Add authorization token to header
	var bearer = "Bearer " + gsc.Token.AccessToken
	req.Header.Add("authorization", bearer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	UnlockToken()
	if err != nil {
		fmt.Println("errDo")
		return err
	}

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("ERROR in Post")
		fmt.Println(url)
		fmt.Println("StatusCode", res.StatusCode)
		fmt.Println(gsc.Token.AccessToken)
		return gsc.PrintError(res)
	}

	defer res.Body.Close()

	fmt.Println("res.Body", res.Body)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

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
