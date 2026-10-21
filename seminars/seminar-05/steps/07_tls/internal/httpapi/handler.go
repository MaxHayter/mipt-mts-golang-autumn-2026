// Package httpapi reuses the tested HTTP handler from step 06.
// TLS changes only how the server is started, not its routes.
package httpapi

import base "http-seminar/internal/httpapi"

type Handler = base.Handler

var New = base.New
