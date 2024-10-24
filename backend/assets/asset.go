package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type Asset struct {
	ID        int    `json:"id"`
	Host      string `json:"host"`
	Comment   string `json:"comment"`
	Owner     string `json:"owner"`
	IPs       []IP   `json:"ips"`
	Ports     []Port `json:"ports"`
	Signature string `json:"signature"`
}

func newAsset(id int, host, comment, owner string) *Asset {
	return &Asset{
		ID:        id,
		Host:      host,
		Comment:   comment,
		Owner:     owner,
		IPs:       []IP{},
		Ports:     []Port{},
		Signature: newSignature(host + comment + owner),
	}
}

func (a *Asset) addIP(address string) {
	a.IPs = append(a.IPs, *newIP(address))
}

func (a *Asset) addPort(portNum int) {
	a.Ports = append(a.Ports, *newPort(portNum))
}

type IP struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

func newIP(address string) *IP {
	return &IP{
		Address:   address,
		Signature: newSignature(address),
	}
}

type Port struct {
	Port      int    `json:"port"`
	Signature string `json:"signature"`
}

func newPort(portNum int) *Port {
	return &Port{
		Port:      portNum,
		Signature: newSignature(strconv.Itoa(portNum)),
	}
}

func newSignature(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
