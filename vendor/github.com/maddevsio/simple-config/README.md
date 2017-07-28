# simple-config
Simple wrapper around viper package

#Install
go get -v github.com/maddevsio/simple-config

#Usage:
```go
config := NewSimpleConfig("./config.test", "yml")
value := config.Get("testkey")
```
#config.test.yaml is:

```yaml
testkey: test value
```

To redefine "testkey" via env set TESTKEY env var
