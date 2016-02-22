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
	"net/http"
	"strings"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	"github.com/intelsdi-x/snap-plugin-collector-glance/types"
	str "github.com/intelsdi-x/snap-plugin-utilities/strings"
	"github.com/rackspace/gophercloud"
)

type CollectorSuite struct {
	suite.Suite
	Token              string
	V1, V2             string
	Tenant1, Tenant2   string
	Images             string
	Img1Size, Img2Size int
}

func (s *CollectorSuite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerTenants(s, "demo", "admin")
	registerImages(s, 1000, 2000)
}

func (s *CollectorSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func (s *CollectorSuite) TestGetMetricTypes() {
	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "me", "secret")

		Convey("When GetMetricTypes() is called", func() {
			collector := New()
			mts, err := collector.GetMetricTypes(cfg)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				metricNames := []string{}
				for _, m := range mts {
					metricNames = append(metricNames, strings.Join(m.Namespace(), "/"))
				}

				So(len(mts), ShouldEqual, 4)
				So(str.Contains(metricNames, "intel/openstack/glance/demo/images/Count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/glance/demo/images/Bytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/glance/admin/images/Count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/glance/admin/images/Bytes"), ShouldBeTrue)
			})
		})
	})
}

//func (s *CollectorSuite) TestCollectMetrics() {
//	Convey("Given set of metric types", s.T(), func() {
//		cfg := setupCfg(th.Endpoint(), "me", "secret")
//		m1 := plugin.PluginMetricType{
//			Namespace_: []string{"intel", "openstack", "cinder", "demo", "limits", "MaxTotalVolumeGigabytes"},
//			Config_: cfg.ConfigDataNode}
//		//m2 := plugin.PluginMetricType{
//		//	Namespace_: []string{"intel", "openstack", "cinder", "demo", "volumes", "count"},
//		//	Config_: &cfg.ConfigDataNode}
//		//m3 := plugin.PluginMetricType{
//		//	Namespace_: []string{"intel", "openstack", "cinder", "demo", "snapshots", "bytes"},
//		//	Config_: &cfg.ConfigDataNode}
//		//
//
//
//		servMock := ServicesMock{}
//		limits := types.Limits{
//			MaxTotalVolumeGigabytes: 333,
//			MaxTotalVolumes: 111,
//		}
//		servMock.On("GetLimits", mock.AnythingOfType("*gophercloud.ProviderClient")).Return(limits, nil)
//
//
//		cmnMock := CommonMock{}
//		tenants := []types.Tenant{types.Tenant{ID: "1fffff", Name: "demo"}}
//		cmnMock.On("GetTenants", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(tenants, nil)
//		cmnMock.On("GetApiVersions", mock.AnythingOfType("*gophercloud.ProviderClient")).Return([]string{"v1.0", "v2.0"}, nil)
//
//		Convey("When ColelctMetrics() is called", func() {
//			collector := New()
//
//			collector.common = cmnMock
//			collector.service = servMock
//
//			mts, err := collector.CollectMetrics([]plugin.PluginMetricType{m1})
//
//			Convey("Then no error should be reported", func() {
//				So(err, ShouldBeNil)
//			})
//
//			Convey("and proper metric types are returned", func() {
//				metricNames := map[string]interface{}{}
//				for _, m := range mts {
//					ns := strings.Join(m.Namespace(), "/")
//					metricNames[ns] = m.Data()
//				}
//
//				So(len(mts), ShouldEqual, 1)
//
//			})
//		})
//	})
//}

func TestCollectorSuite(t *testing.T) {
	collectorTestSuite := new(CollectorSuite)
	suite.Run(t, collectorTestSuite)
}

type ServicesMock struct {
	mock.Mock
}

func (servMock ServicesMock) GetImages(provider *gophercloud.ProviderClient) (types.Images, error) {
	ret := servMock.Mock.Called(provider)
	return ret.Get(0).(types.Images), ret.Error(1)
}

type CommonMock struct {
	mock.Mock
}

func (cmnMock CommonMock) GetTenants(endpoint, user, password string) ([]types.Tenant, error) {
	ret := cmnMock.Mock.Called(endpoint, user, password)
	return ret.Get(0).([]types.Tenant), ret.Error(1)
}

func (cmnMock CommonMock) GetApiVersions(provider *gophercloud.ProviderClient) ([]types.ApiVersion, error) {
	ret := cmnMock.Mock.Called(provider)
	return ret.Get(0).([]types.ApiVersion), ret.Error(1)
}

