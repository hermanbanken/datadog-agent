// +build linux_bpf

package tracer

import (
	"net"
	"unsafe"

	ddconfig "github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/network"
	"github.com/DataDog/datadog-agent/pkg/network/config"
	"github.com/DataDog/datadog-agent/pkg/network/ebpf/probes"
	"github.com/DataDog/datadog-agent/pkg/process/util"
	"github.com/DataDog/datadog-agent/pkg/util/ec2"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/DataDog/ebpf"
	"github.com/DataDog/ebpf/manager"
)

type gatewayLookup struct {
	routeCache          network.RouteCache
	subnetCache         map[int]network.Subnet
	subnetForHwAddrFunc func(net.HardwareAddr) (network.Subnet, error)
}

func gwLookupEnabled(config *config.Config) bool {
	// only enabled on AWS currently
	return config.EnableGatewayLookup && ddconfig.IsCloudProviderEnabled(ec2.CloudProviderName)
}

func newGatewayLookup(config *config.Config, runtimeCompilerEnabled bool, m *manager.Manager) *gatewayLookup {
	if !gwLookupEnabled(config) {
		return nil
	}

	var router network.Router
	if runtimeCompilerEnabled {
		router = newEbpfRouter(m)
	} else {
		router = network.NewNetlinkRouter(config.ProcRoot)
	}

	return &gatewayLookup{
		subnetCache:         make(map[int]network.Subnet),
		routeCache:          network.NewRouteCache(512, router),
		subnetForHwAddrFunc: ec2SubnetForHardwareAddr,
	}
}

func (g *gatewayLookup) Lookup(cs *network.ConnectionStats) {
	r, ok := g.routeCache.Get(cs.Source, cs.Dest, cs.NetNS)
	if !ok {
		return
	}

	// if there is no gateway, we don't need to add subnet info
	// for gateway resolution in the backend
	if util.NetIPFromAddress(r.Gw).IsUnspecified() {
		return
	}

	s, ok := g.subnetCache[r.IfIndex]
	if !ok {
		ifi, err := net.InterfaceByIndex(r.IfIndex)
		if err != nil {
			log.Errorf("error getting index for interface index %d: %s", r.IfIndex, err)
			return
		}

		if len(ifi.HardwareAddr) == 0 {
			// can happen for loopback
			return
		}

		if s, err = g.subnetForHwAddrFunc(ifi.HardwareAddr); err != nil {
			log.Errorf("error getting subnet info for interface index %d: %s", r.IfIndex, err)
			return
		}

		g.subnetCache[r.IfIndex] = s
	}

	cs.Via = &network.Via{
		Subnet: s,
	}
}

func ec2SubnetForHardwareAddr(hwAddr net.HardwareAddr) (network.Subnet, error) {
	snet, err := ec2.GetSubnetForHardwareAddr(hwAddr)
	if err != nil {
		return network.Subnet{}, err
	}

	return network.Subnet{Alias: snet.ID}, nil
}

type ebpfRouter struct {
	gwMp *ebpf.Map
}

func newEbpfRouter(m *manager.Manager) network.Router {
	mp, ok, err := m.GetMap(string(probes.GatewayMap))
	if err != nil || !ok {
		return nil
	}
	return &ebpfRouter{
		gwMp: mp,
	}
}

func (b *ebpfRouter) Route(source, dest util.Address, netns uint32) (network.Route, bool) {
	d := newIPRuoteDest(source, dest, netns)
	gw := &ipRouteGateway{}
	if err := b.gwMp.Lookup(unsafe.Pointer(d), unsafe.Pointer(gw)); err != nil || gw.ifIndex() == 0 {
		return network.Route{}, false
	}

	return network.Route{Gw: gw.gateway(), IfIndex: gw.ifIndex()}, true
}