package testhelper

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsrkzy/jump_in/types/authenticate_types"
	"gopkg.in/resty.v1"
	"net/http"
	"testing"
)

const testDebug = false

// MakeClient
// ML認証を済ませたレスポンスからcookieを取り出し、
// restyのclientに設定して返す
func MakeClient(MLHttpResponse *resty.Response) *resty.Client {
	clientCandidateDelete := resty.New().SetDebug(testDebug)
	clientCandidateDelete.SetCookies(MLHttpResponse.Cookies())
	return clientCandidateDelete
}

func Login(t *testing.T, email string) (*resty.Response, authenticate_types.WhoAmIResponse, error) {
	redirectUri := "http://localhost:80/api/whoami"
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
