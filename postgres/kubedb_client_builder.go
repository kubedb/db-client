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

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-xorm/xorm"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/tools/certholder"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
)

type KubeDBClientBuilder struct {
	kubeClient kubernetes.Interface
	db         *api.Postgres
	url        string
	podName    string
}

func NewKubeDBClientBuilder(kubeClient kubernetes.Interface, dbObj *api.Postgres) *KubeDBClientBuilder {
	return &KubeDBClientBuilder{
		kubeClient: kubeClient,
		db:         dbObj,
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

func (o *KubeDBClientBuilder) GetPostgresXormClient() (*XormClient, error) {
	connector, err := o.getConnectionString()
	if err != nil {
		return nil, err
	}

	engine, err := xorm.NewEngine("postgres", connector)
	if err != nil {
		return nil, fmt.Errorf("failed to generate postgres client using connection string: %v", err)
	}
	_, err = engine.Query("SELECT 1")
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %v", err)
	}
	return &XormClient{engine}, nil
}

func (o *KubeDBClientBuilder) getURL() string {
	return fmt.Sprintf("%s.%s.%s.svc", o.podName, o.db.GoverningServiceName(), o.db.Namespace)
}

func (o *KubeDBClientBuilder) getPostgresAuthCredentials() (string, string, error) {
	if o.db.Spec.AuthSecret == nil {
		return "", "", errors.New("no database secret")
	}
	secret, err := o.kubeClient.CoreV1().Secrets(o.db.Namespace).Get(context.TODO(), o.db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}
	return string(secret.Data[core.BasicAuthUsernameKey]), string(secret.Data[core.BasicAuthPasswordKey]), nil
}

func (o *KubeDBClientBuilder) GetPostgresClient() (*Client, error) {
	connector, err := o.getConnectionString()
	if err != nil {
		return nil, err
	}
	// connect to database
	db, err := sql.Open("postgres", connector)
	if err != nil {
		return nil, err
	}

	// ping to database to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	return &Client{db}, nil
}

func (o *KubeDBClientBuilder) getConnectionString() (string, error) {
	if o.podName != "" {
		o.url = o.getURL()
	}
	dnsName := o.url
	port := 5432

	user, pass, err := o.getPostgresAuthCredentials()
	if err != nil {
		return "", fmt.Errorf("DB basic auth is not found for PostgreSQL %v/%v", o.db.Namespace, o.db.Name)
	}
	cnnstr := ""
	sslMode := o.db.Spec.SSLMode

	//  sslMode == "prefer" and sslMode == "allow"  don't have support for github.com/lib/pq postgres client. as we are using
	// github.com/lib/pq postgres client utils for connecting our server we need to access with  any of require , verify-ca, verify-full or disable.
	// here we have chosen "require" sslmode to connect postgres as a client
	if sslMode == "prefer" || sslMode == "allow" {
		sslMode = "require"
	}
	if o.db.Spec.TLS != nil {
		secretName := o.db.GetCertSecretName(api.PostgresClientCert)

		certSecret, err := o.kubeClient.CoreV1().Secrets(o.db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})

		if err != nil {
			klog.Error(err, "failed to get certificate secret.", secretName)
			return "", err
		}

		certs, _ := certholder.DefaultHolder.ForResource(api.SchemeGroupVersion.WithResource(api.ResourcePluralPostgres), o.db.ObjectMeta)
		paths, err := certs.Save(certSecret)
		if err != nil {
			klog.Error(err, "failed to save certificate")
			return "", err
		}
		if o.db.Spec.ClientAuthMode == api.ClientAuthModeCert {
			cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%d connect_timeout=10 dbname=postgres sslmode=%s sslrootcert=%s sslcert=%s sslkey=%s", user, pass, dnsName, port, sslMode, paths.CACert, paths.Cert, paths.Key)
		} else {
			cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%d connect_timeout=10 dbname=postgres sslmode=%s sslrootcert=%s", user, pass, dnsName, port, sslMode, paths.CACert)
		}
	} else {
		cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%d connect_timeout=10 dbname=postgres sslmode=%s", user, pass, dnsName, port, sslMode)
	}
	return cnnstr, nil
}
