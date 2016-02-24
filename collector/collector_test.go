// +build unit

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
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	str "github.com/intelsdi-x/snap-plugin-utilities/strings"
)

type CollectorSuite struct {
	suite.Suite
	Token              string
	V1, V2             string
	Tenant1, Tenant2   string
	Images             string
	Img1Size, Img2Size int
	Server             *httptest.Server
}

func (s *CollectorSuite) SetupSuite() {
	// for glance calls
	th.SetupHTTP()

	// for identity calls
	router := mux.NewRouter()
	s.Server = httptest.NewServer(router)

	registerIdentityRoot(s, router)
	registerIdentityTokens(s, router)
	registerIdentityTenants(s, router, "demo", "admin")
	registerGlanceApi(s)
	registerGlanceImages(s, 1000, 2000)
}

func (s *CollectorSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func (s *CollectorSuite) TestGetMetricTypes() {
	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		cfg := setupCfg(s.Server.URL, "me", "secret")

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

func (s *CollectorSuite) TestCollectMetrics() {
	Convey("Given set of metric types", s.T(), func() {
		cfg := setupCfg(s.Server.URL, "me", "secret")
		m1 := plugin.PluginMetricType{
			Namespace_: []string{"intel", "openstack", "glance", "demo", "images", "Count"},
			Config_:    cfg.ConfigDataNode}
		m2 := plugin.PluginMetricType{
			Namespace_: []string{"intel", "openstack", "glance", "demo", "images", "Bytes"},
			Config_:    cfg.ConfigDataNode}

		Convey("When ColelctMetrics() is called", func() {
			collector := New()

			mts, err := collector.CollectMetrics([]plugin.PluginMetricType{m1, m2})

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				metricNames := map[string]interface{}{}
				for _, m := range mts {
					ns := strings.Join(m.Namespace(), "/")
					metricNames[ns] = m.Data()
				}
				fmt.Println(metricNames)
				So(len(mts), ShouldEqual, 2)

				val, ok := metricNames["intel/openstack/glance/demo/images/Count"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, 2)

				val, ok = metricNames["intel/openstack/glance/demo/images/Bytes"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, s.Img1Size+s.Img2Size)
			})
		})
	})
}

func TestCollectorSuite(t *testing.T) {
	collectorTestSuite := new(CollectorSuite)
	suite.Run(t, collectorTestSuite)
}

func setupCfg(endpoint, user, password string) plugin.PluginConfigType {
	node := cdata.NewNode()
	node.AddItem("endpoint", ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem("user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("password", ctypes.ConfigValueStr{Value: password})
	return plugin.PluginConfigType{ConfigDataNode: node}
}

func registerIdentityRoot(s *CollectorSuite, r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
				`, s.Server.URL+"/v3/", s.Server.URL+"/v2.0/")
	})
}

func registerIdentityTokens(s *CollectorSuite, r *mux.Router) {
	s.V1 = "v1"
	s.V2 = "v2"
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	r.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
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

func registerIdentityTenants(s *CollectorSuite, r *mux.Router, tenant1 string, tenant2 string) {
	s.Tenant1 = tenant1
	s.Tenant2 = tenant2
	r.HandleFunc("/v2.0/tenants", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Auth-Token", s.Token)
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
	}).Methods("GET")
}

func registerGlanceApi(s *CollectorSuite) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		//th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
		{
			"versions": [
				{
					"id": "v2.3",
					"links": [
						{
							"href": "%s/v2/",
							"rel": "self"
						}
					],
					"status": "CURRENT"
				},
				{
					"id": "v2.2",
					"links": [
						{
							"href": "%s/v2/",
							"rel": "self"
						}
					],
					"status": "SUPPORTED"
				},
				{
					"id": "v2.1",
					"links": [
						{
							"href": "%s/v2/",
							"rel": "self"
						}
					],
					"status": "SUPPORTED"
				},
				{
					"id": "v2.0",
					"links": [
						{
							"href": "%s/v2/",
							"rel": "self"
						}
					],
					"status": "SUPPORTED"
				},
				{
					"id": "v1.1",
					"links": [
						{
							"href": "%s/v1/",
							"rel": "self"
						}
					],
					"status": "SUPPORTED"
				},
				{
					"id": "v1.0",
					"links": [
						{
							"href": "%s/v1/",
							"rel": "self"
						}
					],
					"status": "SUPPORTED"
				}
			]
		}
		`,
			th.Endpoint(),
			th.Endpoint(),
			th.Endpoint(),
			th.Endpoint(),
			th.Endpoint(),
			th.Endpoint(),
		)
	})
}

func registerGlanceImages(s *CollectorSuite, size1 int, size2 int) {
	s.Images = "/" + s.V2 + "/images"
	s.Img1Size = size1
	s.Img2Size = size2

	th.Mux.HandleFunc(s.Images, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		//th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

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
