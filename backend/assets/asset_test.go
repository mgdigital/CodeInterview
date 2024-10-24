package assets

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestAsset(t *testing.T) {
	asset := newAsset(
		10,
		"veritatis.aysZaos.biz",
		"Aut sit voluptatem perferendis consequatur accusantium.",
		"Dr. Marcia Cole",
	)
	asset.addIP("47.21.29.165")
	asset.addPort(36577)
	assert.Equal(t, Asset{
		ID:        10,
		Host:      "veritatis.aysZaos.biz",
		Comment:   "Aut sit voluptatem perferendis consequatur accusantium.",
		Owner:     "Dr. Marcia Cole",
		Signature: "56420776427ead533eb97f9c1b1d9cb5d573276d0104048292fdc70d0d14162f",
		IPs: []IP{
			{
				Address:   "47.21.29.165",
				Signature: "ad43494b46ab8c9e14bcc57c1d349a4883ffed4a8325a6c329e64337f3e02ccf",
			},
		},
		Ports: []Port{
			{
				Port:      36577,
				Signature: "481a66f2881b2d0723b38670d71c554057bd8af3e4bc72056dd5cde0a7cc7e6b",
			},
		},
	}, *asset)
}
