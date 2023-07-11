package sd

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"sync"
)

// md5Count:k=Md5num, v=count
type md5Count struct {
	all   map[string]uint64
	total uint64
}
type PocSDRound struct {
	sync.RWMutex

	logger   hclog.Logger
	blockNum uint64
	seedHash string

	//modelDataCountMap:k=modeHash, v=md5Count
	consensusCountMap map[string]*md5Count

	//modelDataCountMap:k=modeHash, v=Md5num
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
		consensusCountMap:  make(map[string]*md5Count),
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

	r.pocMap.add(msg)

	// add to consensusCountMap
	if md5CountData, modelHashExsit := r.consensusCountMap[msg.ModelHash]; modelHashExsit {
		if count, md5CountExsit := md5CountData.all[msg.Md5num]; md5CountExsit {
			md5CountData.all[msg.Md5num] = count + 1
		} else {
			md5CountData.all[msg.Md5num] = 1
		}
		md5CountData.total += 1
	} else {
		md5CountData = &md5Count{total: 1, all: make(map[string]uint64)}
		md5CountData.all[msg.Md5num] = 1
		r.consensusCountMap[msg.ModelHash] = md5CountData
	}
	return nil
}

func (r *PocSDRound) consensus() error {
	// fill consensusResultMap
	for modelHash, md5CountData := range r.consensusCountMap {
		total := md5CountData.total
		for md5num, count := range md5CountData.all {
			// validate Md5num
			if (count * 2) > total {
				r.consensusResultMap[modelHash] = md5num
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
		if validMd5num, hasValid := r.consensusResultMap[data.ModelHash]; hasValid {
			if validMd5num == data.Md5num {
				validatedPocData = append(validatedPocData, data)
			}
		}
	}
	return validatedPocData
}
