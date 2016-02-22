/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2015 Intel Corporation
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

package collector

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rackspace/gophercloud"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	str "github.com/intelsdi-x/snap-plugin-utilities/strings"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-glance/openstack"
	"github.com/intelsdi-x/snap-plugin-collector-glance/openstack/services"
	"github.com/intelsdi-x/snap-plugin-collector-glance/types"
)

const (
	name    = "glance"
	version = 1
	plgtype = plugin.CollectorPluginType
	vendor  = "intel"
	fs      = "openstack"
)

// New creates initialized instance of Glance collector
func New() *collector {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	providers := map[string]*gophercloud.ProviderClient{}
	return &collector{host: host, tenants: str.InitSet(), providers: providers}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (c *collector) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}
	items, err := config.GetConfigItems(cfg, []string{"endpoint", "user", "password"})
	if err != nil {
		return nil, err
	}

	endpoint := items["endpoint"].(string)
	user := items["user"].(string)
	password := items["password"].(string)

	// retrieve list of all available tenants for provided endpoint, user and password
	cmn := openstackintel.Common{}
	allTenants, err := cmn.GetTenants(endpoint, user, password)
	if err != nil {
		return nil, err
	}

	// Generate available namespace for images
	namespaces := []string{}
	for _, tenant := range allTenants {
		// Construct temporary struct to generate namespace based on tags
		var metrics struct {
			I types.Images `json:"images"`
		}
		current := strings.Join([]string{vendor, fs, name, tenant.Name}, "/")
		ns.FromCompositionTags(metrics, current, &namespaces)
	}

	for _, namespace := range namespaces {
		mts = append(mts, plugin.PluginMetricType{
			Namespace_: strings.Split(namespace, "/"),
			Config_:    cfg.ConfigDataNode,
		})
	}

	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (c *collector) CollectMetrics(metricTypes []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	// iterate over metric types to resolve needed collection calls
	// for requested tenants
	for _, metricType := range metricTypes {
		namespace := metricType.Namespace()
		if len(namespace) < 6 {
			return nil, fmt.Errorf("Incorrect namespace lenth. Expected 6 is %d", len(namespace))
		}

		tenant := namespace[3]
		c.tenants.Add(tenant)
	}

	allImages := map[string]types.Images{}

	for _, tenant := range c.tenants.Elements() {
		fmt.Println("tenant:", tenant)
		if err := c.authenticate(metricTypes[0], tenant); err != nil {
			return nil, err
		}

		provider := c.providers[tenant]

		imgs, err := c.service.GetImages(provider)
		if err != nil {
			panic(err)
		}
		allImages[tenant] = imgs
	}

	metrics := []plugin.PluginMetricType{}
	for _, metricType := range metricTypes {
		namespace := metricType.Namespace()
		tenant := namespace[3]
		// Construct temporary struct to accommodate all gathered metrics
		metricContainer := struct {
			I types.Images `json:"images"`
		}{
			allImages[tenant],
		}

		// Extract values by namespace from temporary struct and create metrics
		metric := plugin.PluginMetricType{
			Source_:    c.host,
			Timestamp_: time.Now(),
			Namespace_: namespace,
			Data_:      ns.GetValueByNamespace(metricContainer, namespace[4:]),
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (c *collector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	return cp, nil
}

// Commenting exported items is very important
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		plgtype,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
	)
}

type collector struct {
	host    string
	tenants str.StringSet
	service services.Service
	common  openstackintel.Commoner
	providers map[string]*gophercloud.ProviderClient
}

func (c *collector) authenticate(cfg interface{}, tenant string) error {
	if _, found := c.providers[tenant]; !found {
		// get credentials and endpoint from configuration
		items, err := config.GetConfigItems(cfg, []string{"endpoint", "user", "password"})
		if err != nil {
			return err
		}

		endpoint := items["endpoint"].(string)
		user := items["user"].(string)
		password := items["password"].(string)

		provider, err := openstackintel.Authenticate(endpoint, user, password, tenant)
		if err != nil {
			return err
		}
		// set provider and dispatch API version based on priority
		c.providers[tenant] = provider
		c.service = services.Dispatch(provider)

		// set Commoner interface
		c.common = openstackintel.Common{}
	}

	return nil
}
