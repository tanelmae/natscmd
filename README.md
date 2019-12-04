# natscmd
Simple NATS client for quick and dirty manual configuration testing.

```
Usage:
  -certs string
    	Path to directory with client.pem, client-key.pem and ca.pem
  -cmd string
    	pub or sub (default "sub")
  -creds string
    	Path to credentials file
  -message string
    	Message to send (default "Hello World")
  -nats string
    	NATS URL (default "nats://localhost:4222")
  -subject string
    	NATS subject to use (default "test")
  -timeout duration
    	Subscriber timeout (default 30s)
```


#### Build:
```
$ go mod vendor
$ go build -mod=readonly natscmd.go
```


#### Release downloads:
https://github.com/tanelmae/natscmd/releases

#### Install with Homebrew:

Binaries provided for Darwin_x86_64 and Linux_x86_64
```
brew install tanelmae/brew/natscmd
