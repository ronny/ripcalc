# ripcalc

ripcalc is a CLI tool and a Go package for calculating IPv4 and IPv6 address blocks.

## Features

- IPv4 and IPv6 support
- Address start, end, usable address count, network count
- Binary representation of addresses
- Structured output: JSON
- Optional coloured output
- No dependencies beyond the Go standard library

## Usage

## IPv4

```sh
ripcalc 192.168.0.0/24
```

```
   Address:	192.168.0.1         	11000000.10101000.00000000. 00000001
   Netmask:	255.255.255.0   = 24	11111111.11111111.11111111. 00000000
  Wildcard:	0.0.0.255           	00000000.00000000.00000000. 11111111
----------------------------------------------------------------------------
   Network:	192.168.0.0/24      	11000000.10101000.00000000. 00000000
First host:	192.168.0.1         	11000000.10101000.00000000. 00000001
 Last host:	192.168.0.254       	11000000.10101000.00000000. 11111110
 Broadcast:	192.168.0.255       	11000000.10101000.00000000. 11111111
Host count:	254                 	Class C, Private Internet
```

```sh
ripcalc -json 192.168.0.1/24
```

```json
{
  "$schema": "https://github.com/ronny/ripcalc/blob/main/schema/ipv4-v1.json",
  "address": "192.168.0.1",
  "netmask": "255.255.255.0",
  "wildcard": "0.0.0.255",
  "network": {
    "address": "192.168.0.0",
    "prefix_length": "24",
    "broadcast": "192.168.0.255",
    "usable_hosts": {
      "min": "192.168.0.1",
      "max": "192.168.0.254"
      "count": 254
    },
    "class": "C",
    "type": "PRIVATE"
  },
```

## IPv6

```sh
ripcalc 2001:db8::/32
```

TODO

## Installation

Download the binary from the [releases page](https://github.com/ronny/ripcalc/releases).

## Other similar tools

- [ipcalc](https://github.com/kjokjo/ipcalc) - the main inspiration for ripcalc
