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

// service contains interface and dispatcher methods for Glance API versions

package services

import (
	"github.com/rackspace/gophercloud"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-glance/openstack"
	glancev1 "github.com/intelsdi-x/snap-plugin-collector-glance/openstack/v1/glance"
	glancev2 "github.com/intelsdi-x/snap-plugin-collector-glance/openstack/v2/glance"
	"github.com/intelsdi-x/snap-plugin-collector-glance/types"
)

// Glancer allows usage of different Glance API versions for metric collection
type Glancer interface {
	GetImages(provider *gophercloud.ProviderClient) (map[string]types.Images, error)
}

// Services serves as a API calls dispatcher
type Service struct {
	glancer Glancer
}

// Set allows to set proper API version implementation
func (c *Service) Set(new Glancer) {
	c.glancer = new
}

// GetImages dispatches call to proper API version calls to collect images metrics
func (s Service) GetImages(provider *gophercloud.ProviderClient) (map[string]types.Images, error) {
	return s.glancer.GetImages(provider)
}

// Dispatch redirects to selected Glance API version based on priority
func Dispatch(provider *gophercloud.ProviderClient) Service {
	cmn := openstackintel.Common{}
	versions, err := cmn.GetApiVersions(provider)
	if err != nil {
		panic(err)
	}

	chosen, err := openstackintel.ChooseVersion(versions)
	if err != nil {
		panic(err)
	}

	service := Service{}
	switch chosen {
	case "v1.0", "v1.1":
		service.Set(glancev1.ServiceV1{})
	case "v2.0", "v2.1", "v2.2", "v2.3":
		service.Set(glancev2.ServiceV2{})
	default:
		panic("Could not select dispatcher")
	}

	return service
}
