# http-mirror-pipeline

http-mirror-pipeline is a tool to mirror HTTP request for continuous testing and benchmarking, replayable logging...

Its modular architecture makes it easy to build powerful pipelines that fit specific needs.

## Example pipelines

## Modules

### Source

#### source.haproxy_spoe

Uses HAProxy SPOE to receive request data. See test/haproxy for an example HAProxy configuration.
Example:

```json
{
  "type": "source.haproxy_spoe",
  "config": {
    "listen_addr": "127.0.0.1:9999"
  }
}
```

| Param         | Value                                    |
| ------------- | ---------------------------------------- |
| `listen_addr` | Can be `<ip>:<port>` or `@<socket_file>` |

### Sinks

#### sink.http

Send the incomming requests to the specified host.

Example:

```json
{
  "type": "sink.http",
  "config": {
    "timeout": "1s",
    "target_url": "http://127.0.0.1:8002"
  }
}
```

| Param        | Value                                                       |
| ------------ | ----------------------------------------------------------- |
| `target_url` | URL to send the requests to. The path in the URL is ignored |
| `timeout`    | Requests timeout. Ex: `1s`, `200ms`, `1m30s`                |

#### sink.file

Writes the requests in a file

Example:

```json
{
  "type": "sink.file",
  "config": {
    "path": "/tmp/my_file",
    "format": "json"
  }
}
```

| Param    | Value                                     |
| -------- | ----------------------------------------- |
| `path`   | Path of the file                          |
| `format` | How to encode requests. `json` or `proto` |

### Control

#### control.fanout

Duplicates messages to all its children. This allows to both write requests to a file and send them in http for example.

Example:

```json
{
  "type": "control.fanout",
  "config": [
    {
      "type": "sink.file",
      "config": { ... }
    },
    {
      "type": "sink.http",
      "config": { ... }
    }
  ]
}
```
