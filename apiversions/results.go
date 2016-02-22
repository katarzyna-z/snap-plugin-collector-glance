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

package apiversions

import (
	"github.com/rackspace/gophercloud"

	"github.com/mitchellh/mapstructure"
)

// APIVersion represents an API version for Cinder.
type APIVersion struct {
	ID     string              `json:"id" mapstructure:"id"`
	Status string              `json:"status" mapstructure:"status"`
	Links  []map[string]string `json:"links" mapstructure:"links"`
}

// GetResult represents the result of a get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r GetResult) Extract() ([]APIVersion, error) {

	var resp struct {
		Versions []APIVersion `mapstructure:"versions"`
	}

	err := mapstructure.Decode(r.Body, &resp)

	return resp.Versions, err
}