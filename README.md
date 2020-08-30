# http-mirror-pipeline

http-mirror-pipeline is a tool to mirror HTTP request for continuous testing and benchmarking, replayable logging...

Its modular architecture makes it easy to build powerful pipelines that fit specific needs.

## Example pipelines

### Simple

A simple pipeline that takes requests from HAProxy and sends them to a test server.

![simple](https://github.com/ShimmerGlass/http-mirror-pipeline/raw/master/docs/simple.png)

`config.json` :

```json
[
  {
    "type": "source.haproxy_spoe",
    "config": {
      "listen_addr": "127.0.0.1:9999"
    }
  },
  {
    "type": "sink.http",
    "config": {
      "timeout": "1s",
      "target_url": "http://127.0.0.1:8002"
    }
  }
]
```

### With Kafka

A more complex use case. In this example we use two pipelines:

- The first one takes requests from HAProxy and sends them to a Kafka queue
- The second one reads the message from Kafka, and both writes them to disk for storage and sends them to a test server

![kafka](https://github.com/ShimmerGlass/http-mirror-pipeline/raw/master/docs/kafka.png)

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

#### source.pcap

Capture and decodes http request from a network interface. Requires root privileges.

Example:

```json
{
  "type": "source.pcap",
  "config": {
    "interface": "lo",
    "port": "80"
  }
}
```

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
