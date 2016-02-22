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

package glance

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-glance/openstack"
)

type GlanceV2Suite struct {
	suite.Suite
	V1, V2             string
	Images             string
	Img1Size, Img2Size int
	Token              string
}

func (s *GlanceV2Suite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerImages(s, 1000, 2000)
}

func (suite *GlanceV2Suite) TearDownSuite() {
	th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	cinderTestSuite := new(GlanceV2Suite)
	suite.Run(t, cinderTestSuite)
}

func (s *GlanceV2Suite) TestGetImages() {
	Convey("Given Glance images are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetImages called", func() {
				dispatch := ServiceV2{}
				imgs, err := dispatch.GetImages(provider)

				Convey("Then proper image values are returned", func() {
					So(imgs.Count, ShouldEqual, 2)
					So(imgs.Bytes, ShouldEqual, s.Img1Size+s.Img2Size)
				})

				Convey("and no error reported", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
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

func registerAuthentication(s *GlanceV2Suite) {
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

func registerImages(s *GlanceV2Suite, size1 int, size2 int) {
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
