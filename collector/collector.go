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

package collector

import (
	"os"
	"strings"
	"time"

	"github.com/rackspace/gophercloud"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"

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
	return &collector{host: host, providers: providers}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (c *collector) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}
	tenant := ""
	item, err := config.GetConfigItem(cfg, "tenant")
	if err != nil {
		tenant = "*"
	} else {
		tenant = item.(string)
	}

	namespaces := []string{}
	// Construct temporary struct to generate namespace based on tags
	var metrics struct {
		I struct {
			Prv types.Images `json:"private"`
			Pub types.Images `json:"public"`
			Sha types.Images `json:"shared"`
		} `json:"images"`
	}

	current := strings.Join([]string{vendor, fs, name, tenant}, "/")
	ns.FromCompositionTags(metrics, current, &namespaces)

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
	//allImages := map[string]types.Images{}

	// get credentials and endpoint from configuration
	items, err := config.GetConfigItems(metricTypes[0], []string{"endpoint", "tenant", "user", "password"})
	if err != nil {
		return nil, err
	}

	endpoint := items["endpoint"].(string)
	tenant := items["tenant"].(string)
	user := items["user"].(string)
	password := items["password"].(string)

	if err := c.authenticate(endpoint, tenant, user, password); err != nil {
		return nil, err
	}

	provider := c.providers[tenant]

	imgs, err := c.service.GetImages(provider)
	if err != nil {
		return nil, err
	}

	metrics := []plugin.PluginMetricType{}
	for _, metricType := range metricTypes {
		namespace := metricType.Namespace()
		// Construct temporary struct to generate namespace based on tags
		metricContainer := struct {
			I struct {
				Prv types.Images `json:"private"`
				Pub types.Images `json:"public"`
				Sha types.Images `json:"shared"`
			} `json:"images"`
		}{
			struct {
				Prv types.Images `json:"private"`
				Pub types.Images `json:"public"`
				Sha types.Images `json:"shared"`
			}{imgs["private"], imgs["public"], imgs["shared"]},
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

// Meta returns plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		plgtype,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.RoutingStrategy(plugin.StickyRouting),
	)
}

type collector struct {
	host      string
	service   services.Service
	common    openstackintel.Commoner
	providers map[string]*gophercloud.ProviderClient
}

func (c *collector) authenticate(endpoint, tenant, user, password string) error {
	if _, found := c.providers[tenant]; !found {
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
