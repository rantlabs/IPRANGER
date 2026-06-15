# IPRANGER

## A command line utility to provide easy access IP address subnet information

```
ipranger -h
Usage: ipranger [options] <address/mask>

  address/mask    CIDR notation (e.g. 192.168.1.0/24)
                  or dotted-mask notation (e.g. 192.168.1.0/255.255.255.0)

Options:
  -broadcast
    	Print only the broadcast address
  -list
    	List all usable host addresses (excludes network and broadcast)
  -network
    	Print only the network address
  -range
    	Print the usable address range in summary form (first - last)
```

## Executables available for download

```
ipranger_linux_32		ipranger_rpi_arm64		ipranger_windows_64.exe
ipranger_linux_64		ipranger_rpi_armv6		
ipranger_intel_mac		ipranger_rpi_armv7
ipranger_mac_arm64		ipranger_windows_32.exe
```

## Go source code
```
main.go
```

#### Command ipranger 192.168.1.0/27 
```
ipranger 192.168.1.0/27 
Address    : 192.168.1.0
Netmask    : 255.255.255.224 = 27
Network    : 192.168.1.0/27
Broadcast  : 192.168.1.31
Gateway    : 192.168.1.1
HostMin    : 192.168.1.1
HostMax    : 192.168.1.30
Hosts/Net  : 30
```

#### Command ipranger -network 192.168.1.129/27
```
ipranger -network 192.168.1.129/27
192.168.1.128
```
#### Command ipranger -broadcast 192.168.1.129/27
```
ipranger -broadcast 192.168.1.129/27
192.168.1.159
```
#### Command ipranger -range 192.168.1.129/27
```
ipranger -range 192.168.1.129/27
192.168.1.129 - 192.168.1.158  (30 hosts)
```
#### Command ipranger -list 192.168.10.68/29 - Lists all useable addresses
```
ipranger -list 192.168.10.68/29 
192.168.10.65
192.168.10.66
192.168.10.67
192.168.10.68
192.168.10.69
192.168.10.70
```