func setupCfg(endpoint, user, password string) plugin.PluginConfigType {
	node := cdata.NewNode()
	node.AddItem("endpoint", ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem("user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("password", ctypes.ConfigValueStr{Value: password})
	return plugin.PluginConfigType{ConfigDataNode: node}
}

func registerRoot() {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"versions": {
						"values": [
							{
								"status": "experimental",
								"id": "v3.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							},
							{
								"status": "stable",
								"id": "v2.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							}
						]
					}
				}
				`, th.Endpoint()+"v3/", th.Endpoint()+"v2.0/")
	})
}

func registerAuthentication(s *CollectorSuite) {
	s.V1 = "v1"
	s.V2 = "v2"
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	th.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"access": {
						"metadata": {
							"is_admin": 0,
							"roles": [
								"3083d61996d648ca88d6ff420542f324"
							]
						},
						"serviceCatalog": [
							{
								"endpoints": [
									{
										"adminURL": "%s",
										"id": "3ffe125aa59547029ed774c10b932349",
										"internalURL": "%s",
										"publicURL": "%s",
										"region": "RegionOne"
									}
								],
								"endpoints_links": [],
								"name": "glance",
								"type": "image"
							}
						],
						"token": {
							"expires": "2016-02-21T14:28:30Z",
							"id": "%s",
							"issued_at": "2016-02-21T13:28:30.656527",
							"tenant": {
								"description": null,
								"enabled": true,
								"id": "97ea299c37bb4e04b3779039ea8aba44",
								"name": "tenant"
							}
						}
					}
				}
			`,
			th.Endpoint(),
			th.Endpoint(),
			th.Endpoint(),
			s.Token)
	})
}

func registerTenants(s *CollectorSuite, tenant1 string, tenant2 string) {
	s.Tenant1 = tenant1
	s.Tenant2 = tenant2
	th.Mux.HandleFunc("/v2.0/tenants", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"tenants": [
					{
						"description": "Test tenat",
						"enabled": true,
						"id": "432534sdfasda",
						"name": "%s"
					},
					{
						"description": "admin tenant",
						"enabled": true,
						"id": "45asdas32",
						"name": "%s"
					}
				],
				"tenants_links": []
			}
		`, s.Tenant1, s.Tenant2)
	})
}

func registerImages(s *CollectorSuite, size1 int, size2 int) {
	s.Images = "/" + s.V2 + "/images"
	s.Img1Size = size1
	s.Img2Size = size2

	th.Mux.HandleFunc(s.Images, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
				{
					"first": "/v2/images",
					"images": [
						{
							"checksum": "eb9139e4942121f22bbc2afc0400b2a4",
							"container_format": "ami",
							"created_at": "2016-02-22T19:06:13Z",
							"disk_format": "ami",
							"file": "/v2/images/5ead7530-3293-40d2-a0ca-f441a33a99e4/file",
							"id": "5ead7530-3293-40d2-a0ca-f441a33a99e4",
							"kernel_id": "e0f483ec-713f-4768-ba1a-220a16b97287",
							"min_disk": 0,
							"min_ram": 0,
							"name": "cirros-0.3.4-x86_64-uec",
							"owner": "ded341b6891c4524b202f08f8808986f",
							"protected": false,
							"ramdisk_id": "95e4ad60-adaf-469d-9711-6baec2ab8a53",
							"schema": "/v2/schemas/image",
							"self": "/v2/images/5ead7530-3293-40d2-a0ca-f441a33a99e4",
							"size": %d,
							"status": "active",
							"tags": [],
							"updated_at": "2016-02-22T19:06:13Z",
							"virtual_size": null,
							"visibility": "public"
						},
						{
							"checksum": "8a40c862b5735975d82605c1dd395796",
							"container_format": "aki",
							"created_at": "2016-02-22T19:06:12Z",
							"disk_format": "aki",
							"file": "/v2/images/e0f483ec-713f-4768-ba1a-220a16b97287/file",
							"id": "e0f483ec-713f-4768-ba1a-220a16b97287",
							"min_disk": 0,
							"min_ram": 0,
							"name": "cirros-0.3.4-x86_64-uec-kernel",
							"owner": "ded341b6891c4524b202f08f8808986f",
							"protected": false,
							"schema": "/v2/schemas/image",
							"self": "/v2/images/e0f483ec-713f-4768-ba1a-220a16b97287",
							"size": %d,
							"status": "active",
							"tags": [],
							"updated_at": "2016-02-22T19:06:12Z",
							"virtual_size": null,
							"visibility": "public"
						}
					],
					"schema": "/v2/schemas/images"
				}
			`, s.Img1Size, s.Img2Size)
	})

}
