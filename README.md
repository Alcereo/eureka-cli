Tiny console util to communicate with [Eureka](https://github.com/Netflix/eureka) 
service discovery server started in [Spring Boot](https://projects.spring.io/spring-boot/) app.

Ready-made app can be taken from [releases](https://github.com/Alcereo/eureka-cli/releases)

#### Use cases

##### Get instances info

```bash
> eureka-cli -u $EUREKA_HOST -p $EUREKA_PORT info

APP NAME            STATUS    ID                                 IP ADDRESS        PORT
EUREKA-CLIENT-2     UP        client-2                           192.168.0.114     38891
EUREKA-CLIENT       UP        192.168.0.114:eureka-client:0      192.168.0.114     35935
EUREKA-CLIENT       UP        192.168.0.114:eureka-client        192.168.0.114     8080
```

##### Get instance url

```bash
> eureka-cli -u $EUREKA_HOST -p $EUREKA_PORT info url EUREKA-CLIENT-2 client-2

http://192.168.0.114:38891
```

##### Wait for instance UP status

```bash
> eureka-cli wait EUREKA-CLIENT-2 client-2

Wait for instanceID: "client-2" app name: "EUREKA-CLIENT-2"...
It took:  7.004323836s
APP NAME            STATUS    ID                                 IP ADDRESS        PORT               
EUREKA-CLIENT-2     UP        client-2                           192.168.0.114     38891 
```


##### Wait for instance and request `/info` hook

```bash
> eureka-cli wait -t 12 EUREKA-CLIENT-2 client-2 && wget -qO- $(eureka-cli info url EUREKA-CLIENT-2 client-2)/info

Wait for instanceID: "client-2" app name: "EUREKA-CLIENT-2"...
It took:  7.004323836s
APP NAME            STATUS    ID                                 IP ADDRESS        PORT               
EUREKA-CLIENT-2     UP        client-2                           192.168.0.114     38891 
{}
```

##### Wait for instance and request `/info` hook

```bash
> eureka-cli wait EUREKA-CLIENT-2 client-2 && wget -qO- $(eureka-cli info url EUREKA-CLIENT-2 client-2)/info

Wait for instanceID: "client-2" app name: "EUREKA-CLIENT-2"...
It took:  7.004323836s
APP NAME            STATUS    ID                                 IP ADDRESS        PORT               
EUREKA-CLIENT-2     UP        client-2                           192.168.0.114     38891 
{}
```

##### Quick check instance 
```bash
> eureka-cli wait -t 0 EUREKA-CLIENT-2 client-2

APP NAME            STATUS    ID                                 IP ADDRESS        PORT               
EUREKA-CLIENT-2     UP        client-2                           192.168.0.114     38891

# Process finished with exit code 0

> eureka-cli wait -t 0 EUREKA-CLIENT-2 client-3

Not found
# Process finished with exit code 2
```