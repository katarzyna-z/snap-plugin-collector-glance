# snap plugin collector - glance

snap plugin for collecting metrics from OpenStack Glance module. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [snap's Global Config](#snaps-global-config)
  * [Task manifest](#task-manifest)
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
 * OpenStack deployment available
 * Supports Glance v1 and v2 APIs
 
### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download glance plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapd` to access.

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

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-glance/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-glance/blob/master/README.md#examples).

#### Suggestions
* It is not recommended to set interval for task less than 20 seconds. This may lead to overloading Glance API with requests.

## Documentation
### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
intel/openstack/glance/\<tenant_name\>/images/public/count | int | Total number of OpenStack public images for given tenant
intel/openstack/glance/\<tenant_name\>/images/private/count | int | Total number of OpenStack private images for given tenant
intel/openstack/glance/\<tenant_name\>/images/shared/count | int | Total number of OpenStack shared images for given tenant
intel/openstack/glance/\<tenant_name\>/images/public/bytes | int | Total number of bytes used by OpenStack private images for given tenant
intel/openstack/glance/\<tenant_name\>/images/private/bytes | int | Total number of bytes used by OpenStack public images for given tenant
intel/openstack/glance/\<tenant_name\>/images/shared/bytes | int | Total number of bytes used by OpenStack shared images for given tenant

### snap's Global Config
Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "glance" in "collector" section and then specify following options:
- `"tenant"` - name of the tenant, this parameter is optional. It can be provided at later stage, in task manifest configuration section for metrics.

See example Global Config in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-glance/blob/master/examples/cfg/).

###Task manifest
User need to provide following parameters in configuration for collector:
- `"endpoint"` - URL for OpenStack Identity endpoint aka Keystone (ex. `"http://keystone.public.org:5000"`)
- `"tenant"` - name of the tenant, this is required if not provided in global config
- `"user"` -  user name which has access to tenant
- `"password"` - user password

See example task manifest in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-glance/blob/master/examples/tasks/).

### Examples
Example running snap-plugin-collector-glance plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

Create Global Config, see example in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-glance/blob/master/examples/cfg/).

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in cfg.json):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0 --config cfg.json
```
In another terminal window:

Load snap-plugin-collector-glance plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-glance
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
```
See available metrics for your system:
```
$ $SNAP_PATH/bin/snapctl metric list
```
Create a task manifest file to use snap-plugin-collector-glance plugin (exemplary file in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-glance/blob/master/examples/tasks/)):
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
		        "/intel/openstack/glance/*/images/public/count": {},
		        "/intel/openstack/glance/*/images/public/bytes": {}
           },
            "config": {
              "/intel/openstack/glance": {
                "endpoint": "http://keystone.public.org:5000",
                "user": "admin",
                "password": "admin",
                "tenant": "test_tenant"
              }
            },
            "process": null,
            "publish": null
        }
    }
}
```
Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t task.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. The full project is at http://github.com/intelsdi-x/snap.
To reach out on other use cases, visit:
* [snap Gitter channel](https://gitter.im/intelsdi-x/snap)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)