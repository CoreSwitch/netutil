// Copyright 2016 CoreSwitch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netutil

import (
	"net"
	"strconv"
	"strings"
)

var MaskBits = []byte{0x00, 0x80, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfe, 0xff}

const (
	AFI_IP = iota
	AFI_IP6
	AFI_MAX
)

type Prefix struct {
	net.IP
	Length int
}

func dupIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

func (p *Prefix) AFI() int {
	if len(p.IP) == net.IPv4len {
		return AFI_IP
	}
	if len(p.IP) == net.IPv6len {
		return AFI_IP6
	}
	return AFI_MAX
}

func NewPrefixAFI(afi int) *Prefix {
	switch afi {
	case AFI_IP:
		return &Prefix{IP: make(net.IP, net.IPv4len), Length: 0}
	case AFI_IP6:
		return &Prefix{IP: make(net.IP, net.IPv6len), Length: 0}
	default:
		return nil
	}
}

func ParsePrefix(s string) (*Prefix, error) {
	i := strings.IndexByte(s, '/')
	if i < 0 {
		return nil, &net.ParseError{Type: "Prefix address", Text: s}
	}

	addr, mask := s[:i], s[i+1:]

	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, &net.ParseError{Type: "Prefix address", Text: s}
	}

	ip4 := ip.To4()
	if ip4 != nil {
		ip = ip4
	}

	length, err := strconv.Atoi(mask)
	if err != nil {
		return nil, err
	}
	return &Prefix{IP: ip, Length: length}, nil
}

func (p *Prefix) String() string {
	return p.IP.String() + "/" + strconv.Itoa(p.Length)
}

func (p *Prefix) ApplyMask() {
	i := p.Length / 8

	if i >= len(p.IP) {
		return
	}

	offset := p.Length % 8
	p.IP[i] &= MaskBits[offset]
	i++

	for i < len(p.IP) {
		p.IP[i] = 0
		i++
	}
}

func (p *Prefix) Copy() *Prefix {
	return &Prefix{IP: dupIP(p.IP), Length: p.Length}
}

func (p *Prefix) Equal(x *Prefix) bool {
	if p.IP.Equal(x.IP) && p.Length == x.Length {
		return true
	} else {
		return false
	}
}

func PrefixFromIPNet(net net.IPNet) *Prefix {
	ip := net.IP.To4()
	if ip == nil {
		ip = net.IP
	}
	len, _ := net.Mask.Size()
	return &Prefix{IP: dupIP(ip), Length: len}
}

func PrefixFromIPPrefixlen(ip net.IP, len int) *Prefix {
	return &Prefix{IP: dupIP(ip), Length: len}
}

func ParseIPv4(s string) net.IP {
	ip := net.ParseIP(s)
	ip4 := ip.To4()
	if ip4 != nil {
		return ip4
	}
	return ip
}