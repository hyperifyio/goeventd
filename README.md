# hyperifyio/goeventd

This repository hosts `goeventd`, a microservice designed to facilitate the 
interaction between NatsIO events and system services. It's a simple, yet 
powerful tool for triggering commands or file system changes based on specific 
events in a NatsIO stream.

![main branch status](https://github.com/hyperifyio/goeventd/actions/workflows/build.yml/badge.svg?branch=main)
![dev branch status](https://github.com/hyperifyio/goeventd/actions/workflows/build.yml/badge.svg?branch=dev)

## Getting Started

To use `goeventd`, simply clone the repository, build the binary using `make`, 
and execute it with the necessary flags:

```bash
make
./goeventd --nats=nats://localhost:4222 --subject=update-nginx --service=ansible-nginx.service
```

## License

Copyright (c) Heusala Group Ltd. All rights reserved.

Each software release of `goeventd` is initially under the HG Evaluation and 
Non-Commercial License for the first two years. This allows use, modification, 
and distribution for non-commercial and evaluation purposes only. Post this 
period, the license transitions to the standard MIT license, permitting broader
usage, including commercial applications. For full details, refer to the 
[LICENSE.md](LICENSE.md) file. 

**Commercial usage licenses can be obtained under separate agreements.**

## Usage and Configuration

To understand the usage and available configuration options for `goeventd`, run:

```bash
./goeventd --help
```

This will display help information, including available flags:

```
Usage of ./goeventd:
  -command-template string
        The command to trigger services. The service name can be provided in %s placeholder. (when event-type is command) (default "systemctl start %s")
  -config string
        Path to configuration file
  -event-type string
        The event triggering mechanism (command/file) (default "command")
  -file-template string
        The path to trigger services by file modifications. The service name can be provided in %s placeholder. (when event-type is file) (default "./%s.trigger")
  -nats string
        The NATS server URL (default "nats://127.0.0.1:4222")
  -once
        Shutdown the service after one service triggering
  -service string
        The SystemD service to trigger
  -subject string
        The NATS subject to subscribe to
```

## Deploying as a SystemD Service

### Installation Steps

1. **Download and Install**:

    Download the latest release and install the binary:

    ```bash
    wget https://github.com/hyperifyio/goeventd/releases/download/v0.0.18/goeventd-v0.0.18-linux-amd64.zip
    unzip goeventd-v0.0.18-linux-amd64.zip
    cp goeventd-v0.0.18-linux-amd64 /usr/local/bin/goeventd
    chmod 755 /usr/local/bin/goeventd
    ```

2. **Create a Service File**:

    You can write your own `goeventd.service` file or use the provided template:

    ```ini
    [Unit]
    Description=GoEventD Service
    After=network.target nats.service

    [Service]
    Type=simple
    User=goeventd
    Group=goeventd
    ExecStart=/usr/local/bin/goeventd --subject=ansible.nginx --service=ansible-nginx
    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
    ```

    To use the provided service file:

    ```bash
    sudo cp goeventd.service /etc/systemd/system/
    ```

3. **Reload and Start the Service**:

    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable goeventd
    sudo systemctl start goeventd
    ```

4. **Check the Service Status**:

    ```bash
    sudo systemctl status goeventd
    ```

This service file example assumes that `goeventd` will communicate with a NATS 
server and trigger a specific SystemD service. Modify the `ExecStart` line in 
the service file according to your needs. 

For more detailed instructions, check out the 
[documentation](https://github.com/hyperifyio/goeventd/wiki).
