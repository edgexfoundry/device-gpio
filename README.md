# GPIO Device Service
## Overview
GPIO Micro Service - device service for connecting GPIO devices to EdgeX

- Function:
  - This device service uses sysfs ABI (default) or chardev ABI (experiment) to control GPIO devices
  For a connected GPIO device, update the configuration files, and then start to read or write data from GPIO device
- This device service **ONLY works on Linux system**
- This device service is contributed by [Jiangxing Intelligence](https://www.jiangxingai.com)


## Usage
- This Device Service runs with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command
- The gpio device service can contains many pre-defined devices which were defined by `configuration.toml` such as `Custom-GPIO-Device`. These devices are created by the GPIO device service in core metadata when the service first initializes
- Device profiles are used to describe the actual GPIO hardware of a device and allow individual gpios to be given human-readable names/aliases
- After the  gpio device service has started,  we can read or write these  corresponding pre-defined devices

```yaml
deviceResources:
  -
    name: "Power"
    description: "system power gpio export"
    attributes: { line: "134" }
    properties:
      value:
        { type: "Bool", readWrite: "RW" }

  -
    name: "Led"
    description: "system led gpio export"
    attributes: { line: "65" }
    properties:
      value:
        { type: "Bool", readWrite: "W" }

  -
    name: "Button"
    description: "system button gpio export"
    attributes: { line: "66" }
    properties:
      value:
        { type: "Bool", readWrite: "R" }
```

- Since GPIO sysfs interface is deprecated after Linux version 4.8, we provide two ABI interfaces: the sysfs version and the new chardev version. By default we set interface to sysfs, and you can change it inside `[DeviceList.Protocols.interface]` section of `configuration.toml`. For the chardev interface, you still need to specify a selected chip, this is also under `[DeviceList.Protocols.interface]` section.

## Guidance
Here we give two step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX.

Since the `edgex-cli` has released, we can use this new approach to operate devices:

`edgex-cli command list -d Custom-GPIO-Device`

If you would prefer the traditional RESTful way to operate, you can try:

`curl http://localhost:48082/api/v1/device/name/Custom-GPIO-Device`

Use the `curl` response to get the command URLs (with device and command ids) to issue commands to the GPIO device via the command service as shown below. You can also use a tool like `Postman` instead of `curl` to issue the same commands.

```json
{
    "id": "d734883a-0c66-4213-9bfb-864e0ce076cc",
    "name": "Custom-GPIO-Device",
    "adminState": "UNLOCKED",
    "operatingState": "ENABLED",
    "labels": [
        "device-custom-gpio"
    ],
    "commands": [
        ......
        {
            "created": 1615972980751,
            "modified": 1615972980751,
            "id": "cd5a77a5-8d3e-4657-8752-a8d6ccae73b7",
            "name": "Power",
            "get": {
                "path": "/api/v1/device/{deviceId}/Power",
                "responses": [
                    {
                        "code": "200",
                        "expectedValues": [
                            "Power"
                        ]
                    },
                    {
                        "code": "500",
                        "description": "service unavailable"
                    }
                ],
                "url": "http://edgex-core-command:48082/api/v1/device/d734883a-0c66-4213-9bfb-864e0ce076cc/command/cd5a77a5-8d3e-4657-8752-a8d6ccae73b7"
            },
            "put": {
                "path": "/api/v1/device/{deviceId}/Power",
                "responses": [
                    {
                        "code": "200"
                    },
                    {
                        "code": "500",
                        "description": "service unavailable"
                    }
                ],
                "url": "http://edgex-core-command:48082/api/v1/device/d734883a-0c66-4213-9bfb-864e0ce076cc/command/cd5a77a5-8d3e-4657-8752-a8d6ccae73b7",
                "parameterNames": [
                    "Power"
                ]
            }
        },
		......
    ]
}
```



### Write value to GPIO
Assume we have a GPIO device (used for power enable) connected to pin 134 on current system. When we write a value to GPIO, this gpio will give a high voltage.

```shell
edgex-cli command put -d Custom-GPIO-Device -n Power -b '{"Power":"true"}'
```

Now if you test pin 134, it is outputting high voltage.


### Read value from GPIO
Assume we have another GPIO device ( used for button detection ) connected to pin 66 on current system. When we read a value from GPIO, this gpio will be exported and set direction to input.

```shell
curl http://localhost:48082/api/v1/device/d734883a-0c66-4213-9bfb-864e0ce076cc/command/852161f4-5ddf-418d-9202-39682cfb1dca
```

The command id `852161f4-5ddf-418d-9202-39682cfb1dca` here is for the `Button` command.

Here, we post some results:

```bash
$ curl http://localhost:48082......852161f4-5ddf-418d-9202-39682cfb1dca
{"device":"Custom-GPIO-Device","origin":1615974593898774961,"readings":[{"origin":1615974593893644001,"device":"Custom-GPIO-Device","name":"Button","value":"false","valueType":"Bool"}],"EncodedEvent":null}

$ curl http://localhost:48082......852161f4-5ddf-418d-9202-39682cfb1dca
{"device":"Custom-GPIO-Device","origin":1615974593898774961,"readings":[{"origin":1615974593893644001,"device":"Custom-GPIO-Device","name":"Button","value":"true","valueType":"Bool"}],"EncodedEvent":null}
```



## API Reference

- `device-custom-gpio-profile.yaml`

  | Core Command | Method | parameters        | Description                                                  | Response |
  | ------------ | ------ | ----------------- | ------------------------------------------------------------ | -------- |
  | Power        | put    | {"Power":<value>} | Set value for the specified gpio<br/><value>: bool, "true" or "false" | 200 ok   |
  |              | get    |                   | Get value of the specified gpio<br/>valueType: Bool          | "true"   |



## License
[Apache-2.0](LICENSE)

