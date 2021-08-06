package main

type PlayerMove struct {
	id 			uint64
	pos			Vector2
	speed 		uint32
	direction 	uint32
	nextpos 	Vector2
	msgMove 	bool
}
