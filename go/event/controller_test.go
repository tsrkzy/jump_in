package event

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsrkzy/jump_in/authenticate"
	"gopkg.in/resty.v1"
	"net/http"
	"testing"
)

const TestDebug = false
const DEFAULT_EMAIL = "tsrmix+echo@gmail.com"

func TestCreate001(t *testing.T) {

	respMl1, w, err := Login(t, DEFAULT_EMAIL)
	accountId := fmt.Sprintf("%s", w.ID)
	assert.NoError(t, err)

	cEc := resty.New().SetDebug(TestDebug)
	cEc.SetCookies(respMl1.Cookies())
	reqEc := CreateRequest{
		Name:      "テスト用イベント名",
		AccountId: accountId,
	}
	resEc := CreateResponse{}
	respEc, err := cEc.R().
		SetBody(reqEc).
		SetResult(&resEc).
		Post("http://localhost:80/api/event/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEc.StatusCode())

	cEl := resty.New().SetDebug(TestDebug)
	cEl.SetCookies(respMl1.Cookies())
	resEl := ListResponse{}
	respEl, err := cEl.R().
		SetQueryParams(map[string]string{
			"account_id": accountId,
		}).
		SetResult(&resEl).
		Get("http://localhost:80/api/event/list")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEl.StatusCode())

	aId := ""
	for i, r := range resEl.EventsOwns {
		if i == 0 {
			aId = r.AccountID
			fmt.Printf("account_id: %s\n", aId)
		} else {
			assert.Equal(t, aId, r.AccountID)
		}
	}
}

func TestAttend001(t *testing.T) {
	respMl, w, err := Login(t, DEFAULT_EMAIL)
	accountId := fmt.Sprintf("%s", w.ID)
	assert.NoError(t, err)

	cEc := resty.New().SetDebug(TestDebug)
	cEc.SetCookies(respMl.Cookies())
	reqEc := CreateRequest{
		Name:      "Attendテスト用イベント名",
		AccountId: accountId,
	}
	resEc := CreateResponse{}
	respEc, err := cEc.R().
		SetBody(reqEc).
		SetResult(&resEc).
		Post("http://localhost:80/api/event/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEc.StatusCode())

	eventId := fmt.Sprintf("%s", resEc.ID)

	/* 参加 */
	cEa1 := resty.New().SetDebug(TestDebug)
	cEa1.SetCookies(respMl.Cookies())
	reqAtt1 := AttendRequest{EventId: eventId, AccountId: accountId}
	respAtt1, err := cEa1.R().
		SetBody(reqAtt1).
		//SetResult(reqAtt1).
		Post("http://localhost:80/api/event/attend")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAtt1.StatusCode())

	/* 二重参加 */
	cEa2 := resty.New().SetDebug(TestDebug)
	cEa2.SetCookies(respMl.Cookies())
	reqAtt2 := AttendRequest{EventId: eventId, AccountId: accountId}
	respAtt2, err := cEa2.R().
		SetBody(reqAtt2).
		//SetResult(reqAtt2).
		Post("http://localhost:80/api/event/attend")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAtt2.StatusCode())

	/* 参加取消 */
	cEl1 := resty.New().SetDebug(TestDebug)
	cEl1.SetCookies(respMl.Cookies())
	reqLe1 := AttendRequest{EventId: eventId, AccountId: accountId}
	respLe1, err := cEl1.R().
		SetBody(reqLe1).
		//SetResult(reqLe1).
		Post("http://localhost:80/api/event/leave")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLe1.StatusCode())

	/* 参加二重取消 */
	cEl2 := resty.New().SetDebug(TestDebug)
	cEl2.SetCookies(respMl.Cookies())
	reqLe2 := AttendRequest{EventId: eventId, AccountId: accountId}
	respLe2, err := cEl2.R().
		SetBody(reqLe2).
		//SetResult(reqLe2).
		Post("http://localhost:80/api/event/leave")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLe2.StatusCode())

}

func Login(t *testing.T, email string) (*resty.Response, authenticate.WhoAmIResponse, error) {
	redirectUri := "http://localhost:80/api/status"
	r := authenticate.Request{
		MailAddress: email,
		RedirectURI: redirectUri,
	}
	/* ML取得、hrefへリダイレクトさせる設定 */
	cA := resty.New().SetDebug(TestDebug)
	authResult := authenticate.Result{}
	respA, err := cA.R().
		SetBody(r).
		SetResult(&authResult).
		Post("http://localhost:80/api/authenticate")
	assert.NoError(t, err)

	/* NoRedirectPolicyをセットして/mlからのレスポンス検証 */
	cMl1 := resty.New().
		SetHTTPMode().
		SetDebug(TestDebug)
	cMl1.SetCookies(respA.Cookies())
	respMl1, err := cMl1.R().
		Get(authResult.MagicLink)
	assert.NoError(t, err)

	cW := resty.New().SetDebug(TestDebug)
	cW.SetCookies(respMl1.Cookies())
	cR := authenticate.WhoAmIResponse{}
	respW, err := cW.R().
		SetResult(&cR).
		Get("http://localhost:80/api/whoami")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respW.StatusCode())

	return respMl1, cR, err
}
