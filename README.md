# apimonitor - checks http response codes from services

apimonitor is a utility to check http responses code, it takes as input a json (stored on consul k/v) where we can describe:

 * services to check
 * expected http return codes
 * define result output messages
 
 Its fully integrated with consul using golang/api library, and expects to read its configuration from consul k/v store.

### command line args & options

when launching the application you can set some arguments in order to change some of the settings:

* **-port** (defaults: "8080") its possible to bind application on a different port by setting this flag.

* **-consul** (defaults: "http://localhost:8500") its possible to set a different consul agent by setting this flag.

* **-branch** (defaults: "master") its possible to set a different path in order to fetch application configuration from consul k/v by setting this flag.


### usage & Configuration

Full example:

```
apimonitor -port=9090 -consul=consulhost:8500 -branch=develop 
```

configuration example:

```

{
    "services": [       
        {
            "name": "myService",
            "type": "service",
            "return_code": "200",
            "return_msg": "myserviceOK",
            "url": "https://myservice.myurl.com/"
        }                    
    ]
}

```

Based on the user agent detection the app will adapt its output, in the following example, doing a curl on a cli enviorment will return the follwoing:

```
user@laptop:~$ curl http://localhost:8080

Service Check Status
 myserviceURL | 200 | myserviceOK  | https://myservice.myurl.com/
 anothesrvURL | 404 | youserviceOK | https://anothesrvURL.myurl.com/

 
```


## Building


* apimonitor can be built from source by firstly cloning the repository `git clone https://github.com/numiralofe/apimonitor.git`. Once cloned the binary can be built using:

```
go build apimonitor.go
```


## Contributing
Contributions to apimonitor are welcome! Please bare that this is a first attempt to solve a problem using golang, as so, a few rocky errors are expected, and improvements are expected :)
