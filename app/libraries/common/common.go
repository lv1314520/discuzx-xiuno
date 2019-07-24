package common

import (
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

// IP2Long IP 转整型
func IP2Long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// Long2IP 整型转 IP
func Long2IP(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}

// GetRandomString 生成指定长度随机字符串
func GetRandomString(t string, l int) string {
	pool := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	switch t {

	case "alpha":
		pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	case "alnum":
		pool = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	case "numeric":
		pool = "0123456789"

	case "nozero":
		pool = "123456789"

	case "hex":
		pool = "0123456789abcdefABCDEF"
	}

	bytes := []byte(pool)
	result := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}
