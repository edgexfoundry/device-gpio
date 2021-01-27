# GPIO Device Service
## Overview
GPIO Micro Service - device service for connecting GPIO devices to EdgeX.

- Function:
  - This device service use sysfs to control GPIO devices. For a connected GPIO device, export the needed pin number, set correct GPIO direction, and then start to read or write data from GPIO device.
- Physical interface: system gpio (/sys/class/gpio)
- Driver protocol: IO



## Usage
- This Device Service runs with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command.
- The gpio device service can contains many pre-defined devices which were defined by `configuration.toml`, such as `gpio65`, `gpio66`. These devices are created by the gpio device service in core metadata when the service first initializes. 
- Each device causes the generation of one value of the same type.  For example,  `gpio65` value: `value` .
- After the  gpio device service has started,  we can read or write the system gpio corresponding to pre-defined device by `GET` or ` PUT` the `value` command.

```yaml
deviceResources:
  -
    name: "value"
    description: "set or get system gpio value"
    properties:
      value:
        { type: "Int8", readWrite: "RW", minimum: "0", maximum: "1", defaultValue: "0" }
      units:
        { type: "String", readWrite: "R", defaultValue: "high:1; low:0" }
```




## Guidance
Here we give two step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX.

Before we actually operate GPIO devices, we need to find out RESTful urls of this pre-defined device. By using

`curl http://localhost:48082/api/v1/device/name/gpio65`

Use the `curl` response to get the command URLs (with device and command ids) to issue commands to the GPIO device via the command service as shown below. You can use a tool like `Postman` instead of `curl` to issue the same commands.

```json
{
    "id": "08c22497-4493-4422-bc46-06f3ba7eb804",
    "name": "gpio65",
    "adminState": "UNLOCKED",
    "operatingState": "ENABLED",
    "labels": [
        "gpio65"
    ],
    "commands": [
        {
            "created": 1614825879590,
            "modified": 1614825879590,
            "id": "ead832aa-43b2-4140-b3f2-78b627855997",
            "name": "value",
            "get": {
                "path": "/api/v1/device/{deviceId}/value",
                "responses": [
                    {
                        "code": "200",
                        "expectedValues": [
                            "value"
                        ]
                    },
                    {
                        "code": "500",
                        "description": "service unavailable"
                    }
                ],
                "url": "http://edgex-core-command:48082/api/v1/device/08c22497-4493-4422-bc46-06f3ba7eb804/command/ead832aa-43b2-4140-b3f2-78b627855997"
            },
            "put": {
                "path": "/api/v1/device/{deviceId}/value",
                "responses": [
                    {
                        "code": "200"
                    },
                    {
                        "code": "500",
                        "description": "service unavailable"
                    }
                ],
                "url": "http://edgex-core-command:48082/api/v1/device/08c22497-4493-4422-bc46-06f3ba7eb804/command/ead832aa-43b2-4140-b3f2-78b627855997",
                "parameterNames": [
                    "value"
                ]
            }
        }
    ]
}
```



### Write value to GPIO
Assume we have a GPIO device connected to pin 65 on current system. When we write a value to GPIO, this gpio will be exported and set direction to output.

```shell
curl -H "Content-Type: application/json" -X PUT -d '{"value":"1"}' http://localhost:48082/api/v1/device/08c22497-4493-4422-bc46-06f3ba7eb804/command/ead832aa-43b2-4140-b3f2-78b627855997
```

Now if you test pin 65, it is outputing high voltage.


### Read value from GPIO
Assume we have another GPIO device connected to pin 134 on current system. When we read a value from GPIO, this gpio will be exported and set direction to input.

```shell
curl http://localhost:48082/api/v1/device/08c22497-4493-4422-bc46-06f3ba7eb804/command/ead832aa-43b2-4140-b3f2-78b627855997
```

The command id `ead832aa-43b2-4140-b3f2-78b627855997` here is for the value command.

Here, we post some results:

```bash
$ curl http://localhost...78b627855997`
{"device":"gpio","origin":1611752289806843150,"readings":[{"origin":1611752289806307945,"device":"gpio","name":"value","value":"{\"gpio\":65,\"value\":0}","valueType":"String"}],"EncodedEvent":null}

$ curl http://localhost...78b627855997`
{"device":"gpio","origin":1611752309686651113,"readings":[{"origin":1611752309686212741,"device":"gpio","name":"value","value":"{\"gpio\":65,\"value\":1}","valueType":"String"}],"EncodedEvent":null}
```


## API Reference

| Method | Core Command | parameters        | Description                                                  | Response                    |
| ------ | ------------ | ----------------- | ------------------------------------------------------------ | --------------------------- |
| put    | value        | {"value":<value>} | Set value for the exported gpio<br/><value>: string, "1" or "0" | 200 ok                      |
| get    | value        |                   | Get value of the specified gpio                              | "{\"gpio\":65,\"value\":1}" |



## License
[Apache-2.0](LICENSE)

