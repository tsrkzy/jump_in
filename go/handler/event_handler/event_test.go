package event_handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsrkzy/jump_in/helper/testhelper"
	"github.com/tsrkzy/jump_in/types/authenticate_types"
	"github.com/tsrkzy/jump_in/types/candidate_types"
	"github.com/tsrkzy/jump_in/types/event_types"
	"gopkg.in/resty.v1"
	"net/http"
	"testing"
)

const testDebug = false
const DEFAULT_EMAIL = "tsrmix+echo@gmail.com"

func TestCreate001(t *testing.T) {

	respMl1, w, err := Login(t, DEFAULT_EMAIL)
	accountId := fmt.Sprintf("%s", w.ID)
	assert.NoError(t, err)

	cEc := testhelper.MakeClient(respMl1)
	reqEc := event_types.CreateRequest{
		Name:        "テスト用イベント名",
		Description: "/Users/tsrkzy/dev/go/github.com/tsrkzy/jump_in/go/handler/event_handler/event_test.go - Event 001",
		AccountID:   accountId,
	}
	resEc := event_types.CreateResponse{}
	respEc, err := cEc.R().
		SetBody(reqEc).
		SetResult(&resEc).
		Post("http://localhost:80/api/event/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEc.StatusCode())

	cEl := testhelper.MakeClient(respMl1)
	resEl := event_types.ListResponse{}
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

	eventNameUpdated := "test_event_name_updated"
	cEn := testhelper.MakeClient(respMl1)
	reqEn := event_types.UpdateNameRequest{
		EventID:   resEc.Event.ID,
		EventName: eventNameUpdated,
	}
	resEn := event_types.UpdateNameResponse{}
	respEn, err := cEn.R().
		SetBody(reqEn).
		SetResult(&resEn).
		Post("http://localhost:80/api/event/name/update")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEn.StatusCode())
	assert.Equal(t, eventNameUpdated, resEn.Event.Name)

	cCu := testhelper.MakeClient(respMl1)

	reqCu := candidate_types.CreateRequest{
		EventID:   resEc.Event.ID,
		AccountID: accountId,
		OpenAt:    "202212271200",
	}
	resCu := candidate_types.CreateResponse{}
	respCu, err := cCu.R().
		SetBody(reqCu).
		SetResult(&resCu).
		Post("http://localhost:80/api/candidate/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCu.StatusCode())
}

func TestAttend001(t *testing.T) {
	respMl, w, err := Login(t, DEFAULT_EMAIL)
	accountId := fmt.Sprintf("%s", w.ID)
	assert.NoError(t, err)

	cEc := testhelper.MakeClient(respMl)
	reqEc := event_types.CreateRequest{
		Name:        "Attendテスト用イベント名",
		Description: "/Users/tsrkzy/dev/go/github.com/tsrkzy/jump_in/go/handler/event_handler/event_test.go Attend 001",
		AccountID:   accountId,
	}
	resEc := event_types.CreateResponse{}
	respEc, err := cEc.R().
		SetBody(reqEc).
		SetResult(&resEc).
		Post("http://localhost:80/api/event/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEc.StatusCode())

	eventId := fmt.Sprintf("%s", resEc.ID)

	/* 参加 */
	cEa1 := testhelper.MakeClient(respMl)
	reqAtt1 := event_types.AttendRequest{EventID: eventId, AccountID: accountId, Comment: "1回目"}
	respAtt1, err := cEa1.R().
		SetBody(reqAtt1).
		//SetResult(reqAtt1).
		Post("http://localhost:80/api/event/attend")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAtt1.StatusCode())

	/* 二重参加 */
	cEa2 := testhelper.MakeClient(respMl)
	reqAtt2 := event_types.AttendRequest{EventID: eventId, AccountID: accountId, Comment: "2回目"}
	respAtt2, err := cEa2.R().
		SetBody(reqAtt2).
		//SetResult(reqAtt2).
		Post("http://localhost:80/api/event/attend")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAtt2.StatusCode())

	/* 参加取消 */
	cEl1 := testhelper.MakeClient(respMl)
	reqLe1 := event_types.LeaveRequest{EventID: eventId, AccountID: accountId}
	respLe1, err := cEl1.R().
		SetBody(reqLe1).
		//SetResult(reqLe1).
		Post("http://localhost:80/api/event/leave")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLe1.StatusCode())

	/* 参加二重取消 */
	cEl2 := testhelper.MakeClient(respMl)
	reqLe2 := event_types.LeaveRequest{EventID: eventId, AccountID: accountId}
	respLe2, err := cEl2.R().
		SetBody(reqLe2).
		//SetResult(reqLe2).
		Post("http://localhost:80/api/event/leave")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLe2.StatusCode())

}

func Login(t *testing.T, email string) (*resty.Response, authenticate_types.WhoAmIResponse, error) {
	redirectUri := "http://localhost:80/api/status"
	r := authenticate_types.Request{
		MailAddress: email,
		RedirectURI: redirectUri,
	}
	/* ML取得、hrefへリダイレクトさせる設定 */
	cA := resty.New().SetDebug(testDebug)
	authResult := authenticate_types.Result{}
	respA, err := cA.R().
		SetBody(r).
		SetResult(&authResult).
		Post("http://localhost:80/api/authenticate")
	assert.NoError(t, err)

	/* NoRedirectPolicyをセットして/mlからのレスポンス検証 */
	cMl1 := resty.New().
		SetHTTPMode().
		SetDebug(testDebug)
	cMl1.SetCookies(respA.Cookies())
	respMl1, err := cMl1.R().
		Get(authResult.MagicLink)
	assert.NoError(t, err)

	cW := resty.New().SetDebug(testDebug)
	cW.SetCookies(respMl1.Cookies())
	cR := authenticate_types.WhoAmIResponse{}
	respW, err := cW.R().
		SetResult(&cR).
		Get("http://localhost:80/api/whoami")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respW.StatusCode())

	return respMl1, cR, err
}
