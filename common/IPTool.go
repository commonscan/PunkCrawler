package common

import (
	"github.com/asaskevich/govalidator"
	"net"
)

// cidr to ip range
func GenerateIP(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IsCIDR(string2 string) bool {
	return govalidator.IsCIDR(string2)
}
func IsIp(string2 string) bool {
	return govalidator.IsIP(string2)
}
