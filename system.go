package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"unsafe"
)

var (
	NoDefaultRoute error = errors.New("no default route")
)

func System() (map[string]interface{}, error) {
	var err error
	hostname := "localhost"
	ipv4addr := "127.0.0.1"
	ipv4net := "127.0.0.0/8"
	ipv4gw := ""

	if hostname, err = os.Hostname(); err != nil {
		return nil, err
	}

	if ifi, gw, err := DefaultIPv4Route(); err == nil {
		ipv4gw = gw.String()
		if addrs, err := ifi.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ipv4 := ipnet.IP.To4(); ipv4 != nil {
						ipv4addr = ipv4.String()
						ipv4net = (&net.IPNet{ipnet.IP.Mask(ipnet.Mask), ipnet.Mask}).String()
						break
					}
				}
			}
		}
	} else if err != NoDefaultRoute {
		return nil, err
	}

	return map[string]interface{}{
		"hostname": hostname,
		"address":  ipv4addr,
		"network":  ipv4net,
		"gateway":  ipv4gw,
	}, nil
}

// DefaultIPv4Route retrieves the default IPv4 route and returns its interface
// and gateway IP. An error is returned if none can be determined.
func DefaultIPv4Route() (*net.Interface, net.IP, error) {
	data, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		return nil, net.IP{}, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines[1:] {
		columns := strings.Split(line, "\t")
		if len(columns) < 3 || columns[0] == "" || columns[1] != "00000000" {
			continue
		}
		ifi, err := net.InterfaceByName(columns[0])
		if err != nil {
			continue
		}
		ip, err := HexToIP(columns[2])
		if err != nil {
			return nil, net.IP{}, err
		}
		return ifi, ip, nil
	}
	return nil, net.IP{}, NoDefaultRoute
}

// IsLittleEndian returns `true` if the host system is little endian.
func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

// HextToIP takes a hex string from the routing table and returns an IP.
func HexToIP(hexIP string) (net.IP, error) {
	bytes, err := hex.DecodeString(hexIP)
	if err != nil {
		return net.IPv4(0, 0, 0, 0), err
	}
	if IsLittleEndian() {
		ipInt := binary.LittleEndian.Uint32(bytes)
		binary.BigEndian.PutUint32(bytes, ipInt)
	}
	return net.IP(bytes), nil
}
