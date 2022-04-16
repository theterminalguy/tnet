package util

import (
	"net/url"
	"strings"
)

// IsRedirectURISecureStrict is stricter than IsRedirectURISecure and it does not allow custom-scheme
// URLs because they can be hijacked for native apps. Use claimed HTTPS redirects instead.
// See discussion in https://github.com/ory/fosite/pull/489.
func IsRedirectURISecureStrict(redirectURI *url.URL) bool {
	return redirectURI.Scheme == "https" || (redirectURI.Scheme == "http" && IsLocalhost(redirectURI))
}

func IsRedirectURISecure(redirectURI *url.URL) bool {
	return !(redirectURI.Scheme == "http" && !IsLocalhost(redirectURI))
}

func IsLocalhost(redirectURI *url.URL) bool {
	hn := redirectURI.Hostname()
	return strings.HasSuffix(hn, ".localhost") || hn == "127.0.0.1" || hn == "::1" || hn == "localhost"
}

// IsValidRedirectURI validates a redirect_uri as specified in:
//
// * https://tools.ietf.org/html/rfc6749#section-3.1.2
//   * The redirection endpoint URI MUST be an absolute URI as defined by [RFC3986] Section 4.3.
//   * The endpoint URI MUST NOT include a fragment component.
// * https://tools.ietf.org/html/rfc3986#section-4.3
//   absolute-URI  = scheme ":" hier-part [ "?" query ]
// * https://tools.ietf.org/html/rfc6819#section-5.1.1
func IsValidRedirectURI(redirectURI *url.URL) bool {
	// We need to explicitly check for a scheme
	if !IsRequestURL(redirectURI.String()) {
		return false
	}

	if redirectURI.Fragment != "" {
		// "The endpoint URI MUST NOT include a fragment component."
		return false
	}

	return true
}

// IsRequestURL check if the string rawurl, assuming
// it was received in an HTTP request, is a valid
// URL confirm to RFC 3986
func IsRequestURL(rawurl string) bool {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false //Couldn't even parse the rawurl
	}
	if len(url.Scheme) == 0 {
		return false //No Scheme found
	}
	return true
}
