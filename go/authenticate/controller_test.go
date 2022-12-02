package authenticate

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
	"net/http"
	"testing"
)

const TestDebug = false

/* TestAuthenticate000
 * ML取得APIのバリデーションエラー
 */
func TestAuthenticate000(t *testing.T) {
	/* 不正なメールアドレスとURI */
	email := "tsrmix+echogmail.com"
	redirectUri := "http:/localhost:80/api/status"
	r := Request{
		MailAddress: email,
		RedirectURI: redirectUri,
	}

	cA := resty.New().SetDebug(TestDebug)
	authResult := Result{}
	respA, err := cA.R().
		SetBody(r).
		SetResult(&authResult).
		Post("http://localhost:80/api/authenticate")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respA.StatusCode())
}

/* TestAuthenticate001
 * ML取得、そこへ同じデバイスでアクセスした想定
 * NoRedirectPolicy を使用し、 /ml のレスポンス自体を検証
 */
func TestAuthenticate001(t *testing.T) {
	/* logout */
	cLo := resty.New().SetDebug(TestDebug)
	respLo, err := cLo.R().Get("http://localhost:80/api/logout")
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLo.StatusCode())

	/* ログインしていないので401 */
	cLS := resty.New().SetDebug(TestDebug)
	respLS, err := cLS.R().Get("http://localhost:80/api/status")
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, respLS.StatusCode())

	email := "tsrmix+echo@gmail.com"
	redirectUri := "http://localhost:80/api/status"
	r := Request{
		MailAddress: email,
		RedirectURI: redirectUri,
	}
	/* ML取得、hrefへリダイレクトさせる設定 */
	cA := resty.New().SetDebug(TestDebug)
	authResult := Result{}
	respA, err := cA.R().
		SetBody(r).
		SetResult(&authResult).
		Post("http://localhost:80/api/authenticate")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respA.StatusCode())

	/* NoRedirectPolicyをセットして/mlからのレスポンス検証 */
	cMl1 := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetDebug(TestDebug)
	cMl1.SetCookies(respA.Cookies())
	respMl1, err := cMl1.R().
		Get(authResult.MagicLink)
	/* policy により3xxはエラー扱いになる */
	assert.Error(t, err)
	assert.Equal(t, http.StatusFound, respMl1.StatusCode())

	/* 2回目はMLが無効になっているのでNG (404) */
	cMl2 := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetDebug(TestDebug)
	cMl2.SetCookies(respA.Cookies())
	respMl2, err := cMl2.R().
		Get(authResult.MagicLink)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respMl2.StatusCode())
}

/* TestAuthenticate002
 * ML取得、そこへ同じデバイスでアクセスした想定
 * リダイレクト先で認証済みのレスポンスを受け取れるかを干渉
 */
func TestAuthenticate002(t *testing.T) {
	/* logout */
	cLo := resty.New().SetDebug(TestDebug)
	respLo, err := cLo.R().Get("http://localhost:80/api/logout")
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLo.StatusCode())

	email := "tsrmix+echo@gmail.com"
	redirectUri := "http://localhost:80/api/status"
	r := Request{
		MailAddress: email,
		RedirectURI: redirectUri,
	}
	/* ML取得 */
	cA := resty.New().SetDebug(TestDebug)
	authResult := Result{}
	respA, err := cA.R().
		SetBody(r).
		SetResult(&authResult).
		Post("http://localhost:80/api/authenticate")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respA.StatusCode())

	/* /ml からリダイレクト
	 * REST モードだとリダイレクト時に Header をコピーしないため HTTP モードに */
	cMl1 := resty.New().
		SetHTTPMode().
		SetDebug(TestDebug)
	cMl1.SetCookies(respA.Cookies())
	respMl1, err := cMl1.R().
		Get(authResult.MagicLink)
	/* policy により3xxはエラー扱いになる */
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMl1.StatusCode())

}

/* TestAuthenticate003
 * ML取得、session cookie を持たない別デバイスでアクセスした想定
 */
func TestAuthenticate003(t *testing.T) {
	/* NG */
	email := "tsrmix+echo@gmail.com"
	redirectUri := "http://localhost:80/api/status"
	r := Request{
		MailAddress: email,
		RedirectURI: redirectUri,
	}

	cAu := resty.New().SetDebug(TestDebug)
	authResult := Result{}
	respA, err := cAu.R().
		SetBody(r).
		SetResult(&authResult).
		Post("http://localhost:80/api/authenticate")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respA.StatusCode())

	cMl := resty.New().SetDebug(TestDebug)

	/* cookieを空にする = session の coookie を付与せずMLにアクセスする */
	empty := make([]*http.Cookie, 0)
	cMl.SetCookies(empty)
	respMl, err := cMl.R().
		Get(authResult.MagicLink)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respMl.StatusCode())
}
