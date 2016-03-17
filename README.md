# snap-plugin-collector-glance

snap plugin for collecting metrics from OpenStack Glance module. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

Plugin collects metrics by communicating with OpenStack by REST API.
It can run locally on the host, or in proxy mode (communicating with the host via HTTP(S)). 

### System Requirements

 - Linux
 - OpenStack deployment available
 - Supports Glance v1 and v2 APIs 

### Installation
#### Download glance plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-glance
Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-glance
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
intel/openstack/glance/\<tenant_name\>/images/Count | int | Total number of OpenStack images for given tenant
intel/openstack/glance/\<tenant_name\>/images/Bytes | int | Total number of bytes used by OpenStack images for given tenant

### snap's Global Config
Global configuration files are described in snap's documentation. You have to add section "glance" in "collector" section and then specify following options:
- `"tenant"` - name of the tenant, this parameter is optional 

### Examples
It is not suggested to set interval below 20 seconds. This may lead to overloading Keystone with authentication requests. 

User need to provide following parameters in configuration for collector
- `"endpoint"` - URL for OpenStack Identity endpoint aka Keystone (ex. `"http://keystone.public.org:5000"`)
- `"tenant"` - name of the tenant, this parameter is optionadl 
- `"user"` -  user name which has access to tenant
- `"password"` - user password 

Example task manifest to use <glance> plugin:
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "60s"
    },
    "workflow": {
        "collect": {
            "metrics": {
		        "/intel/openstack/glance/demo/images/public/count": {},
		        "/intel/openstack/glance/demo/images/public/bytes": {}
           },
            "config": {
                "endpoint": "http://keystone.public.org:5000",
                "user": "admin",
                "password": "admin",
                "tenant": "test_tenant"
            },
            "process": null,
            "publish": null
        }
    }
}
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Marcin Krolik](https://github.com/marcin-krolik)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.