# ippygo

This is a re-write of https://github.com/shivammathur/IPpy which is a parallel testing of IP addresses and domains in python. Reads IP addresses and domains from a CSV file and gives two lists of working and not working ones.
Did it just for fun to have a Go -version and maybe benchmark a little bit.

## About

Written in Go
Testing of IPs and domains is done in parallel.

Currently Mac and Linux are supported (haven't tried Windows)
Supports both IPv4 and IPv6 IPs and domain names.

### Examples

127.0.0.1
::1
localhost

## Install

cd $GOPATH
go get -u github.com/dvwallin/ippygo

## Run

$GOPATH/bin/ippygo <the-ip-file>
