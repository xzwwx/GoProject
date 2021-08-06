package main

import (
	"usercmd"
)

// byte:cmd(2)|frame(4)|movesNum(4)|moves(...)|
func msgSceneToBytes(cmd uint16, msg *usercmd.MsgScene, buf []byte) int {
	var pos = 0
	//var ismove bool
	buf[0] = byte(cmd)
	buf[1] = byte(cmd>>8)
	pos = 2
	pos = PutUvarint(buf, pos, msg.Frame)

	if len(msg.Moves) > 0 {
		pos = PutUvarint(buf, pos, uint32(len(msg.Moves)))
		for _, player := range msg.Moves {
			pos = PutUvarint64(buf, pos, player.Id)
			pos = PutInt16(buf, pos, int16(player.X))
			pos = PutInt16(buf, pos, int16(player.Y))
			pos = PutInt16(buf, pos, int16(player.Nx))
			pos = PutInt16(buf, pos, int16(player.Ny))

		}
	}else{
		buf[pos] = 0
		pos ++
	}
	return pos
}

func PutInt16(b []byte, i int, v int16) int {
	b[i] = byte(v>>8)
	b[i+1] = byte(v)
	return i + 2
}

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, i int, x uint32) int {
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}
func PutUvarint64(buf []byte, i int, x uint64) int {
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

