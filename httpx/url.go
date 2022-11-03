// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net/http"
	"net/url"

	"github.com/ory/x/stringsx"
)

// IncomingRequestURL returns the URL of the incoming HTTP request by looking at the host, TLS, and X-Forwarded-* headers.
func IncomingRequestURL(r *http.Request) *url.URL {
	source := *r.URL
	source.Host = stringsx.Coalesce(source.Host, r.Header.Get("X-Forwarded-Host"), r.Host)

	if proto := r.Header.Get("X-Forwarded-Proto"); len(proto) > 0 {
		source.Scheme = proto
	}

	if source.Scheme == "" {
		source.Scheme = "https"
		if r.TLS == nil {
			source.Scheme = "http"
		}
	}

	return &source
}
