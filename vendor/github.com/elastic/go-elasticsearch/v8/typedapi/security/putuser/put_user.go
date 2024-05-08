// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/4316fc1aa18bb04678b156f23b22c9d3f996f9c9

// Adds and updates users in the native realm. These users are commonly referred
// to as native users.
package putuser

import (
	gobytes "bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"fmt"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
)

const (
	usernameMask = iota + 1
)

// ErrBuildPath is returned in case of missing parameters within the build of the request.
var ErrBuildPath = errors.New("cannot build path, check for missing path parameters")

type PutUser struct {
	transport elastictransport.Interface

	headers http.Header
	values  url.Values
	path    url.URL

	buf *gobytes.Buffer

	req *Request
	raw json.RawMessage

	paramSet int

	username string
}

// NewPutUser type alias for index.
type NewPutUser func(username string) *PutUser

// NewPutUserFunc returns a new instance of PutUser with the provided transport.
// Used in the index of the library this allows to retrieve every apis in once place.
func NewPutUserFunc(tp elastictransport.Interface) NewPutUser {
	return func(username string) *PutUser {
		n := New(tp)

		n.Username(username)

		return n
	}
}

// Adds and updates users in the native realm. These users are commonly referred
// to as native users.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-put-user.html
func New(tp elastictransport.Interface) *PutUser {
	r := &PutUser{
		transport: tp,
		values:    make(url.Values),
		headers:   make(http.Header),
		buf:       gobytes.NewBuffer(nil),
	}

	return r
}

// Raw takes a json payload as input which is then passed to the http.Request
// If specified Raw takes precedence on Request method.
func (r *PutUser) Raw(raw json.RawMessage) *PutUser {
	r.raw = raw

	return r
}

// Request allows to set the request property with the appropriate payload.
func (r *PutUser) Request(req *Request) *PutUser {
	r.req = req

	return r
}

// HttpRequest returns the http.Request object built from the
// given parameters.
func (r *PutUser) HttpRequest(ctx context.Context) (*http.Request, error) {
	var path strings.Builder
	var method string
	var req *http.Request

	var err error

	if r.raw != nil {
		r.buf.Write(r.raw)
	} else if r.req != nil {
		data, err := json.Marshal(r.req)

		if err != nil {
			return nil, fmt.Errorf("could not serialise request for PutUser: %w", err)
		}

		r.buf.Write(data)
	}

	r.path.Scheme = "http"

	switch {
	case r.paramSet == usernameMask:
		path.WriteString("/")
		path.WriteString("_security")
		path.WriteString("/")
		path.WriteString("user")
		path.WriteString("/")
		path.WriteString(url.PathEscape(r.username))

		method = http.MethodPut
	}

	r.path.Path = path.String()
	r.path.RawQuery = r.values.Encode()

	if r.path.Path == "" {
		return nil, ErrBuildPath
	}

	if ctx != nil {
		req, err = http.NewRequestWithContext(ctx, method, r.path.String(), r.buf)
	} else {
		req, err = http.NewRequest(method, r.path.String(), r.buf)
	}

	if r.buf.Len() > 0 {
		req.Header.Set("content-type", "application/vnd.elasticsearch+json;compatible-with=8")
	}

	req.Header.Set("accept", "application/vnd.elasticsearch+json;compatible-with=8")

	if err != nil {
		return req, fmt.Errorf("could not build http.Request: %w", err)
	}

	return req, nil
}

// Do runs the http.Request through the provided transport.
func (r PutUser) Do(ctx context.Context) (*http.Response, error) {
	req, err := r.HttpRequest(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.transport.Perform(req)
	if err != nil {
		return nil, fmt.Errorf("an error happened during the PutUser query execution: %w", err)
	}

	return res, nil
}

// Header set a key, value pair in the PutUser headers map.
func (r *PutUser) Header(key, value string) *PutUser {
	r.headers.Set(key, value)

	return r
}

// Username The username of the User
// API Name: username
func (r *PutUser) Username(v string) *PutUser {
	r.paramSet |= usernameMask
	r.username = v

	return r
}

// Refresh If `true` (the default) then refresh the affected shards to make this
// operation visible to search, if `wait_for` then wait for a refresh to make
// this operation visible to search, if `false` then do nothing with refreshes.
// API name: refresh
func (r *PutUser) Refresh(enum refresh.Refresh) *PutUser {
	r.values.Set("refresh", enum.String())

	return r
}
