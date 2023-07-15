package helper

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"net"
)
import "github.com/shirou/gopsutil/v3/mem"

type CpuInfo struct {
	Cpus      int
	VendorId  string
	Family    string
	Model     string
	Cores     int32
	ModelName string
	Mhz       float64
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
