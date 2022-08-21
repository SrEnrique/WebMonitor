package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// Disk is a structure for disk information
type Disk struct {
	Device string  `json:"Device"`
	Total  float64 `json:"Total"`
	Used   float64 `json:"Used"`
	Free   float64 `json:"Free"`
}

// SystemInformatoin is a structure for system information
type SystemInformatoin struct {
	Host   string `json:"host"`
	Memory struct {
		Total float64 `json:"Total"`
		Used  float64 `json:"Used"`
		Free  float64 `json:"Free"`
	} `json:"memory"`
	CPU struct {
		Percent float64 `json:"Percent"`
	} `json:"cpu"`
	Disk []Disk `json:"Disk"`
}

// Windows Service
const serviceName = "GoWebServerMonitor"
const serviceDescription = "Simple web server monitor"

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	for {
		r := gin.Default()

		//Return Json SystemInformation
		r.GET("/", func(c *gin.Context) {
			hostStat, _ := host.Info()
			vmStat, _ := mem.VirtualMemory()

			//Create Json
			system_info := new(SystemInformatoin)

			//Set Hostname
			system_info.Host = hostStat.Hostname

			//Set Memory Info
			system_info.Memory.Total = float64(vmStat.Total / 1024 / 1024)
			system_info.Memory.Used = float64(vmStat.Used / 1024 / 1024)
			system_info.Memory.Free = float64(vmStat.Free / 1024 / 1024)

			//Set one secon percen cpu usage
			system_info.CPU.Percent = getCpuInfo()[0]

			//Set Partitions or disk units in the system
			diskpartitions, _ := disk.Partitions(true)

			disks := []Disk{}

			//iterate in partitions for set information
			for _, parts := range diskpartitions {

				var ddisk Disk
				//Set PArtitions info

				//Set device name
				ddisk.Device = parts.Device
				devicesStats, _ := disk.Usage(parts.Device)
				//Set Usage espesific partition
				ddisk.Total = float64(devicesStats.Total / 1024 / 1024)
				ddisk.Used = float64(devicesStats.Used / 1024 / 1024)
				ddisk.Free = float64(devicesStats.Free / 1024 / 1024)
				disks = append(disks, ddisk)
			}

			//save in Json
			system_info.Disk = disks
			c.JSON(http.StatusOK, gin.H{"data": system_info})
		})

		//Run Server
		r.Run("0.0.0.0:2221")
	}
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}
}

// getCpuInfo is a function for get one second percent of cpu usage
func getCpuInfo() []float64 {
	_, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
	}

	//CPU utilization
	percent, _ := cpu.Percent(time.Second, false)
	return percent

}
