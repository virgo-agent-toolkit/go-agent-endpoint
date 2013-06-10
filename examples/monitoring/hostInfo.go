package main

import (
	"encoding/json"
	"net"
)

type systemInfoResponse struct {
	SysInfo siSysInfo `json:"sysinfo"`
	NetIFs  []siNetIF `json:"netifs"`
	CPUs    []siCPU   `json:"cpus"`
}

type siSysInfo struct {
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	Vendor         string `json:"vendor"`
	VendorCodeName string `json:"vendor_code_name"`
	VendorVersion  string `json:"vendor_version"`
	VendorName     string `json:"vendor_name"`
	PatchLevel     string `json:"patch_level"`
	Machine        string `json:"machine"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

type siNetIF struct {
	Info  siNetIFInfo  `json:"info"`
	Usage siNetIFUsage `json:"usage"`
}

type siCPU struct {
	Info siCPUInfo `json:"info"`
	Data siCPUData `json:"data"`
}

type siNetIFInfo struct {
	Metric      int              `json:"metric"`
	MTU         int              `json:"mtu"`
	Flags       int              `json:"flags"`
	Type        string           `json:"type"`
	Name        string           `json:"name"`
	Broadcast   *net.IPAddr      `json:"broadcast"`
	Address     *net.IPAddr      `json:"address"`
	Address6    *net.IPAddr      `json:"address6"`
	HWAddr      net.HardwareAddr `json:"hwaddr"`
	Destination *net.IPAddr      `json:"destination"`
	Netmask     net.IPMask       `json:"netmask"`
}

type _helper_siNetIFInfo struct {
	Metric      int    `json:"metric"`
	MTU         int    `json:"mtu"`
	Flags       int    `json:"flags"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Broadcast   string `json:"broadcast"`
	Address     string `json:"address"`
	Address6    string `json:"address6"`
	HWAddr      string `json:"hwaddr"`
	Destination string `json:"destination"`
	Netmask     string `json:"netmask"`
}

func (i *siNetIFInfo) UnmarshalJSON(data []byte) (err error) {
	var tmp _helper_siNetIFInfo
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return
	}

	i.Metric = tmp.Metric
	i.MTU = tmp.MTU
	i.Flags = tmp.Flags
	i.Type = tmp.Type
	i.Name = tmp.Name

	i.Broadcast, err = net.ResolveIPAddr("ip", tmp.Broadcast)
	if err != nil {
		return
	}
	i.Address, err = net.ResolveIPAddr("ip", tmp.Address)
	if err != nil {
		return
	}
	i.Address6, err = net.ResolveIPAddr("ip", tmp.Address6)
	if err != nil {
		return
	}
	i.HWAddr, err = net.ParseMAC(tmp.HWAddr)
	if err != nil {
		return
	}
	i.Destination, err = net.ResolveIPAddr("ip", tmp.Destination)
	if err != nil {
		return
	}
	i.Netmask = net.IPMask(net.ParseIP(tmp.Netmask))
	return
}

type siNetIFUsage struct {
	TxDropped    int `json:"tx_dropped"`
	TxErrors     int `json:"tx_errors"`
	TxOverruns   int `json:"tx_overruns"`
	RxErrors     int `json:"rx_errors"`
	TxCarrier    int `json:"tx_carrier"`
	TxCollisions int `json:"tx_collisions"`
	RxBytes      int `json:"rx_bytes"`
	TxPackets    int `json:"tx_packets"`
	RxPackets    int `json:"rx_packets"`
	RxOverruns   int `json:"rx_overruns"`
	TxBytes      int `json:"tx_bytes"`
	RxDropped    int `json:"rx_dropped"`
	RxFrame      int `json:"rx_frame"`
}

type siCPUInfo struct {
	CacheSize      int    `json:"cache_size"`
	CoresPerSocket int    `json:"cores_per_socket"`
	Vendor         string `json:"vendor"`
	Model          string `json:"model"`
	TotalSockets   int    `json:"total_sockets"`
	TotalCores     int    `json:"total_cores"`
	MHz            int    `json:"mhz"`
}

type siCPUData struct {
	SoftIrq int `json:"soft_irq"`
	User    int `json:"user"`
	Idle    int `json:"idle"`
	Nice    int `json:"nice"`
	Total   int `json:"total"`
	Stolen  int `json:"stolen"`
	Irq     int `json:"irq"`
	Sys     int `json:"sys"`
	Wait    int `json:"wait"`
}
