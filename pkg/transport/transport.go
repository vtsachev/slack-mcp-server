package transport

import "net/http"

type UserAgentTransport struct {
	roundTripper http.RoundTripper
	userAgent    string
	cookie       string
	dsCookie     string
}

func New(roundTripper http.RoundTripper, userAgent string, cookie string, dsCookie string) *UserAgentTransport {
	return &UserAgentTransport{
		roundTripper: roundTripper,
		userAgent:    userAgent,
		cookie:       cookie,
		dsCookie:     dsCookie,
	}
}

// RoundTrip implements the RoundTripper interface.
func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Set("User-Agent", t.userAgent)
	clonedReq.Header.Set("Cookie", "d="+t.cookie+";d-s="+t.dsCookie)

	return t.roundTripper.RoundTrip(clonedReq)
}
