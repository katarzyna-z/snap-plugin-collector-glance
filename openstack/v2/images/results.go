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

package images

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// Image represents an Glance image
type Image struct {
	Checksum        string              `json:"checksum" mapstructure:"checksum"`
	ContainerFormat string              `json:"container_format" mapstructure:"container_format"`
	CreatedAt       string              `json:"created_at" mapstructure:"created_at"`
	DirectURL       string              `json:"direct_url" mapstructure:"direct_url"`
	DiskFormat      string              `json:"disk_format" mapstructure:"disk_format"`
	File            string              `json:"file" mapstructure:"file"`
	ID              string              `json:"id" mapstructure:"id"`
	MinDisk         int                 `json:"min_disk" mapstructure:"min_disk"`
	MinRam          int                 `json:"min_ram" mapstructure:"min_ram"`
	Name            string              `json:"name" mapstructure:"name"`
	Owner           string              `json:"owner" mapstructure:"owner"`
	Protected       bool                `json:"protected" mapstructure:"protected"`
	Schema          string              `json:"schema" mapstructure:"schema"`
	Self            string              `json:"self" mapstructure:"self"`
	Size            int                 `json:"size" mapstructure:"size"`
	Status          string              `json:"status" mapstructure:"status"`
	Tags            []map[string]string `json:"tags" mapstructure:"tags"`
	UpdatedAt       string              `json:"updated_at" mapstructure:"updated_at"`
	VirtualSize     string              `json:"virtual_size" mapstructure:"virtual_size"`
	Visibility      string              `json:"visibility" mapstructure:"visibility"`
}

// GetResult represents the result of a get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r GetResult) Extract() ([]Image, error) {

	var resp struct {
		First  string  `mapstructure:"first"`
		Images []Image `json:"images" mapstructure:"images"`
		Schema string  `mapstructure:"schema"`
	}

	err := mapstructure.Decode(r.Body, &resp)

	return resp.Images, err
}
