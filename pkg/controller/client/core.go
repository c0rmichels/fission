/*
Copyright 2016 The Fission Authors.

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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fission/fission/pkg/info"
)

func (c *Client) SecretGet(m *metav1.ObjectMeta) (*apiv1.Secret, error) {
	relativeUrl := fmt.Sprintf("secrets/%v", m.Name)
	relativeUrl += fmt.Sprintf("?namespace=%v", m.Namespace)

	resp, err := c.get(relativeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var secret apiv1.Secret
	err = json.Unmarshal(body, &secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func (c *Client) ConfigMapGet(m *metav1.ObjectMeta) (*apiv1.ConfigMap, error) {
	relativeUrl := fmt.Sprintf("configmaps/%v", m.Name)
	relativeUrl += fmt.Sprintf("?namespace=%v", m.Namespace)

	resp, err := c.get(relativeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var configMap apiv1.ConfigMap
	err = json.Unmarshal(body, &configMap)
	if err != nil {
		return nil, err
	}

	return &configMap, nil
}

func (c *Client) GetSvcURL(label string) (string, error) {
	resp, err := c.proxy(http.MethodGet, "svcname?"+label, nil)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.Errorf("failed to find service for given label: %v", label)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	storageSvc := string(body)

	return storageSvc, err
}

func (c *Client) ServerInfo() (*info.ServerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := ctxhttp.Get(ctx, &http.Client{}, c.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &info.ServerInfo{}
	err = json.Unmarshal(body, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
