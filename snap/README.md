# EdgeX GPIO Device Service Snap
[![snap store badge](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/edgex-device-gpio)

This folder contains snap packaging for the EdgeX GPIO Device Service Snap

The snap currently supports the following architectures: `amd64`, `arm64`, `armhf`

## Installation
The snap is published in the snap store at https://snapcraft.io/edgex-device-gpio.
You can see the current revisions available for your machine's architecture by running the command:

```bash
$ snap info edgex-device-gpio
```

Note that the application requires access to GPIO on the device.
When running in confined environments, the snap only allows that access via the [gpio interface](https://snapcraft.io/docs/gpio-interface).
For more details, refer to the [GPIO Access](GPIO-Access) section below.

The latest stable version of the snap can be installed using:

```bash
$ sudo snap install edgex-device-gpio
```

The latest development version of the snap can be installed using:

```bash
$ sudo snap install edgex-device-gpio --edge
```

## Snap configuration
### GPIO Access
This snap is strictly confined which means that the access to interfaces are subject to various security measures.

On a Linux distribution without snap confinement for GPIO (e.g. Raspberry Pi OS 11), the snap may be able to access the GPIO directly, without any snap interface and manual connections.

On Linux distributions with snap confinement for GPIO such as Ubuntu Core, the GPIO access is possible via the [gpio interface](https://snapcraft.io/docs/gpio-interface), provided by a gadget snap. 
The official [Raspberry Pi Ubuntu Core](https://ubuntu.com/download/raspberry-pi-core) image includes that gadget.
It is NOT possible to use this snap on Linux distributions that have the GPIO confinement but not the interface (e.g. Ubuntu Server 20.04), unless for development purposes.

In development environments, it is possible to install the snap in dev mode (using `--devmode` flag which disables security confinement and automatic upgrades) and allows direct GPIO access.

The `gpio` interface provides slots for each GPIO channel. The slots can be listed using:
```bash
$ sudo snap interface gpio
name:    gpio
summary: allows access to specific GPIO pin
plugs:
  - edgex-device-gpio
slots:
  - pi:bcm-gpio-0
  - pi:bcm-gpio-1
  - pi:bcm-gpio-10
  ...
```

The slots are not connected automatically. For example, to connect GPIO-17:
```
$ sudo snap connect edgex-device-gpio:gpio pi:bcm-gpio-17
```

Check the list of connections:
```
$ sudo snap connections
Interface        Plug                            Slot              Notes
gpio             edgex-device-gpio:gpio          pi:bcm-gpio-17    manual
…
```

### Startup
Device services implement a service dependency check on startup which ensures that all of the runtime dependencies of a particular service are met before the service transitions to active state.

Snapd doesn't support orchestration between services in different snaps. It is therefore possible on a reboot for a device service to come up faster than all of the required services running in the main edgexfoundry snap. If this happens, it's possible that the device service repeatedly fails startup, and if it exceeds the systemd default limits, then it might be left in a failed state. This situation might be more likely on constrained hardware (e.g. RPi).

This snap therefore implements a basic retry loop with a maximum duration and sleep interval. If the dependent services are not available, the service sleeps for the defined interval (default: 1s) and then tries again up to a maximum duration (default: 60s). These values can be overridden with the following commands:
    
To change the maximum duration, use the following command:

```bash
$ sudo snap set edgex-device-gpio startup-duration=60
```

To change the interval between retries, use the following command:

```bash
$ sudo snap set edgex-device-gpio startup-interval=1
```

The service can then be started as follows. The "--enable" option
ensures that as well as starting the service now, it will be automatically started on boot:

```bash
$ sudo snap start --enable edgex-device-gpio
```

### Vault token
When running this snap and the edgexfoundry snap on the same machine, the vault token can be provisioned via the `edgex-secretstore-token` content interface. 
For details, please refer [here](https://github.com/edgexfoundry/edgex-go/tree/jakarta/snap#interfaces).

### Using a content interface to set device configuration

The `device-config` content interface allows another snap to seed this snap with configuration directories under `$SNAP_DATA/config/device-gpio`.

Note that the `device-config` content interface does NOT support seeding of the Secret Store Token because that file is expected at a different path.

Please refer to [edgex-config-provider](https://github.com/canonical/edgex-config-provider), for an example and further instructions.

### Autostart
By default, the edgex-device-gpio disables its service on install, as the expectation is that the default profile configuration files will be customized, and thus this behavior allows the profile `configuration.toml` files in $SNAP_DATA to be modified before the service is first started.

This behavior can be overridden by setting the `autostart` configuration setting to "true". This is useful when configuration and/or device profiles are being provided via configuration or gadget snap content interface.

**Note** - this option is typically set from a gadget snap.

### Rich Configuration
While it's possible on Ubuntu Core to provide additional profiles via gadget 
snap content interface, quite often only minor changes to existing profiles are required. 

These changes can be accomplished via support for EdgeX environment variable 
configuration overrides via the snap's configure hook.
If the service has already been started, setting one of these overrides currently requires the
service to be restarted via the command-line or snapd's REST API. 
If the overrides are provided via the snap configuration defaults capability of a gadget snap, 
the overrides will be picked up when the services are first started.

The following syntax is used to specify service-specific configuration overrides:


```
env.<stanza>.<config option>
```

For instance, to setup an override of the service's Port use:
```
$ sudo snap set edgex-device-gpio env.service.port=2112
```

And restart the service:
```
$ sudo snap restart edgex-device-gpio.device-gpio
```

**Note** - at this time changes to configuration values in the [Writable] section are not supported.
For details on the mapping of configuration options to Config options, please refer to "Service Environment Configuration Overrides".

### Service Environment Configuration Overrides
**Note** - all of the configuration options below must be specified with the prefix: 'env.'

```
[Service]
service.health-check-interval   // Service.HealthCheckInterval
service.host                    // Service.Host
service.server-bind-addr        // Service.ServerBindAddr
service.port                    // Service.Port
service.max-result-count        // Service.MaxResultCount
service.max-request-size        // Service.MaxRequestSize
service.startup-msg             // Service.StartupMsg
service.request-timeout         // Service.RequestTimeout

[SecretStore]
secret-store.secrets-file               // SecretStore.SecretsFile
secret-store.disable-scrub-secrets-file // SecretStore.DisableScrubSecretsFile

[Clients.core-data]
clients.core-data.port          // Clients.core-data.Port

[Clients.core-metadata]
clients.core-metadata.port      // Clients.core-metadata.Port

[Device]
device.update-last-connected    // Device.UpdateLastConnected
device.use-message-bus          // Device.UseMessageBus
```
