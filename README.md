# GPIO Device Service
## Overview
GPIO Micro Service - device service for connecting GPIO devices to EdgeX.

- Function:
  - This device service use sysfs to control GPIO devices. For a connected GPIO device, export the needed pin number, set correct GPIO direction, and then start to read or write data from GPIO device.
- Physical interface: system gpio (/sys/class/gpio)
- Driver protocol: IO



## Usage
- This Device Service runs with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command.
- The gpio device service can contains many pre-defined devices which were defined by `configuration.toml`, such as `Custom-GPIO-Device`. These devices are created by the gpio device service in core metadata when the service first initializes. 
- Each pre-defined device has its own profile, such as `device-custom-gpio-profile.yaml`, which choose the sysfs gpios to export by `deviceResources`. Then it add `coreCommands`  and `coreCommands` with the  name same to `deviceResources`.
- After the  gpio device service has started,  we can read or write the system gpio corresponding to pre-defined device by specific command.

```yaml
deviceResources:
  -
    name: "Power"
    description: "system power gpio export"
    attributes: { number: "/134" }
    properties:
      value:
        { type: "Bool", readWrite: "RW" }

  -
    name: "Led"
    description: "system led gpio export"
    attributes: { number: "/65" }
    properties:
      value:
        { type: "Bool", readWrite: "W" }

  -
    name: "Button"
    description: "system button gpio export"
    attributes: { number: "/66" }
    properties:
      value:
        { type: "Bool", readWrite: "R" }
```

- Since GPIO sysfs interface is deprecated after Linux version 4.8, we provide two ABI interfaces: the sysfs version and the new chardev version. Please change GOTAG in Makefile, and select the interface you like. For the sysfs version, in device resource, you need to provide attribute number like  "/123", which 123 is the pin number you selected. As for chardev version, the number should be like "0/17", which the first number is for gpiochip, and the second number is for selected line number. 


## Guidance
Here we give two step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX.

Before we actually operate GPIO devices, we need to find out RESTful urls of this pre-defined device. By using

`curl http://localhost:48082/api/v1/device/name/Custom-GPIO-Device`

Use the `curl` response to get the command URLs (with device and command ids) to issue commands to the GPIO device via the command service as shown below. You can use a tool like `Postman` instead of `curl` to issue the same commands.

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
Assume we have a GPIO device ( used for power enable ) connected to pin 134 on current system. When we write a value to GPIO, this gpio will be exported and set direction to output.

```shell
curl -H "Content-Type: application/json" -X PUT -d '{"Power":"true"}' http://localhost:48082/api/v1/device/d734883a-0c66-4213-9bfb-864e0ce076cc/command/cd5a77a5-8d3e-4657-8752-a8d6ccae73b7
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
  | Power        | put    | {"Power":<value>} | Set value for the exported gpio<br/><value>: bool, "true" or "false" | 200 ok   |
  |              | get    |                   | Get value of the specified gpio<br/>valueType: Bool          | "true"   |



## License
[Apache-2.0](LICENSE)

