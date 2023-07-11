package sd

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRound1(t *testing.T) {
	round := NewPocSDRound(hclog.NewNullLogger(), 10000, "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c")

	cases := []*PocSdData{
		{
			NodeId:    "n1",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1c", //Md5num for ModelHash case1
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n2",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1c", //Md5num for ModelHash case1
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n3",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1c", //Md5num for ModelHash case1
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n4",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1a", // invalid Md5num
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n5",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1a", // invalid Md5num
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n6",
			ModelHash: "e6415c4890",                       // ModelHash case2
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1b", // Md5num for ModelHash case2
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n7",
			ModelHash: "e6415c4890",                       // ModelHash case2
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1b", // Md5num for ModelHash case2
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n8",
			ModelHash: "e6415c4890",                       // ModelHash case2
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1d", // invalid Md5num
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
	}

	for _, data := range cases {
		round.AddPocData(data)
	}

	allValidData, err := round.CompleteRound()
	if err != nil {
		t.Error(err)
		return
	}
	for _, data := range allValidData {
		t.Log(fmt.Sprintf("data: %v", data))
	}
	assert.Equal(t, 5, len(allValidData))
}
func TestRound2(t *testing.T) {
	round := NewPocSDRound(hclog.NewNullLogger(), 10000, "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c")

	cases := []*PocSdData{
		{
			NodeId:    "n1",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1c", //Md5num for ModelHash case1
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
	}

	for _, data := range cases {
		round.AddPocData(data)
	}

	allValidData, err := round.CompleteRound()
	if err != nil {
		t.Error(err)
		return
	}
	for _, data := range allValidData {
		t.Log(fmt.Sprintf("data: %v", data))
	}
	assert.Equal(t, 1, len(allValidData))
}

func TestRound3(t *testing.T) {
	round := NewPocSDRound(hclog.NewNullLogger(), 10000, "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c")

	cases := []*PocSdData{
		{
			NodeId:    "n1",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1a", //Md5num for ModelHash case1
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
		{
			NodeId:    "n2",
			ModelHash: "e6415c4892",                       // ModelHash case1
			Md5num:    "2ed218ca9d99263dae3ee9e7ab5a8d1c", //Md5num for ModelHash case1
			SeedHash:  "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c",
			BlockNum:  10000,
			Power:     150000,
		},
	}

	for _, data := range cases {
		round.AddPocData(data)
	}

	allValidData, err := round.CompleteRound()
	if err != nil {
		t.Error(err)
		return
	}
	for _, data := range allValidData {
		t.Log(fmt.Sprintf("data: %v", data))
	}
	assert.Equal(t, 0, len(allValidData))
}
