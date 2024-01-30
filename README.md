# hyperifyio/goeventd

Simple microservice to trigger SystemD services from NatsIO events.

```
make
./goeventd --nats=nats://localhost:4222 --subject=update-nginx --service=ansible-nginx.service
```

### LICENSE

First 2 years for each release of the software the license is restricted 
for non-commercial use cases. After the restriction period it transitions 
to standard MIT license. See the full license at [./LICENSE.md](LICENSE.md). 

Permission to use the software commercially may be granted under a separate 
commercial agreement.

### USAGE

```
$ ./goeventd --help
Usage of ./goeventd:
  -config string
        Path to configuration file
  -nats string
        The NATS server URL (default "nats://localhost:4222")
  -once
        Shutdown the service after one successful SystemD service execution
  -service string
        The SystemD service to trigger
  -subject string
        The NATS subject to subscribe to
```
