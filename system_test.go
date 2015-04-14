package main

import (
	"net"
	"testing"
)

func TestSystem(t *testing.T) {
	if sys, err := System(); err == nil {
		t.Logf("%+v\n", sys)
	} else {
		t.Error(err)
	}
}

func TestDefaultIPv4Route(t *testing.T) {
	checkIfi := func(ifi *net.Interface) {
		if ifi == nil {
			t.Error("interface is nil")
		}
	}
	checkGatewayIP := func(gwip net.IP) {
	}

	ifi, gwip, err := DefaultIPv4Route()
	if err == nil {
		t.Logf("iface:   %s\n", ifi.Name)
		checkIfi(ifi)
		t.Logf("gateway: %s\n", gwip.String())
		checkGatewayIP(gwip)
	} else {
		t.Error(err)
	}
}
