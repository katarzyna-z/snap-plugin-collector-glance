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

// Glance package contains wrapper functions designed to collect required metrics
// All functions are dependant on OpenStack ImageService API Version 2

package glance

import (
	"github.com/rackspace/gophercloud"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-glance/openstack"
	"github.com/intelsdi-x/snap-plugin-collector-glance/openstack/v1/images"
	"github.com/intelsdi-x/snap-plugin-collector-glance/types"
)

// ServiceV2 serves as dispatcher for Glance API version 2.0
type ServiceV1 struct{}

// GetLimits collects images by sending REST call to glancehost:9292/v1/images
func (s ServiceV1) GetImages(provider *gophercloud.ProviderClient) (types.Images, error) {
	imgsGlance := types.Images{}

	client, err := openstackintel.NewImageService(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return imgsGlance, err
	}

	imgs, err := images.Get(client).Extract()
	if err != nil {
		return imgsGlance, err
	}

	for _, img := range imgs {
		imgsGlance.Count += 1
		imgsGlance.Bytes += img.Size
	}

	return imgsGlance, nil
}