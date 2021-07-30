package common

import "net/http"

type RetMsg []byte

func (p RetMsg) Size() int {
	return len([]byte(p))
}

func (p RetMsg) MarshalTo(data []byte) (int, error) {
	copy(data, []byte(p))
	return 0, nil
}

func (p RetMsg) UnMarshal(data []byte) error {
	return nil
}

func (p RetMsg) Reset() {

}

//------Response Writer
type ResWrite struct {
	Buf []byte
}

func (this *ResWrite) Header() http.Header {
	return nil
}

func (this *ResWrite) Write(data []byte) (int, error) {
	this.Buf = data
	return 0, nil
}

func (this *ResWrite) WriteHeader (statusCode int) {

}

type HttpBody struct {
	Buf []byte
}

func (this *HttpBody)Read(p []byte) (n int, err error) {
	copy(p, this.Buf)
	return len(this.Buf), nil
}

func (this *HttpBody) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (this *HttpBody) Close() error {
	return nil
}


// Request room
type ReqRoom struct {
	UserId 	uint64
	IsNew 	bool
	//	others
}

// Reply get room
type RetRoom struct {
	ServerId	uint16
	Address 	string
	RoomId		uint32
	EndTime		uint32
	NewSync		bool
}
