package main

import "net/http"

// Middleware Pattern in Go:
//
// Middleware is a function that wraps an http.Handler, adding extra behavior
// before and/or after the handler executes. This is the "decorator pattern."
//
// Common middleware use cases:
//   - Logging request/response
//   - Authentication/authorization
//   - Rate limiting
//   - Caching headers
//   - Compression
//
// The chain looks like:
//   Client → cacheControl → FileServer → Static Files
//              (adds headers)   (serves files)
//
// The http.Handler interface:
//   type Handler interface {
//       ServeHTTP(http.ResponseWriter, *http.Request)
//   }
//
// Any function can be converted to a Handler using http.HandlerFunc.

// cacheControl adds Cache-Control headers to responses.
// This is a "decorator" - it takes a Handler and returns a new Handler.
//
// The returned HandlerFunc closes over 'next', allowing it to call
// next.ServeHTTP(w, r) after (or before) adding headers.
func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add cache headers BEFORE calling the wrapped handler.
		// The response will include these headers when sent to the client.
		//
		// "public" = can be cached by shared caches (CDNs)
		// "max-age=31536000" = cache for 1 year (in seconds)
		// This is appropriate for static assets (JS, CSS, images).
		w.Header().Set("cache-control", "public,max-age=31536000")

		// Pass control to the next handler in the chain.
		// This is crucial - without it, the request would stop here.
		next.ServeHTTP(w, r)
	})
}
