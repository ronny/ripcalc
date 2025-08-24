package ipv4

import (
	"fmt"
	"net"
)

type addressType int

const (
	addressTypePublic addressType = iota
	addressTypePrivate
	addressTypeSharedAddressSpace
	addressTypeLinkLocal
	addressTypeLoopback
	addressTypeMulticast
)

func (at addressType) String() string {
	switch at {
	case addressTypePublic:
		return "Public Internet"
	case addressTypePrivate:
		return "Private Internet"
	case addressTypeSharedAddressSpace:
		return "Shared Address Space"
	case addressTypeLinkLocal:
		return "Link Local"
	case addressTypeLoopback:
		return "Loopback"
	case addressTypeMulticast:
		return "Multicast"
	default:
		return "Unknown"
	}
}

type addressRange struct {
	network *net.IPNet
	typ     addressType
}

var specialRanges = []addressRange{
	{mustParseCIDR("192.168.0.0/16"), addressTypePrivate},
	{mustParseCIDR("172.16.0.0/12"), addressTypePrivate},
	{mustParseCIDR("10.0.0.0/8"), addressTypePrivate},
	{mustParseCIDR("100.64.0.0/10"), addressTypeSharedAddressSpace},
	{mustParseCIDR("169.254.0.0/16"), addressTypeLinkLocal},
	{mustParseCIDR("127.0.0.0/8"), addressTypeLoopback},
	{mustParseCIDR("224.0.0.0/4"), addressTypeMulticast},
}

func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(fmt.Sprintf("invalid CIDR in mustParseCIDR: %s", cidr))
	}

	return network
}

type Network struct {
	Address      net.IP
	PrefixLength int
	Netmask      net.IPMask
	Wildcard     net.IP
	Network      net.IP
	Broadcast    net.IP
	HostMin      net.IP
	HostMax      net.IP
	HostCount    uint32
	Class        string
	Type         string
}

func ParseCIDR(cidr string) (*Network, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("net.ParseCIDR: %w", err)
	}

	if ip.To4() == nil {
		return nil, fmt.Errorf("%w: not an IPv4 address", ErrInvalidAddress)
	}

	prefixLen, _ := ipNet.Mask.Size()

	return &Network{
		Address:      ip.To4(),
		PrefixLength: prefixLen,
	}, nil
}

func (n *Network) String() string {
	return fmt.Sprintf("%s/%d", n.Address, n.PrefixLength)
}

func (n *Network) Calculate() error {
	if n.Address == nil {
		return fmt.Errorf("%w: address is nil", ErrInvalidAddress)
	}

	n.Netmask = net.CIDRMask(n.PrefixLength, 32)
	n.Wildcard = invertMask(net.IP(n.Netmask))
	n.Network = n.Address.Mask(n.Netmask)
	n.Broadcast = calculateBroadcast(n.Network, n.Wildcard)
	n.HostMin, n.HostMax = calculateHostRange(n.Network, n.Broadcast)
	n.HostCount = calculateHostCount(n.PrefixLength)
	n.Class = classifyAddress(n.Address)
	n.Type = classifyAddressType(n.Address).String()

	return nil
}

func (n *Network) FormattedText() string {
	addressBinary := FormatBinary(n.Address)
	netmaskBinary := FormatBinary(net.IP(n.Netmask))
	wildcardBinary := FormatBinary(n.Wildcard)
	networkBinary := FormatBinary(n.Network)
	hostMinBinary := FormatBinary(n.HostMin)
	hostMaxBinary := FormatBinary(n.HostMax)
	broadcastBinary := FormatBinary(n.Broadcast)

	typeStr := n.Type

	return fmt.Sprintf("Address:   %-15s  %s\n"+
		"Netmask:   %-15s = %-2d   %s\n"+
		"Wildcard:  %-15s  %s\n"+
		"--------------------------------------------------------------------\n"+
		"Network:   %s/%-2d       %s\n"+
		" HostMin:  %-15s  %s\n"+
		" HostMax:  %-15s  %s\n"+
		"Broadcast: %-15s  %s\n"+
		"Hosts:     %-15d  Class %s, %s",
		n.Address.String(), addressBinary,
		net.IP(n.Netmask).String(), n.PrefixLength, netmaskBinary,
		n.Wildcard.String(), wildcardBinary,
		n.Network.String(), n.PrefixLength, networkBinary,
		n.HostMin.String(), hostMinBinary,
		n.HostMax.String(), hostMaxBinary,
		n.Broadcast.String(), broadcastBinary,
		n.HostCount, n.Class, typeStr,
	)
}

func invertMask(mask net.IP) net.IP {
	wildcard := make(net.IP, 4)
	for i := range 4 {
		wildcard[i] = ^mask[i]
	}

	return wildcard
}

func calculateBroadcast(network, wildcard net.IP) net.IP {
	broadcast := make(net.IP, 4)
	for i := range 4 {
		broadcast[i] = network[i] | wildcard[i]
	}

	return broadcast
}

func calculateHostRange(network, broadcast net.IP) (net.IP, net.IP) {
	hostMin := make(net.IP, 4)
	hostMax := make(net.IP, 4)

	copy(hostMin, network)
	copy(hostMax, broadcast)

	// Host min is network + 1
	hostMin[3]++

	// Host max is broadcast - 1
	hostMax[3]--

	return hostMin, hostMax
}

func calculateHostCount(prefixLen int) uint32 {
	hostBits := 32 - prefixLen
	if hostBits <= 1 {
		return 0
	}

	return (1 << hostBits) - 2 // -2 for network and broadcast
}

func classifyAddress(ip net.IP) string {
	firstOctet := ip[0]
	switch {
	case firstOctet <= 127:
		return "A"
	case firstOctet <= 191:
		return "B"
	case firstOctet <= 223:
		return "C"
	case firstOctet <= 239:
		return "D"
	default:
		return "E"
	}
}

func classifyAddressType(ip net.IP) addressType {
	for _, r := range specialRanges {
		if r.network.Contains(ip) {
			return r.typ
		}
	}

	return addressTypePublic
}

func FormatBinary(ip net.IP) string {
	if len(ip) != 4 {
		return ""
	}

	return fmt.Sprintf("%08b.%08b.%08b.%08b", ip[0], ip[1], ip[2], ip[3])
}
