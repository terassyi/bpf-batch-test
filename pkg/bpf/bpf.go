package bpf

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cilium/ebpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang XdpProg ../../bpf/main.c -g -- -I../../bpf/include

const (
	MAP_ID_FILE string = "map_id"
	MAX_ENTRIES int = 1024 * 1024 * 100
)

func LoadAndGetMap() (*ebpf.Map, error) {
	obj := XdpProgObjects{}
	if err := LoadXdpProgObjects(&obj, &ebpf.CollectionOptions{}); err != nil {
		return nil, err
	}

	file, err := os.Create(MAP_ID_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := obj.TestMap.Info()
	if err != nil {
		return nil, err
	}
	i, ok := info.ID()
	if !ok {
		return nil, fmt.Errorf("failed to get map id")
	}
	file.Write([]byte(strconv.Itoa(int(i))))

	fmt.Println("map id is ", i)
	return obj.TestMap, nil
}

func UnLoad() error {
	obj := XdpProgObjects{}
	obj.Close()
	return nil
}

func GetMap() (*ebpf.Map, error) {
	file, err := os.Open(MAP_ID_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b := make([]byte, 10)
	i, err := file.Read(b)
	if err != nil {
		fmt.Println("failed to read file")
		return nil, err
	}
	n, err := strconv.Atoi(string(b[:i]))
	if err != nil {
		return nil, err
	}

	fmt.Println("get from id: ", ebpf.MapID(n))
	return ebpf.NewMapFromID(ebpf.MapID(n))
}
