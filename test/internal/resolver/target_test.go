package resolver_test

import (
	"MinecraftCrawler/internal/resolver"
	"errors"
	"net"
	"testing"
)

const (
	resolverTestHost       = "mc.example.net"
	resolverTargetErrFmt   = "ResolveTarget() error = %v"
)

func TestResolveTargetExplicitPort(t *testing.T) {
	res, err := resolver.ResolveTarget(resolverTestHost+":25570", nil)
	if err != nil {
		t.Fatalf(resolverTargetErrFmt, err)
	}
	if res.Host != resolverTestHost || res.Port != 25570 {
		t.Fatalf("ResolveTarget() host/port = %s:%d, want %s:25570", res.Host, res.Port, resolverTestHost)
	}
	if !res.ExplicitPort {
		t.Fatalf("ExplicitPort = false, want true")
	}
}

func TestResolveTargetDirectIP(t *testing.T) {
	res, err := resolver.ResolveTarget("192.168.1.10", nil)
	if err != nil {
		t.Fatalf(resolverTargetErrFmt, err)
	}
	if !res.DirectIP {
		t.Fatalf("DirectIP = false, want true")
	}
	if res.SRVLookupAttempted {
		t.Fatalf("SRVLookupAttempted = true, want false")
	}
	if res.Port != resolver.DefaultMinecraftPort {
		t.Fatalf("Port = %d, want %d", res.Port, resolver.DefaultMinecraftPort)
	}
}

func TestResolveTargetMalformedIPv6(t *testing.T) {
	_, err := resolver.ResolveTarget("2001:db8::1:25565", nil)
	if err == nil {
		t.Fatalf("ResolveTarget() error = nil, want malformed ipv6 error")
	}
}

func TestResolveTargetSRVSuccess(t *testing.T) {
	lookup := func(service, proto, name string) (string, []*net.SRV, error) {
		if service != "minecraft" || proto != "tcp" || name != resolverTestHost {
			return "", nil, errors.New("unexpected lookup arguments")
		}
		return "", []*net.SRV{{Target: "srv.example.net.", Port: 25566}}, nil
	}

	res, err := resolver.ResolveTarget(resolverTestHost, lookup)
	if err != nil {
		t.Fatalf(resolverTargetErrFmt, err)
	}
	if !res.SRVLookupAttempted || !res.UsedSRV {
		t.Fatalf("SRV flags = attempted:%v used:%v, want true/true", res.SRVLookupAttempted, res.UsedSRV)
	}
	if res.Host != "srv.example.net" || res.Port != 25566 {
		t.Fatalf("ResolveTarget() host/port = %s:%d, want srv.example.net:25566", res.Host, res.Port)
	}
}

func TestResolveTargetSRVNotFound(t *testing.T) {
	lookup := func(service, proto, name string) (string, []*net.SRV, error) {
		return "", nil, &net.DNSError{IsNotFound: true}
	}

	res, err := resolver.ResolveTarget(resolverTestHost, lookup)
	if err != nil {
		t.Fatalf(resolverTargetErrFmt, err)
	}
	if !res.SRVNotFound {
		t.Fatalf("SRVNotFound = false, want true")
	}
	if res.Port != resolver.DefaultMinecraftPort {
		t.Fatalf("Port = %d, want %d", res.Port, resolver.DefaultMinecraftPort)
	}
}

func TestResolveTargetSRVLookupError(t *testing.T) {
	lookup := func(service, proto, name string) (string, []*net.SRV, error) {
		return "", nil, errors.New("dns failed")
	}

	_, err := resolver.ResolveTarget(resolverTestHost, lookup)
	if err == nil {
		t.Fatalf("ResolveTarget() error = nil, want dns error")
	}
}
