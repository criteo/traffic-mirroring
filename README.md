# http-mirror-pipeline

http-mirror-pipeline is a tool to mirror HTTP request for continuous testing and benchmarking, replayable logging...

Its modular architecture makes it easy to build powerful pipelines that fit specific needs.

## Example pipelines

### Simple

A simple pipeline that takes requests from HAProxy and sends them to a test server.

![simple](https://github.com/criteo/traffic-mirroring/raw/master/docs/simple.png)

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

![kafka](https://github.com/criteo/traffic-mirroring/raw/master/docs/kafka.png)

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
| `mapping`     | See Mapping                              |

##### Mapping

This config key allows to customize the mapping between SPOE arguments and traffic-mirroring, and add metas to the request.
Left side is the SPOE name, right side is traffic-mirroring.
User provided values will be merged with :

```json
{
  "method": "method",
  "ver": "ver",
  "path": "path",
  "headers": "headers",
  "body": "body"
}
```

Meta can be added like so :

```json
{
  "my_spoe_var": "meta.my_spoe_var"
}
```

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
| `parallel`   | How many requests to send in parallel. Default: 1           |

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

| Param         | Value                                                |
| ------------- | ---------------------------------------------------- |
| `path`        | Path of the file                                     |
| `format`      | How to encode requests. `json` or `proto`            |
| `buffer_size` | Buffer n bytes before writing to file. Default: 1024 |

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

#### control.decouple

http-mirror-pipeline is blocking, meaning a slow sink will slow down upstream modules. While this is useful when reading from a database for example, there are situations where it is preferable to drop requests instead if the output is struggling.

In the following example we decouple the sink from the source to avoid slowing down haproxy :

```json
[
  {
    "type": "source.haproxy_spoe",
    "config": {
      "listen_addr": "127.0.0.1:9999"
    }
  },
  {
    "type": "control.decouple",
    "config": {}
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

| Param   | Value                                                 |
| ------- | ----------------------------------------------------- |
| `quiet` | Do not log dropped messages summary. Default: `false` |

#### control.rate_limit

Rate limits the flow of requests to the specified value. Note that this will block upstream modules, to drop requets exceeding, use a [`control.decouple`](#controldecouple).

In the following example we use a `control.rate_limit` in coordination with a `control.decouple` to perform rate limiting without slowing down haproxy :

```json
[
  {
    "type": "source.haproxy_spoe",
    "config": {
      "listen_addr": "127.0.0.1:9999"
    }
  },
  {
    "type": "control.decouple",
    "config": {}
  },
  {
    "type": "control.rate_limit",
    "config": {
      "rps": 300
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

| Param | Value                       |
| ----- | --------------------------- |
| `rps` | Maximum requests per second |

#### control.split_by

Splits the requets into multiple pipelines based on an arbitrary value. This is useful for example to apply rate limiting on a per-host basis.

In the following example we will do just that :

```json
[
  {
    "type": "source.haproxy_spoe",
    "config": {
      "listen_addr": "127.0.0.1:9999"
    }
  },
  {
    "type": "control.split_by",
    "config": {
      "expr": "{req.header('Host')}",
      "pipeline": {
        "type": "control.seq",
        "config": [
          {
            "type": "control.decouple",
            "config": {}
          },
          {
            "type": "control.rate_limit",
            "config": {
              "rps": "{req.meta.rps.int}"
            }
          }
        ]
      }
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

| Param | Value                       |
| ----- | --------------------------- |
| `rps` | Maximum requests per second |
