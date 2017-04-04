/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
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

// common contains shared functions for general purposes, like Authentication, choosing version etc.

package openstack

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"

	"github.com/intelsdi-x/snap-plugin-collector-glance/apiversions"
	"github.com/intelsdi-x/snap-plugin-collector-glance/types"
)

var apiPriority = map[string]int{
	"v1.0": 0,
	"v1.1": 1,
	"v2.0": 2,
	"v2.1": 3,
	"v2.2": 4,
	"v2.3": 5,
}

// Commoner provides abstraction for shared functions mainly for mocking
type Commoner interface {
	GetTenants(endpoint, user, password string) ([]types.Tenant, error)
	GetApiVersions(provider *gophercloud.ProviderClient) ([]types.ApiVersion, error)
}

// Common is a receiver for Commoner interface
type Common struct{}

// GetTenants is used to retrieve list of available tenant for authenticated user
// List of tenants can then be used to authenticate user for each given tenant
func (c Common) GetTenants(endpoint, user, password string) ([]types.Tenant, error) {
	tnts := []types.Tenant{}

	provider, err := Authenticate(endpoint, user, password, "", "", "")
	if err != nil {
		return nil, err
	}

	client := openstack.NewIdentityV2(provider)

	opts := tenants.ListOpts{}
	pager := tenants.List(client, &opts)

	page, err := pager.AllPages()
	if err != nil {
		return tnts, err
	}

	tenantList, err := tenants.ExtractTenants(page)
	if err != nil {
		return tnts, err
	}

	for _, t := range tenantList {
		tnts = append(tnts, types.Tenant{Name: t.Name, ID: t.ID})
	}

	return tnts, nil
}

// GetApiVersions is used to retrieve list of available Cinder API versions
// List of api version is then used to dispatch calls to proper API version based on defined priority
func (c Common) GetApiVersions(provider *gophercloud.ProviderClient) ([]types.ApiVersion, error) {
	apis := []types.ApiVersion{}

	client, err := NewImageService(provider, gophercloud.EndpointOpts{
		Availability: gophercloud.AvailabilityInternal},
	)

	if err != nil {
		return apis, err
	}

	apiVersions, err := apiversions.Get(client).Extract()
	if err != nil {
		return apis, err
	}

	for _, apiVersion := range apiVersions {
		link := apiVersion.Links[0]
		apis = append(apis, types.ApiVersion{
			ID:   apiVersion.ID,
			Link: link["href"],
		})
	}

	return apis, nil
}

// Authenticate is used to authenticate user for given tenant. Request is send to provided Keystone endpoint
// Returns authenticated provider client, which is used as a base for service clients.
func Authenticate(endpoint, user, password, tenant, domain_name, domain_id string) (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: endpoint,
		Username:         user,
		Password:         password,
		TenantName:       tenant,
		AllowReauth:      true,
	}
	if domain_name != "" && domain_id == "" {
		authOpts.DomainName = domain_name
	}
	if domain_id != "" && domain_name == "" {
		authOpts.DomainID = domain_id
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// ChooseVersion returns chosen Cinder API version based on defined priority
func ChooseVersion(recognized []types.ApiVersion) (string, error) {
	if len(recognized) < 1 {
		return "", fmt.Errorf("No recognized API versions provided")
	}
	chosen := recognized[0].ID
	for _, ver := range recognized[1:] {
		chosenPriority, ok1 := apiPriority[chosen]
		verPriority, ok2 := apiPriority[ver.ID]
		if ok1 && ok2 {
			if chosenPriority < verPriority {
				chosen = ver.ID
			}
		}
	}
	return chosen, nil
}
