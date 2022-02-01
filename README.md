# lens-examples

**These are several examples on how to implement [Lens](https://github.com/strangelove-ventures/lens) into outside projects**

## **TOC**
* **[Query Balance](https://github.com/strangelove-ventures/lens-examples/tree/main/query_balance)**
    * This is about as basic as it gets. Querys balance of wallet address.

* **[Send Transaction](https://github.com/strangelove-ventures/lens-examples/tree/main/send_transaction)**
    * Create key and send transaction to another address.


* **Indexing Example** - TODO


---


## **Go Mod Setup**

Because Lens replaces some modules in its mod file, manual steps are needed to properly import Lens into a project. 

First, install lens to your `GOPATH/pkg/mod` directory. Run:
```bash
go install github.com/strangelove-ventures/lens@latest
```

Add the following lines to your mod file:
```
require github.com/gogo/protobuf v1.3.3
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
```

Tidy your mod file. Run:
```bash
go mod tidy
```

No you can properly import lens to your go file:
```go
import (
	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
)
````