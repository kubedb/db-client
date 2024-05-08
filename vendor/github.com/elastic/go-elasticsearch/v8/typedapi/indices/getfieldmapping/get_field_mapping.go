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

// Returns mapping for one or more fields.
package getfieldmapping

import (
	gobytes "bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fmt"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
)

const (
	fieldsMask = iota + 1

	indexMask
)

// ErrBuildPath is returned in case of missing parameters within the build of the request.
var ErrBuildPath = errors.New("cannot build path, check for missing path parameters")

type GetFieldMapping struct {
	transport elastictransport.Interface

	headers http.Header
	values  url.Values
	path    url.URL

	buf *gobytes.Buffer

	paramSet int

	fields string
	index  string
}

// NewGetFieldMapping type alias for index.
type NewGetFieldMapping func(fields string) *GetFieldMapping

// NewGetFieldMappingFunc returns a new instance of GetFieldMapping with the provided transport.
// Used in the index of the library this allows to retrieve every apis in once place.
func NewGetFieldMappingFunc(tp elastictransport.Interface) NewGetFieldMapping {
	return func(fields string) *GetFieldMapping {
		n := New(tp)

		n.Fields(fields)

		return n
	}
}

// Returns mapping for one or more fields.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-get-field-mapping.html
func New(tp elastictransport.Interface) *GetFieldMapping {
	r := &GetFieldMapping{
		transport: tp,
		values:    make(url.Values),
		headers:   make(http.Header),
		buf:       gobytes.NewBuffer(nil),
	}

	return r
}

// HttpRequest returns the http.Request object built from the
// given parameters.
func (r *GetFieldMapping) HttpRequest(ctx context.Context) (*http.Request, error) {
	var path strings.Builder
	var method string
	var req *http.Request

	var err error

	r.path.Scheme = "http"

	switch {
	case r.paramSet == fieldsMask:
		path.WriteString("/")
		path.WriteString("_mapping")
		path.WriteString("/")
		path.WriteString("field")
		path.WriteString("/")
		path.WriteString(url.PathEscape(r.fields))

		method = http.MethodGet
	case r.paramSet == indexMask|fieldsMask:
		path.WriteString("/")
		path.WriteString(url.PathEscape(r.index))
		path.WriteString("/")
		path.WriteString("_mapping")
		path.WriteString("/")
		path.WriteString("field")
		path.WriteString("/")
		path.WriteString(url.PathEscape(r.fields))

		method = http.MethodGet
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

	req.Header.Set("accept", "application/vnd.elasticsearch+json;compatible-with=8")

	if err != nil {
		return req, fmt.Errorf("could not build http.Request: %w", err)
	}

	return req, nil
}

// Do runs the http.Request through the provided transport.
func (r GetFieldMapping) Do(ctx context.Context) (*http.Response, error) {
	req, err := r.HttpRequest(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.transport.Perform(req)
	if err != nil {
		return nil, fmt.Errorf("an error happened during the GetFieldMapping query execution: %w", err)
	}

	return res, nil
}

// IsSuccess allows to run a query with a context and retrieve the result as a boolean.
// This only exists for endpoints without a request payload and allows for quick control flow.
func (r GetFieldMapping) IsSuccess(ctx context.Context) (bool, error) {
	res, err := r.Do(ctx)

	if err != nil {
		return false, err
	}
	io.Copy(ioutil.Discard, res.Body)
	err = res.Body.Close()
	if err != nil {
		return false, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return true, nil
	}

	return false, nil
}

// Header set a key, value pair in the GetFieldMapping headers map.
func (r *GetFieldMapping) Header(key, value string) *GetFieldMapping {
	r.headers.Set(key, value)

	return r
}

// Fields A comma-separated list of fields
// API Name: fields
func (r *GetFieldMapping) Fields(v string) *GetFieldMapping {
	r.paramSet |= fieldsMask
	r.fields = v

	return r
}

// Index A comma-separated list of index names
// API Name: index
func (r *GetFieldMapping) Index(v string) *GetFieldMapping {
	r.paramSet |= indexMask
	r.index = v

	return r
}

// AllowNoIndices Whether to ignore if a wildcard indices expression resolves into no concrete
// indices. (This includes `_all` string or when no indices have been specified)
// API name: allow_no_indices
func (r *GetFieldMapping) AllowNoIndices(b bool) *GetFieldMapping {
	r.values.Set("allow_no_indices", strconv.FormatBool(b))

	return r
}

// ExpandWildcards Whether to expand wildcard expression to concrete indices that are open,
// closed or both.
// API name: expand_wildcards
func (r *GetFieldMapping) ExpandWildcards(value string) *GetFieldMapping {
	r.values.Set("expand_wildcards", value)

	return r
}

// IgnoreUnavailable Whether specified concrete indices should be ignored when unavailable
// (missing or closed)
// API name: ignore_unavailable
func (r *GetFieldMapping) IgnoreUnavailable(b bool) *GetFieldMapping {
	r.values.Set("ignore_unavailable", strconv.FormatBool(b))

	return r
}

// IncludeDefaults Whether the default mapping values should be returned as well
// API name: include_defaults
func (r *GetFieldMapping) IncludeDefaults(b bool) *GetFieldMapping {
	r.values.Set("include_defaults", strconv.FormatBool(b))

	return r
}

// Local Return local information, do not retrieve the state from master node
// (default: false)
// API name: local
func (r *GetFieldMapping) Local(b bool) *GetFieldMapping {
	r.values.Set("local", strconv.FormatBool(b))

	return r
}
