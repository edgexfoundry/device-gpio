# GPIO Device Service
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-gpio/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-gpio/job/main/) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/device-gpio)](https://goreportcard.com/report/github.com/edgexfoundry/device-gpio) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-gpio?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/device-gpio/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-gpio?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/device-gpio)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/device-gpio) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/device-gpio)](https://github.com/edgexfoundry/device-gpio/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/device-gpio)](https://github.com/edgexfoundry/device-gpio/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/device-gpio-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/device-gpio)](https://github.com/edgexfoundry/device-gpio/commits)

## Overview
GPIO Micro Service - device service for connecting GPIO devices to EdgeX

- Function:
  - This device service uses sysfs ABI (default) or chardev ABI (experiment) to control GPIO devices
  For a connected GPIO device, update the configuration files, and then start to read or write data from GPIO device
- This device service **ONLY works on Linux system**
- This device service is contributed by [Jiangxing Intelligence](https://www.jiangxingai.com)


## Usage
- This Device Service runs with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command
- The gpio device service can contains many pre-defined devices which were defined by `res/devices/device.custom.gpio.toml` such as `GPIO-Device01`. These devices are created by the GPIO device service in core metadata when the service first initializes
- Device profiles(`/res/profiles/device.custom.gpio.yaml`) are used to describe the actual GPIO hardware of a device and allow individual gpios to be given human-readable names/aliases
- After the gpio device service has started, we can read or write these corresponding pre-defined devices

```yaml
name: "Custom-GPIO-Device"
manufacturer: "Jiangxing Intelligence"
model: "SP-01"
labels:
  - "device-custom-gpio"
description: "Example of custom gpio device"

deviceResources:
  -
    name: "Power"
    isHidden: false
    description: "mocking power button"
    attributes: { line: 17 }
    properties:
      valueType: "Bool"
      readWrite: "RW"

  -
    name: "LED"
    isHidden: false
    description: "mocking LED"
    attributes: { line: 27 }
    properties:
      valueType: "Bool"
      readWrite: "W"

  -
    name: "Switch"
    isHidden: false
    description: "mocking switch"
    attributes: { line: 22 }
    properties:
      valueType: "Bool"
      readWrite: "R"
```

- Since GPIO sysfs interface is **deprecated after Linux version 4.8**, we provide two ABI interfaces: the sysfs version and the new chardev version. By default we set interface to sysfs, and you can change it inside `[DeviceList.Protocols.interface]` section of `configuration.toml`. For the chardev interface, you still need to specify a selected chip, this is also under `[DeviceList.Protocols.interface]` section.

## Guidance
Here we give two step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX (please notice that, you still need to use Core Command service rather than directly interact with GPIO device service).

Since the `edgex-cli` has released, we can use this new approach to operate devices:

`edgex-cli command list -d GPIO-Device01`

If you would prefer the traditional RESTful way to operate, you can try:

`curl http://localhost:59882/api/v2/device/name/GPIO-Device01`

Use the `curl` response to get the command URLs (with device and command ids) to issue commands to the GPIO device via the command service as shown below. You can also use a tool like `Postman` instead of `curl` to issue the same commands.

```json
{
    "apiVersion": "v2",
    "statusCode": 200,
    "deviceCoreCommand": {
        "deviceName": "GPIO-Device01",
        "profileName": "Custom-GPIO-Device",
        "coreCommands": [
            {
                "name": "Power",
                "get": true,
                "set": true,
                "path": "/api/v2/device/name/GPIO-Device01/Power",
                "url": "http://edgex-core-command:59882",
                "parameters": [
                    {
                        "resourceName": "Power",
                        "valueType": "Bool"
                    }
                ]
            },
            {
                "name": "LED",
                "set": true,
                "path": "/api/v2/device/name/GPIO-Device01/LED",
                "url": "http://edgex-core-command:59882",
                "parameters": [
                    {
                        "resourceName": "LED",
                        "valueType": "Bool"
                    }
                ]
            },
            {
                "name": "Switch",
                "get": true,
                "path": "/api/v2/device/name/GPIO-Device01/Switch",
                "url": "http://edgex-core-command:59882",
                "parameters": [
                    {
                        "resourceName": "Switch",
                        "valueType": "Bool"
                    }
                ]
            }
        ]
    }
}
```



### Write value to GPIO
Assume we have a GPIO device (used for power enable) connected to gpio17 on current system of raspberry pi 4b. When we write a value to GPIO, this gpio will give a high voltage.

```shell
# Set the 'Power' gpio to high
$ curl -X PUT -d   '{"Power":"true"}' http://localhost:59882/api/v2/device/name/GPIO-Device01/Power
{"apiVersion":"v2","statusCode":200}
$ cat /sys/class/gpio/gpio17/direction ; cat /sys/class/gpio/gpio17/value
out
1

# Set the 'Power' gpio to low
$ curl -X PUT -d   '{"Power":"false"}' http://localhost:59882/api/v2/device/name/GPIO-Device01/Power
{"apiVersion":"v2","statusCode":200}
$ cat /sys/class/gpio/gpio17/direction ; cat /sys/class/gpio/gpio17/value
out
0
```

Now if you test gpio17 of raspberry pi 4b , it is outputting high voltage.


### Read value from GPIO
Assume we have another GPIO device (used for button detection) connected to pin 66 on current system. When we read a value from GPIO, this gpio will be exported and set direction to input.

```shell
$ curl http://localhost:59882/api/v2/device/name/GPIO-Device01/Power
```

Here, we post some results:

```bash
{
    "apiVersion": "v2",
    "statusCode": 200,
    "event": {
        "apiVersion": "v2",
        "id": "66e3916f-bac2-4dc6-a53f-befc09a0b888",
        "deviceName": "GPIO-Device01",
        "profileName": "Custom-GPIO-Device",
        "sourceName": "Power",
        "origin": 1631010353930524856,
        "readings": [
            {
                "id": "53e13da0-e2a4-42a0-8a68-54d0cbacbc12",
                "origin": 1631010353930524856,
                "deviceName": "GPIO-Device01",
                "resourceName": "Power",
                "profileName": "Custom-GPIO-Device",
                "valueType": "Bool",
                "binaryValue": null,
                "mediaType": "",
                "value": "false"
            }
        ]
    }
}
```



### docker-compose.yml 

Add the `device-gpio` to the docker-compose.yml of edgex foundry 2.0-Ireland.

```yml
...
	device-gpio:
        container_name: edgex-device-gpio
        depends_on:
        - consul
        - data
        - metadata
        environment:
          CLIENTS_CORE_COMMAND_HOST: edgex-core-command
          CLIENTS_CORE_DATA_HOST: edgex-core-data
          CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
          CLIENTS_SUPPORT_NOTIFICATIONS_HOST: edgex-support-notifications
          CLIENTS_SUPPORT_SCHEDULER_HOST: edgex-support-scheduler
          DATABASES_PRIMARY_HOST: edgex-redis
          EDGEX_SECURITY_SECRET_STORE: "false"
          MESSAGEQUEUE_HOST: edgex-redis
          REGISTRY_HOST: edgex-core-consul
          SERVICE_HOST: edgex-device-gpio
        hostname: edgex-device-gpio
        image: edgexfoundry/device-gpio:0.0.0-dev
        networks:
          edgex-network: {}
        ports:
        - 49994:49994/tcp
        read_only: false
        privileged: true
        volumes:
        - "/sys:/sys"
        - "/dev:/dev"
        security_opt:
        - no-new-privileges:false
        user: root:root
...
```



## API Reference

- `device-custom-gpio-profile.yaml`

  | Core Command | Method | parameters        | Description                                                  | Response |
  | ------------ | ------ | ----------------- | ------------------------------------------------------------ | -------- |
  | Power        | put    | {"Power":<value>} | Set value for the specified gpio<br/><value>: bool, "true" or "false" | 200 ok   |
  |              | get    |                   | Get value of the specified gpio<br/>valueType: Bool          | "true"   |



## License
[Apache-2.0](LICENSE)

