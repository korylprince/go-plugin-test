This is a very simple test of [HashiCorp's go-plugin](https://github.com/hashicorp/go-plugin) with bidirectional communication over GRPC.

**How to Run**

```
git clone github.com/korylprince/go-plugin-test.git
cd go-plugin-test
cd cmd/greeter-plugin
go build
cd ../greeter/
go run greeter.go "Mr. Anderson"

<snip debug logs>

Yo, what's up Mr. Anderson? <nil>
```
