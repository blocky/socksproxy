# SOCKS Proxy

A SOCKS5 server that, when used in conjunction with a
[viproxy](https://github.com/blocky/viproxy), grants networking capabilities to
a program running in a Nitro Server.

The current iteration, however, is coded to either just allow connections to
[Let's Encrypt](https://letsencrypt.org/) for the purpose of getting
self-signed certificates, or it allows all outgoing HTTP connections.

## Usage

The easy way to set up development is to run the server in verbose mode:

    go run main.go --verbose

This will start a SOCKS5 proxy on the default port `:1080`.  From a different
shell, you can use curl to make a request though the proxy with the following
command:

    curl --socks5-hostname localhost:1080 https://acme-v02.api.letsencrypt.org

This command should succeed and return a website.

Next, we can try to make a request that should be blocked:

	curl --socks5-hostname localhost:1080 https://example.com

This command should fail with some message such as:

    curl: (97) Can't complete SOCKS5 connection to example.com. (2)

Next, let's set up our proxy so that it allows all requests

    go run main.go --verbose --fqdn-allow-all

Now, try again, both `curl` requests should now succeed.
