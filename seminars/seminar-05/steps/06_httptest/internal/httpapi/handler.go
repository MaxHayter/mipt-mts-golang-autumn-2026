// Package httpapi reuses the finished handler from the upload step.
// This step adds a test; production HTTP code does not change.
package httpapi

import base "http-seminar/internal/httpapi"

type Handler = base.Handler

var New = base.New
