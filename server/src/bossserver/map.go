package main

import "sync"

type Map struct {
	gamemap 	[][]int32
	mapPool		[]map[uint32]interface{}
	needMutex 	bool
	mutexPool 	sync.RWMutex

}

func GenerateRandMap() *map[uint32]*Obstacle {
	obstacle := make(map[uint32]*Obstacle)
	for i := 0; i < 30; i++ {
		o := &Obstacle{
			Id:    uint32(i),
			pos : Pos{0,0},
			OType: 1,
		}
		obstacle[uint32(i)] = o
	}
	return &obstacle

}
