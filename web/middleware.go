package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) error

type Middleware func(Handler) Handler

type Chain struct {
	middlewares []Middleware
}

func NewChain(middlewares ...Middleware) Chain {
	return Chain{append(([]Middleware)(nil), middlewares...)}
}

func (m Chain) Then(middleware ...Middleware) Chain {
	return NewChain(append(m.middlewares, middleware...)...)
}

func (m Chain) Finally(h Handler) http.HandlerFunc {
	for i := len(m.middlewares) - 1; i >= 0; i = i - 1 {
		h = m.middlewares[i](h)
	}
	return handleError(h)
}

func Method(method string) Middleware {
	return func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if r.Method == method {
				return h(w, r)
			}
			return ErrMethodNotAllowed
		}
	}
}

func handleError(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			var errorCode int
			switch err := err.(type) {
			case Error:
				errorCode = err.Code
			default:
				errorCode = http.StatusInternalServerError
			}

			w.WriteHeader(errorCode)
			json.NewEncoder(w).Encode(&struct {
				Error string `json:"error"`
				Type  string `json:"type"`
			}{err.Error(), fmt.Sprintf("%T", err)})
		}
	}
}

func MakeTorrentHandler(h func(string, http.ResponseWriter) error) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var torrentInfo struct {
			Hash string `json:"hash"`
		}
		err := json.NewDecoder(r.Body).Decode(&torrentInfo)
		if err != nil {
			return err
		}
		return h(torrentInfo.Hash, w)
	}
}
