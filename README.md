# Go Proxy setup tool

This is a simple tool for toggling company proxy settings on and off in the shell environment on MacOS.

I previously made a version of [this tool in python](https://gist.github.com/igilham/f2fb41190258f385eebe48a6c46421e5).

Use with `eval` in the shell to make changes to your configuration, e.g.
`eval $(goproxy on)`.

## Supported commands

* `on`: turn on proxy environment settings
* `off`: turn off proxy environment settings
* `status`: print the current network location (proxies request to `networksetup`)
* `reset`: reconfigure current proxy environment settings without changing location

## Configuration

Edit `myvars.go` to ensure the correct hosts, port and network locations are set before building the binary.

## Usage in a shell

I have the following in my `~/.profile` file. I run the `reset` command after setting up my variables etc. so that every new shell session loads with the correct environment by the time it shows an interactive prompt.

```shell
alias proxy="~/bin/goproxy"

proxy_on() {
  eval "$(proxy on)"
}
proxy_off() {
  eval "$(proxy off)"
}
proxy_reset() {
  eval "$(proxy reset)"
}
# various setup actions...

proxy_reset
```
