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
	Checksum        string            `json:"checksum" mapstructure:"checksum"`
	ContainerFormat string            `json:"container_format" mapstructure:"container_format"`
	CreatedAt       string            `json:"created_at" mapstructure:"created_at"`
	Deleted         bool              `json:"deleted" mapstructure:"deleted"`
	DeletedAt       string            `json:"deleted_at" mapstructure:"deleted_at"`
	DiskFormat      string            `json:"disk_format" mapstructure:"disk_format"`
	ID              string            `json:"id" mapstructure:"id"`
	IsPublic        bool              `json:"is_public" mapstructure:"is_public"`
	MinDisk         int               `json:"min_disk" mapstructure:"min_disk"`
	MinRam          int               `json:"min_ram" mapstructure:"min_ram"`
	Name            string            `json:"name" mapstructure:"name"`
	Owner           string            `json:"owner" mapstructure:"owner"`
	Properties      map[string]string `json:"properties" mapstructure:"properties"`
	Size            int               `json:"size" mapstructure:"size"`
	Status          string            `json:"status" mapstructure:"status"`
	UpdatedAt       string            `json:"updated_at" mapstructure:"updated_at"`
	VirtualSize     string            `json:"virtual_size" mapstructure:"virtual_size"`
}

// GetResult represents the result of a get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r GetResult) Extract() ([]Image, error) {

	var resp struct {
		Images []Image `json:"images" mapstructure:"images"`
	}

	err := mapstructure.Decode(r.Body, &resp)

	return resp.Images, err
}
