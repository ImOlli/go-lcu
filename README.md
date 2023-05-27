# GO-LCU

This library automatically retrieves the port and auth-token for the LCU(League of Legends Client)

[What's the LCU?](https://riot-api-libraries.readthedocs.io/en/latest/lcu.html)

> **Important!** This library currently only works on windows.

## Install

```shell
go get github.com/ImOlli/go-lcu
```

## Usage

### Get port and auth-token of LeagueClient

```go
package main

import (
	"github.com/ImOlli/go-lcu/lcu"
	"log"
)

func main() {
	info, err := lcu.FindLCUConnectInfo()

	if err != nil {
		if lcu.IsProcessNotFoundError(err) {
			log.Println("No LeagueClient process found!")
			return
		}

		panic(err)
	}

	log.Printf("LeagueClient is running on port %s and you can authenticate with following token: %s", info.Port, info.AuthToken)
}
```

### Create Proxy

If you don't want to use SSL and import the self-signed certificate of the league client you can use the proxy method to
create a reverse proxy. The proxy also automatically adds the authentication header, so you don't have to add it yourself.
This proxy runs then under http and the specified hostname.

```go
package main

import (
	"github.com/ImOlli/go-lcu/proxy"
	"log"
)

func main() {
	// Automatically resolves the port and auth token of the lcu and creates a reverese proxy
	p, err := proxy.CreateProxy(":8080")

	if err != nil {
		panic(err)
	}

	// Optionally you can disable cors or the certificate check
	p.DisableCORS = true
	p.DisableCertCheck = true

	log.Fatal(p.ListenAndServe())
}
```
