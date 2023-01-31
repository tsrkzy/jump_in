package admin_handler

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsrkzy/jump_in/helper/testhelper"
	"github.com/tsrkzy/jump_in/types/admin_types"
	"gopkg.in/resty.v1"
	"net/http"
	"testing"
)

const DEFAULT_EMAIL = "tsrmix+echo@gmail.com"

func TestLogin(t *testing.T) {
	respMl1, w, err := testhelper.Login(t, DEFAULT_EMAIL)
	accountId := fmt.Sprintf("%s", w.Account.ID)
	assert.NoError(t, err)

	pass := "fushianasan"
	true_hash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))

	/* 管理者ログイン */
	AdminLogin(t, accountId, true_hash, respMl1)

	/* 二重管理者ログイン */
	AdminLogin(t, accountId, true_hash, respMl1)

	AdminLogout(t, accountId, respMl1)
}

func AdminLogin(t *testing.T, accountId string, true_hash string, respMl1 *resty.Response) {
	clientAdminLogin := testhelper.MakeClient(respMl1)

	reqAdminLogin := admin_types.LoginRequest{
		AccountID: accountId,
		PassHash:  true_hash,
	}
	resAdminLogin := admin_types.LoginResponse{}
	response, err := clientAdminLogin.R().
		SetBody(&reqAdminLogin).
		SetResult(&resAdminLogin).
		Post("http://localhost:80/api/admin/login")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode())
}

func AdminLogout(t *testing.T, accountId string, respMl1 *resty.Response) {
	clientAdminLogout := testhelper.MakeClient(respMl1)

	reqAdminLogout := admin_types.LogoutRequest{
		AccountID: accountId,
	}
	resAdminLogout := admin_types.LogoutResponse{}
	response, err := clientAdminLogout.R().
		SetBody(&reqAdminLogout).
		SetResult(&resAdminLogout).
		Post("http://localhost:80/api/admin/logout")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode())
}
