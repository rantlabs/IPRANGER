package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func uint32ToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

func parseCIDR(input string) (*net.IPNet, net.IP, error) {
	// Support dotted-decimal mask notation: 192.168.1.0/255.255.255.0
	if strings.Count(input, "/") == 1 {
		parts := strings.SplitN(input, "/", 2)
		mask := parts[1]
		// If mask looks like an IP address, convert to prefix length
		if strings.Contains(mask, ".") {
			maskIP := net.ParseIP(mask).To4()
			if maskIP == nil {
				return nil, nil, fmt.Errorf("invalid mask: %s", mask)
			}
			ones, _ := net.IPMask(maskIP).Size()
			input = fmt.Sprintf("%s/%d", parts[0], ones)
		}
	}

	ip, network, err := net.ParseCIDR(input)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid address/mask: %s", err)
	}
	return network, ip, nil
}

func main() {
	listAll := flag.Bool("list", false, "List all usable host addresses (excludes network and broadcast)")
	listNet := flag.Bool("network", false, "Print only the network address")
	listBcast := flag.Bool("broadcast", false, "Print only the broadcast address")
	listRange := flag.Bool("range", false, "Print the usable address range in summary form (first - last)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: ipranger [options] <address/mask>\n\n")
		fmt.Fprintf(os.Stderr, "  address/mask    CIDR notation (e.g. 192.168.1.0/24)\n")
		fmt.Fprintf(os.Stderr, "                  or dotted-mask notation (e.g. 192.168.1.0/255.255.255.0)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Support "IP mask" as two separate args in addition to "IP/mask"
	var input string
	if flag.NArg() == 2 {
		input = flag.Arg(0) + "/" + flag.Arg(1)
	} else {
		input = flag.Arg(0)
	}
	network, hostIP, err := parseCIDR(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	mask := network.Mask
	netAddr := network.IP.To4()
	ones, bits := mask.Size()
	totalIPs := uint32(1) << uint(bits-ones)

	netInt := ipToUint32(netAddr)
	bcastInt := netInt | ^binary.BigEndian.Uint32(mask)
	bcastAddr := uint32ToIP(bcastInt)

	gatewayInt := netInt + 1
	gatewayAddr := uint32ToIP(gatewayInt)

	firstUsable := gatewayAddr
	lastUsableInt := bcastInt - 1
	lastUsable := uint32ToIP(lastUsableInt)

	var usableCount int64
	if totalIPs >= 2 {
		usableCount = int64(totalIPs) - 2
	}

	// Handle special cases for /31 and /32
	is31 := ones == 31
	is32 := ones == 32

	// Exclusive modes
	if *listRange {
		if is32 {
			fmt.Printf("%s - %s  (%d host)\n", netAddr, netAddr, 1)
		} else if is31 {
			fmt.Printf("%s - %s  (%d hosts)\n", netAddr, bcastAddr, 2)
		} else {
			fmt.Printf("%s - %s  (%d hosts)\n", firstUsable, lastUsable, usableCount)
		}
		return
	}
	if *listNet {
		fmt.Println(netAddr)
		return
	}
	if *listBcast {
		if is32 {
			fmt.Println(netAddr)
		} else {
			fmt.Println(bcastAddr)
		}
		return
	}
	if *listAll {
		if is32 {
			fmt.Println(netAddr)
		} else if is31 {
			fmt.Println(netAddr)
			fmt.Println(bcastAddr)
		} else {
			for i := gatewayInt; i <= lastUsableInt; i++ {
				fmt.Println(uint32ToIP(i))
			}
		}
		return
	}

	// Default: summary output
	_ = hostIP
	maskDotted := net.IP(mask).String()

	fmt.Printf("Address    : %s\n", hostIP)
	fmt.Printf("Netmask    : %s = %d\n", maskDotted, ones)
	fmt.Printf("Network    : %s/%d\n", netAddr, ones)
	if !is32 {
		fmt.Printf("Broadcast  : %s\n", bcastAddr)
	}
	fmt.Printf("Gateway    : %s\n", gatewayAddr)
	if !is31 && !is32 {
		fmt.Printf("HostMin    : %s\n", firstUsable)
		fmt.Printf("HostMax    : %s\n", lastUsable)
		fmt.Printf("Hosts/Net  : %d\n", usableCount)
	} else if is31 {
		fmt.Printf("Hosts/Net  : 2  (point-to-point)\n")
	} else {
		fmt.Printf("Hosts/Net  : 1  (single host)\n")
	}
}
