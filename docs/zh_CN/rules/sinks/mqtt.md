# MQTT动作

该操作用于将输出消息发布到 MQTT 服务器中。

| 属性名称 | 是否可选 | 说明                                          |
| ------------- | -------- | ---------------------------------------------------- |
| server        | 否    | mqtt  服务器的代理地址，例如 `tcp：//127.0.0.1：1883` |
| topic          | 否    | mqtt 主题，例如``分析/结果``                         |
| clientId      | 是     | mqtt 连接的客户端 ID。 如果未指定，将使用一个 uuid   |
| protocolVersion   | 是    | 3.1 (也被称为 MQTT 3) 或者 3.1.1 (也被称为 MQTT 4)。 如果未指定，缺省值为 3.1。|
| qos               | 是    | 消息转发的服务质量                               |
| username          | 是    | 连接用户名                            |
| password          | 是    | 连接密码                             |
| certificationPath | 是    | 证书路径。可以为绝对路径，也可以为相对路径。如果指定的是相对路径，那么父目录为执行 `server` 命令的路径。比如，如果你在 `/var/kuiper` 中运行 `bin/server` ，那么父目录为 `/var/kuiper`; 如果运行从 `/var/kuiper/bin` 中运行`./server`，那么父目录为 `/var/kuiper/bin`。 |
| privateKeyPath    | 是    | 私钥路径。可以为绝对路径，也可以为相对路径。更详细的信息，请参考 `certificationPath`. |

以下为使用 SAS 连接到 Azure IoT Hub 的样例。
```json
    {
      "mqtt": {
        "server": "ssl://xyz.azure-devices.net:8883",
        "topic": "devices/demo_001/messages/events/",
        "protocolVersion": "3.1.1",
        "qos": 1,
        "clientId": "demo_001",
        "username": "xyz.azure-devices.net/demo_001/?api-version=2018-06-30",
        "password": "SharedAccessSignature sr=*******************"
      }
    }
```

以下为使用证书和私钥连接到 AWS IoT的另一个样例。

```json
    {
      "mqtt": {
        "server": "ssl://xyz-ats.iot.us-east-1.amazonaws.com:8883",
        "topic": "devices/result",
        "qos": 1,
        "clientId": "demo_001",
        "certificationPath": "keys/d3807d9fa5-certificate.pem",
        "privateKeyPath": "keys/d3807d9fa5-private.pem.key"
      }
    }
```

