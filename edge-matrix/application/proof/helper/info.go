package helper

import (
	"encoding/json"
	"fmt"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"net"
)

type CpuInfo struct {
	Cpus      int
	VendorId  string
	Family    string
	Model     string
	Cores     int32
	ModelName string
	Mhz       float64
}

type GpuInfo struct {
	Gpus         int      `json:"gpus"`
	GraphicsCard []string `json:"graphics_card"`
}

func GetGpuInfo() string {
	gpuInfo := GpuInfo{}
	gpu, err := ghw.GPU()
	if err != nil {
		//fmt.Printf("Error getting GPU info: %v", err)
		return ""
	}

	//fmt.Printf("%v\n", gpu)
	gpuInfo.Gpus = len(gpu.GraphicsCards)
	gpuInfo.GraphicsCard = make([]string, gpuInfo.Gpus)
	for i, card := range gpu.GraphicsCards {
		gpuInfo.GraphicsCard[i] = card.String()
	}
	marshaledJson, err := json.Marshal(gpuInfo)
	if err != nil {
		return ""
	}
	return string(marshaledJson)
}

func GetMemInfo() string {
	v, _ := mem.VirtualMemory()

	if v == nil {
		return ""
	}
	// almost every return value is a struct
	return fmt.Sprintf(`{"total": %v, "free":%v, "used_percent":%f}`, v.Total, v.Free, v.UsedPercent)
}

func GetLocalMac() (mac string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, inter := range interfaces {
		//fmt.Println(inter.Name)
		ha := inter.HardwareAddr
		//fmt.Println("MAC ===== ", mac)
		mac = ha.String()
		if len(mac) > 0 {
			return
		}
	}
	return "", nil
}

func GetCpuInfo() string {
	v, err := cpu.Info()
	if err != nil {
		log.Println(fmt.Sprintf("error %s", err.Error()))
	}
	if len(v) == 0 {
		log.Println("could not get CPU Info")
	}
	vv := v[0]
	if vv.ModelName == "" {
		log.Println("could not get CPU ModelName")
	}
	cpu := CpuInfo{
		Cpus:      len(v),
		VendorId:  vv.VendorID,
		Family:    vv.Family,
		Model:     vv.Model,
		Cores:     vv.Cores,
		ModelName: vv.ModelName,
		Mhz:       vv.Mhz,
	}

	marshaledJson, err := json.Marshal(cpu)
	if err != nil {
		return ""
	}
	return string(marshaledJson)
}
