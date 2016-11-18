# Snap plugin collector - glance

Snap plugin for collecting metrics from OpenStack Glance module.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Snap's Global Config](#snaps-global-config)
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
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-glance/releases) page. Download the plugin from the latest release and load it into `snapd` (`/opt/snap/plugins` is the default location for Snap packages).

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-glance
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-glance.git
```

Build the Snap glance plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started).
* Create Global Config, see description in [Snap's Global Config] (#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](#examples).

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

### Snap's Global Config
Global configuration files are described in [Snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "glance" in "collector" section and then specify following options:
- `"tenant"` - name of the tenant, this parameter is optional. It can be provided at later stage, in task manifest configuration section for metrics.

See example Global Config in [examples/cfg] (examples/cfg/cfg.json).

###Task manifest
User need to provide following parameters in configuration for collector:
- `"endpoint"` - URL for OpenStack Identity endpoint aka Keystone (ex. `"http://keystone.public.org:5000"`)
- `"tenant"` - name of the tenant, this is required if not provided in global config
- `"user"` -  user name which has access to tenant
- `"password"` - user password

See example task manifest in [examples/task] (examples/tasks/task.json).

### Examples
Example of running Snap glance collector and writing data to file.

Download an [example Snap global config](examples/cfg/cfg.json) file.
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-glance/master/examples/cfg/cfg.json
```
Ensure to provide your Keystone instance address and credentials.

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap) with provided configuration file:
* command line: `snapd -l 1 -t 0 --config cfg.json&`

Download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-glance/latest/linux/x86_64/snap-plugin-collector-glance
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snapctl plugin load snap-plugin-collector-glance
$ snapctl plugin load snap-plugin-publisher-file
```

See all available metrics:

```
$ snapctl metric list
```

Download an [example task file](examples/tasks/task.json) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-glance/master/examples/tasks/task.json
$ snapctl task create -t task.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See realtime output from `snapctl task watch <task_id>` (CTRL+C to exit)
```
$ snapctl task watch 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

This data is published to a file `/tmp/published_glance.log` per task specification

Stop task:
```
$ snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)