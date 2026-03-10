package resolver

import (
	"fmt"
	"net"
	"strings"
)

const DefaultMinecraftPort = 25565

type SRVLookupFunc func(service, proto, name string) (string, []*net.SRV, error)

type TargetResolution struct {
	Host               string
	Port               int
	ExplicitPort       bool
	DirectIP           bool
	SRVLookupAttempted bool
	UsedSRV            bool
	SRVNotFound        bool
}

func ResolveTarget(target string, lookupSRV SRVLookupFunc) (TargetResolution, error) {
	res := TargetResolution{Host: target, Port: DefaultMinecraftPort}

	if lookupSRV == nil {
		lookupSRV = net.LookupSRV
	}

	explicitRes, explicit, err := resolveExplicitHostPort(target)
	if err != nil {
		return TargetResolution{}, err
	}
	if explicit {
		res = explicitRes
	} else if net.ParseIP(target) != nil {
		res.DirectIP = true
	} else if strings.Count(target, ":") > 1 {
		return TargetResolution{}, fmt.Errorf("dirección IPv6 mal formateada: usa el formato '[addr]:port'")
	} else if err := resolveSRV(&res, target, lookupSRV); err != nil {
		return TargetResolution{}, err
	}

	if err := validatePort(res.Port); err != nil {
		return TargetResolution{}, err
	}

	return res, nil
}

func resolveExplicitHostPort(target string) (TargetResolution, bool, error) {
	h, p, err := net.SplitHostPort(target)
	if err != nil {
		return TargetResolution{}, false, nil
	}

	res := TargetResolution{Host: h, Port: DefaultMinecraftPort, ExplicitPort: true}
	if _, err := fmt.Sscanf(p, "%d", &res.Port); err != nil {
		return TargetResolution{}, true, fmt.Errorf("puerto inválido en el objetivo: %v", err)
	}
	return res, true, nil
}

func resolveSRV(res *TargetResolution, target string, lookupSRV SRVLookupFunc) error {
	res.SRVLookupAttempted = true
	_, addrs, err := lookupSRV("minecraft", "tcp", target)
	if err == nil && len(addrs) > 0 {
		res.Host = strings.TrimSuffix(addrs[0].Target, ".")
		res.Port = int(addrs[0].Port)
		res.UsedSRV = true
		return nil
	}

	if err == nil {
		res.SRVNotFound = true
		return nil
	}

	if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
		res.SRVNotFound = true
		return nil
	}

	return fmt.Errorf("error de resolución DNS: %v", err)
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("puerto fuera de rango: %d", port)
	}
	return nil
}
