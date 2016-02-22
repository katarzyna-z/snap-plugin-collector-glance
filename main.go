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

package main

import (
    openstackintel "github.com/intelsdi-x/snap-plugin-collector-glance/openstack"
    "github.com/intelsdi-x/snap-plugin-collector-glance/openstack/v2/glance"
    "fmt"
)

const (
    endpoint = "http://192.168.20.2:5000"
    user = "admin"
    password = "admin"
    tenant = "admin"
)

func main(){
    cmn := openstackintel.Common{}
    provider, err := openstackintel.Authenticate(endpoint, user, password, tenant)
    if err != nil {
        panic(err)
    }

    apiv, err := cmn.GetApiVersions(provider)
    if err != nil {
        panic(err)
    }

    chosen, err := openstackintel.ChooseVersion(apiv)
    if err != nil {
        panic(err)
    }
    fmt.Println(chosen)

    srv := glance.ServiceV2{}
    img, err := srv.GetImages(provider)

    fmt.Println(img)
}
