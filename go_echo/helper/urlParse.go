package helper

import (
	"fmt"
	"net/url"
)

// ExtractOriginFromURI
// returns "http(s)://domain:port"
// @SEE https://qiita.com/nyamage/items/80ca5480baad8da581df
func ExtractOriginFromURI(redirectUri string) (string, error) {
	u, err := url.Parse(redirectUri)
	if err != nil {
		return "", err
	}
	scheme := u.Scheme
	host := u.Host

	redirectOrigin := fmt.Sprintf("%s://%s", scheme, host)
	return redirectOrigin, nil
}
