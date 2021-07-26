package usercmd

func encodeVarint(data []byte, offset *int, v uint64){
	for v >= 1 <<7{
		data[*offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		*offset++
	}
	data[*offset] = uint8(v)
	*offset++
}

func encodeVarintSize(x uint64)(n int){
	for{
		n++
		x>>=7;
		if x == 0{
			break
		}
	}
	return n
}
