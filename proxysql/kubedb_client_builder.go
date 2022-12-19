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

package proxysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"xorm.io/xorm"
)

type KubeDBClientBuilder struct {
	kc      client.Client
	db      *api.ProxySQL
	url     string
	podName string
}

func NewKubeDBClientBuilder(kc client.Client, db *api.ProxySQL) *KubeDBClientBuilder {
	return &KubeDBClientBuilder{
		kc: kc,
		db: db,
	}
}

func (o *KubeDBClientBuilder) WithURL(url string) *KubeDBClientBuilder {
	o.url = url
	return o
}

func (o *KubeDBClientBuilder) WithPod(podName string) *KubeDBClientBuilder {
	o.podName = podName
	return o
}

func (o *KubeDBClientBuilder) GetProxySQLClient() (*Client, error) {
	connector, err := o.getConnectionString()
	if err != nil {
		return nil, err
	}

	// connect to database
	db, err := sql.Open("mysql", connector)
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

func (o *KubeDBClientBuilder) GetProxySQLXormClient() (*XormClient, error) {
	connector, err := o.getConnectionString()
	if err != nil {
		return nil, err
	}
	engine, err := xorm.NewEngine("mysql", connector)
	if err != nil {
		return nil, err
	}
	_, err = engine.Query("SELECT 1")
	if err != nil {
		return nil, err
	}
	return &XormClient{
		engine,
	}, nil
}

func (o *KubeDBClientBuilder) getURL() string {
	return fmt.Sprintf("%s.%s.%s.svc", o.podName, o.db.GoverningServiceName(), o.db.Namespace)
}

func (o *KubeDBClientBuilder) getProxySQLRootCredentials() (string, string, error) {
	db := o.db
	var secretName string
	if db.Spec.AuthSecret != nil {
		secretName = db.GetAuthSecretName()
	}
	var secret core.Secret
	err := o.kc.Get(context.Background(), client.ObjectKey{Namespace: db.Namespace, Name: secretName}, &secret)
	if err != nil {
		return "", "", err
	}
	user, ok := secret.Data[core.BasicAuthUsernameKey]
	if !ok {
		return "", "", fmt.Errorf("DB root user is not set")
	}
	pass, ok := secret.Data[core.BasicAuthPasswordKey]
	if !ok {
		return "", "", fmt.Errorf("DB root password is not set")
	}
	return string(user), string(pass), nil
}

func (o *KubeDBClientBuilder) getConnectionString() (string, error) {
	user, pass, err := o.getProxySQLRootCredentials()
	if err != nil {
		return "", err
	}

	if o.podName != "" {
		o.url = o.getURL()
	}

	//tlsConfig := ""
	//if o.db.Spec.TLS != nil {
	//	// get client-secret
	//	var clientSecret core.Secret
	//	err := o.kc.Get(context.TODO(), client.ObjectKey{Namespace: o.db.GetNamespace(), Name: o.db.GetCertSecretName(api.ProxySQLClientCert)}, &clientSecret)
	//	if err != nil {
	//		return "", err
	//	}
	//	cacrt := clientSecret.Data["ca.crt"]
	//	certPool := x509.NewCertPool()
	//	certPool.AppendCertsFromPEM(cacrt)
	//
	//	crt := clientSecret.Data["tls.crt"]
	//	key := clientSecret.Data["tls.key"]
	//	cert, err := tls.X509KeyPair(crt, key)
	//	if err != nil {
	//		return "", err
	//	}
	//	var clientCert []tls.Certificate
	//	clientCert = append(clientCert, cert)
	//
	//	// tls custom setup
	//	if o.db.Spec.RequireSSL {
	//		err = sql_driver.RegisterTLSConfig(api.MySQLTLSConfigCustom, &tls.Config{
	//			RootCAs:      certPool,
	//			Certificates: clientCert,
	//		})
	//		if err != nil {
	//			return "", err
	//		}
	//		tlsConfig = fmt.Sprintf("tls=%s", api.MySQLTLSConfigCustom)
	//	} else {
	//		tlsConfig = fmt.Sprintf("tls=%s", api.MySQLTLSConfigSkipVerify)
	//	}
	//}

	connector := fmt.Sprintf("%v:%v@tcp(%s:%d)/", user, pass, o.url, 6032)
	return connector, nil
}
