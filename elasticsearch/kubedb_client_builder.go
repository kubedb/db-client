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

package elasticsearch

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	esv6 "github.com/elastic/go-elasticsearch/v6"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type KubeDBClientBuilder struct {
	kubeClient    kubernetes.Interface
	dynamicClient dynamic.Interface
	db            *api.Elasticsearch
	url           string
	podName       string
}

func NewKubeDBClientBuilder(kubeClient kubernetes.Interface, dClient dynamic.Interface, db *api.Elasticsearch) *KubeDBClientBuilder {
	return &KubeDBClientBuilder{
		kubeClient:    kubeClient,
		dynamicClient: dClient,
		db:            db,
	}
}

func (o *KubeDBClientBuilder) WithPod(podName string) *KubeDBClientBuilder {
	o.podName = podName
	return o
}

func (o *KubeDBClientBuilder) WithURL(url string) *KubeDBClientBuilder {
	o.url = url
	return o
}

func (o *KubeDBClientBuilder) GetElasticClient() (*Client, error) {
	if o.podName != "" {
		o.url = o.getURL()
	}
	var username, password string
	if !o.db.Spec.DisableSecurity && o.db.Spec.AuthSecret != nil {
		secret, err := o.kubeClient.CoreV1().Secrets(o.db.Namespace).Get(context.TODO(), o.db.Spec.AuthSecret.Name, metav1.GetOptions{})
		if err != nil {
			klog.Errorf("Failed to get secret: %s for Elasticsearch: %s/%s with: %s", o.db.Spec.AuthSecret.Name, o.db.Namespace, o.db.Name, err.Error())
			return nil, errors.Wrap(err, "failed to get the secret")
		}

		if value, ok := secret.Data[core.BasicAuthUsernameKey]; ok {
			username = string(value)
		} else {
			klog.Errorf("Failed for secret: %s/%s, username is missing", secret.Namespace, secret.Name)
			return nil, errors.New("username is missing")
		}

		if value, ok := secret.Data[core.BasicAuthPasswordKey]; ok {
			password = string(value)
		} else {
			klog.Errorf("Failed for secret: %s/%s, password is missing", secret.Namespace, secret.Name)
			return nil, errors.New("password is missing")
		}
	}

	// get Elasticsearch version from Elasticsearch version objects
	gvr := schema.GroupVersionResource{
		Group:    "catalog.kubedb.com",
		Version:  "v1alpha1",
		Resource: "elasticsearchversions",
	}
	versionObj, err := o.dynamicClient.Resource(gvr).Get(context.Background(), o.db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get elasticsearch version: %v", err)
	}
	version, found, err := unstructured.NestedFieldNoCopy(versionObj.Object, "spec", "version")
	if err != nil || !found {
		return nil, fmt.Errorf("failed to get version field from ElasticsearchVersion object: %v", err)
	}
	esVersion := version.(string)

	switch {
	// for Elasticsearch 6.x.x
	case strings.HasPrefix(esVersion, "6."):
		client, err := esv6.NewClient(esv6.Config{
			Addresses:         []string{o.url},
			Username:          username,
			Password:          password,
			EnableDebugLogger: true,
			DisableRetry:      true,
			Transport: &http.Transport{
				IdleConnTimeout: 3 * time.Second,
				DialContext: (&net.Dialer{
					Timeout: 30 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					MaxVersion:         tls.VersionTLS12,
				},
			},
		})
		if err != nil {
			klog.Errorf("Failed to create HTTP client for Elasticsearch: %s/%s with: %s", o.db.Namespace, o.db.Name, err.Error())
			return nil, err
		}
		// do a manual health check to test client
		res, err := client.Cluster.Health(
			client.Cluster.Health.WithPretty(),
		)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
		}
		return &Client{
			&ESClientV6{client: client},
		}, nil

	// for Elasticsearch 7.x.x
	case strings.HasPrefix(esVersion, "7."):
		client, err := esv7.NewClient(esv7.Config{
			Addresses:         []string{o.url},
			Username:          username,
			Password:          password,
			EnableDebugLogger: true,
			DisableRetry:      true,
			Transport: &http.Transport{
				IdleConnTimeout: 3 * time.Second,
				DialContext: (&net.Dialer{
					Timeout: 30 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					MaxVersion:         tls.VersionTLS12,
				},
			},
		})
		if err != nil {
			klog.Errorf("Failed to create HTTP client for Elasticsearch: %s/%s with: %s", o.db.Namespace, o.db.Name, err.Error())
			return nil, err
		}
		// do a manual health check to test client
		res, err := client.Cluster.Health(
			client.Cluster.Health.WithPretty(),
		)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.IsError() {
			return nil, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
		}
		return &Client{
			&ESClientV7{client: client},
		}, nil
	}

	return nil, fmt.Errorf("unknown database verseion: %s", o.db.Spec.Version)
}

func (o *KubeDBClientBuilder) getURL() string {
	return fmt.Sprintf("%v://%s.%s.%s.svc:%d", o.db.GetConnectionScheme(), o.podName, o.db.ServiceName(), o.db.GetNamespace(), api.ElasticsearchRestPort)
}
