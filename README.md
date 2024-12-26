# A2S Monitoring

A Go application that retrieves information from an Source server using the A2S (Source Engine Query) protocol and outputs the server details in JSON format.

# Config

| Key | Envvar | Default | Description |
| --- | --- | --- | --- |
| address | A2SMON_ADDR | :9112 | The address to listen. |
| status | A2SMON_STATUS | /status | The path to the server status. |
| server | A2SMON_SERVER | :27015 | The server address to monitoring. |

Example:
```shell
a2s-monitoring -address :9112 -status /status -server :27015
```
