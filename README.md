# lens-examples

Helpful info to implement [Lens](https://github.com/strangelove-ventures/lens) into your projects

## **Example Implementations:**
* **[Query Balance](https://github.com/strangelove-ventures/lens-examples/tree/main/query_balance)**
    * This is about as basic as it gets. Querys balance of wallet address.

* **[Send Transaction](https://github.com/strangelove-ventures/lens-examples/tree/main/send_transaction)**
    * Restore key, build transaction, broadcast transaction


* **Indexing Example** - TODO  


---  


## **Go.mod Setup**

**Note:** This is only necessary if you are importing Lens into your project. These steps have already been taken in the above examples. 

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