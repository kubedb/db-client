/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package restproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

type Client struct {
	*resty.Client
	Config *Config
}

type Response struct {
	Code   int
	Header http.Header
	Body   io.ReadCloser
}

type ResponseBody struct {
	Status  string `json:"status"`
	Brokers []int  `json:"checks"`
}

type Config struct {
	host             string
	api              string
	connectionScheme string
	transport        *http.Transport
}

func (cc *Client) GetKafkaBrokerList() (*Response, error) {
	req := cc.Client.R().SetDoNotParseResponse(true)
	res, err := req.Get(cc.Config.api)
	if err != nil {
		klog.Error(err, "Failed to send http request")
		return nil, err
	}
	response := &Response{
		Code:   res.StatusCode(),
		Header: res.Header(),
		Body:   res.RawBody(),
	}

	if response.Code != http.StatusOK {
		return response, fmt.Errorf("failed to get broker list from rest proxy")
	}

	return response, nil
}

// IsBrokerAvailableForRequest parse health response in json from server and
// return overall status of the server
func (cc *Client) IsBrokerAvailableForRequest(response *Response) (bool, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err1 := errors.Wrap(err, "failed to parse response body")
			if err1 != nil {
				return
			}
			return
		}
	}(response.Body)

	var responseBody ResponseBody
	body, _ := io.ReadAll(response.Body)
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse response body")
	}

	if len(responseBody.Brokers) == 0 {
		return false, fmt.Errorf("no brokers found to serve request with rest proxy")
	}

	return true, nil
}
