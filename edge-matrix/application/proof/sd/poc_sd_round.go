package sd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"sync"
)

const (
	// min difference bit count
	MinDifferenceBitCount = 1
)

// hashCount:k=ImageHash, v=count
type hashCount struct {
	all   map[string]uint64
	total uint64
}
type PocSDRound struct {
	sync.RWMutex

	logger   hclog.Logger
	blockNum uint64
	seedHash string

	//modelDataCountMap:k=modeHash, v=hashCount
	consensusCountMap map[string]*hashCount

	//modelDataCountMap:k=modeHash, v=imageHash
	consensusResultMap map[string]string

	// poc lookup map keeping track of all
	// poc present in the pool
	pocMap pocLookupMap
}

type PocSdData struct {
	NodeId    string
	ModelHash string
	SeedHash  string
	Md5num    string
	ImageHash string
	BlockNum  uint64
	Power     int64
}

func NewPocSDRound(
	logger hclog.Logger,
	blockNum uint64,
	seedHash string,
) *PocSDRound {
	return &PocSDRound{
		logger:             logger,
		blockNum:           blockNum,
		seedHash:           seedHash,
		consensusCountMap:  make(map[string]*hashCount),
		consensusResultMap: make(map[string]string),
		pocMap:             pocLookupMap{all: make(map[string]*PocSdData)},
	}
}

func (r *PocSDRound) CompleteRound() ([]*PocSdData, error) {
	r.Lock()
	defer r.Unlock()

	err := r.consensus()
	if err != nil {
		return nil, err
	}

	total := r.pocMap.len()
	validatedData := r.validate()

	r.logger.Info("CompleteRound", "valid", len(validatedData), "len", total)
	return validatedData, nil
}

func (r *PocSDRound) Print() error {
	all := r.pocMap.getAll()
	marshalAll, err := json.Marshal(all)
	if err != nil {
		return err
	}
	r.logger.Info("PocSDRound:", "pocMap", string(marshalAll))

	consensusMarshal, err := json.Marshal(r.consensusResultMap)
	if err != nil {
		return err
	}
	r.logger.Info("PocSDRound:", "consensusResultMap", string(consensusMarshal))

	for modelHash, consensusCount := range r.consensusCountMap {
		consensucCountAll, err := json.Marshal(consensusCount.all)
		if err != nil {
			return err
		}
		r.logger.Info("PocSDRound", "consensusCount", fmt.Sprintf("modelHash:%s, total:%d, all:%s", modelHash, consensusCount.total, string(consensucCountAll)))
	}

	return nil
}

func (r *PocSDRound) GetRoundSeed() (blockNum uint64, seedHash string) {
	r.Lock()
	defer r.Unlock()

	blockNum = r.blockNum
	seedHash = r.seedHash
	return
}

func (r *PocSDRound) AddPocData(msg *PocSdData) error {
	r.Lock()
	defer r.Unlock()

	if msg.BlockNum != r.blockNum {
		return errors.New("invalid round blockNum")
	}
	if msg.SeedHash != r.seedHash {
		return errors.New("invalid round seed hash")
	}

	add := r.pocMap.add(msg)
	if !add {
		return errors.New("poc data exist")
	}
	// add to consensusCountMap
	if hashCountData, modelHashExsit := r.consensusCountMap[msg.ModelHash]; modelHashExsit {
		if count, hashCountExsit := hashCountData.all[msg.ImageHash]; hashCountExsit {
			hashCountData.all[msg.ImageHash] = count + 1
		} else {
			hashCountData.all[msg.ImageHash] = 1
		}
		hashCountData.total += 1
	} else {
		hashCountData = &hashCount{total: 1, all: make(map[string]uint64)}
		hashCountData.all[msg.ImageHash] = 1
		r.consensusCountMap[msg.ModelHash] = hashCountData
	}
	return nil
}

func (r *PocSDRound) consensus() error {
	// fill consensusResultMap
	for modelHash, hashCountData := range r.consensusCountMap {
		total := hashCountData.total
		for imageHash, count := range hashCountData.all {
			// validate imageHash count
			if (count * 2) > total {
				r.consensusResultMap[modelHash] = imageHash
				break
			}
		}
	}
	return nil
}

func (r *PocSDRound) validate() []*PocSdData {
	validatedPocData := make([]*PocSdData, 0)

	// validate pocData by consensusResultMap
	for _, data := range r.pocMap.all {
		if validImageHash, hasValid := r.consensusResultMap[data.ModelHash]; hasValid {
			if validImageHash == data.ImageHash {
				validatedPocData = append(validatedPocData, data)
			} else {
				count, err := DifferenceBitCount(validImageHash, data.ImageHash)
				if err == nil && count <= MinDifferenceBitCount {
					validatedPocData = append(validatedPocData, data)
				}
			}
		}
	}
	return validatedPocData
}
