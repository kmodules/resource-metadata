/*
Copyright AppsCode Inc. and Contributors.

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

package identity

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"path"

	kmapi "kmodules.xyz/client-go/api/v1"
	identityapi "kmodules.xyz/resource-metadata/apis/identity/v1alpha1"

	"go.bytebuilders.dev/license-verifier/info"
	"gomodules.xyz/sync"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

type Client struct {
	baseURL string
	token   string
	caCert  []byte
	client  *http.Client

	clusterUID string
}

func NewClient(baseURL, token string, caCert []byte, clusterUID string) (*Client, error) {
	c := &Client{
		baseURL: baseURL,
		token:   token,
		caCert:  caCert,
	}
	if len(caCert) == 0 {
		c.client = http.DefaultClient
	} else {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		c.client = &http.Client{Transport: transport}
	}
	return c, nil
}

func (c *Client) Identify(clusterUID string) (*kmapi.ClusterMetadata, error) {
	u, err := info.APIServerAddress(c.baseURL)
	if err != nil {
		return nil, err
	}
	apiEndpoint := u.String()
	u.Path = path.Join(u.Path, "api/v1/clustersv2/identity", clusterUID)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apierrors.NewGenericServerResponse(
			resp.StatusCode,
			http.MethodGet,
			schema.GroupResource{Group: identityapi.GroupName, Resource: identityapi.ResourceClusterIdentities},
			"",
			string(body),
			0,
			false,
		)
	}
	var md kmapi.ClusterMetadata
	err = json.Unmarshal(body, &md)
	if err != nil {
		return nil, err
	}

	md.APIEndpoint = apiEndpoint
	md.CABundle = string(c.caCert)

	return &md, nil
}

func (c *Client) GetToken() (string, error) {
	u, err := info.APIServerAddress(c.baseURL)
	if err != nil {
		return "", err
	}

	id, err := c.GetIdentity()
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "api/v1/agent", id.Status.Name, id.Status.UID, "token")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// add authorization header to the req
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

const SelfName = "self"

var (
	identity          *identityapi.ClusterIdentity
	once              sync.Once
	idError           error
	creationTimestamp = metav1.Now()
)

func (c *Client) GetIdentity() (*identityapi.ClusterIdentity, error) {
	once.Do(func() error {
		var md *kmapi.ClusterMetadata
		md, idError = c.Identify(c.clusterUID)
		if idError != nil {
			return idError
		}
		identity = &identityapi.ClusterIdentity{
			ObjectMeta: metav1.ObjectMeta{
				UID:               types.UID("cid-" + c.clusterUID),
				Name:              SelfName,
				CreationTimestamp: creationTimestamp,
				Generation:        1,
			},
			Status: *md,
		}
		idError = nil
		return idError
	})
	return identity, idError
}
