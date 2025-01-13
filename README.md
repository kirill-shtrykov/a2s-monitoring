# A2S Monitoring

A Go application that retrieves information from an Source server using the A2S (Source Engine Query) protocol and outputs the server details in JSON format.

# Config

| Key | Envvar | Default | Description |
| --- | --- | --- | --- |
| address | A2SMON_ADDR | :9112 | The address to listen. |
| server | A2SMON_SERVER | :27015 | The server address to monitoring. |
| debug | A2SMON_DEBUG | false | Enables debug mode. |

Example:
```shell
a2s-monitoring -address :9112 -server :27015 -debug
```
