// +build small

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

type GlanceV1Suite struct {
	suite.Suite
	V1, V2             string
	Images             string
	Img1Size, Img2Size int
	Token              string
}

func (s *GlanceV1Suite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerImages(s, 1000, 2000)
}

func (suite *GlanceV1Suite) TearDownSuite() {
	th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	cinderTestSuite := new(GlanceV1Suite)
	suite.Run(t, cinderTestSuite)
}

func (s *GlanceV1Suite) TestGetImages() {
	Convey("Given Glance images are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant", "", "")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetImages called", func() {
				dispatch := ServiceV1{}
				imgs, err := dispatch.GetImages(provider)

				Convey("Then proper image values are returned", func() {
					public := imgs["public"]
					So(public.Count, ShouldEqual, 1)
					So(public.Bytes, ShouldEqual, s.Img2Size)

					private := imgs["private"]
					So(private.Count, ShouldEqual, 1)
					So(private.Bytes, ShouldEqual, s.Img1Size)
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

func registerAuthentication(s *GlanceV1Suite) {
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

func registerImages(s *GlanceV1Suite, size1 int, size2 int) {
	s.Images = "/" + s.V1 + "/images/detail"
	s.Img1Size = size1
	s.Img2Size = size2

	th.Mux.HandleFunc(s.Images, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
				{
					"images": [
						{
							"checksum": "19e5f96b987929ef0d56759c2eedf611",
							"container_format": "bare",
							"created_at": "2016-02-25T10:46:13.000000",
							"deleted": false,
							"deleted_at": null,
							"disk_format": "raw",
							"id": "31bbc179-5a75-4d52-98ea-f4f5f6c76279",
							"is_public": false,
							"min_disk": 10,
							"min_ram": 4,
							"name": "AdminVM",
							"owner": "d98e06adf5db49ad9f372625cad7840b",
							"properties": {
								"description": "Private VM for admin"
							},
							"protected": false,
							"size": %d,
							"status": "active",
							"updated_at": "2016-02-25T10:47:15.000000",
							"virtual_size": null
						},
						{
							"checksum": "ee1eca47dc88f4879d8a229cc70a07c6",
							"container_format": "bare",
							"created_at": "2016-02-05T16:04:01.000000",
							"deleted": false,
							"deleted_at": null,
							"disk_format": "qcow2",
							"id": "e256d524-bbd7-40af-9bfa-463d86917459",
							"is_public": true,
							"min_disk": 0,
							"min_ram": 64,
							"name": "TestVM",
							"owner": "76cd5afce159466b885a4731c06998cb",
							"properties": {},
							"protected": false,
							"size": %d,
							"status": "active",
							"updated_at": "2016-02-05T16:04:02.000000",
							"virtual_size": null
						}
					],
					"schema": "/v1/schemas/images"
				}
			`, s.Img1Size, s.Img2Size)
	})

}
