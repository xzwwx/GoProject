package main
import (
	glog "glog-master"
	"gonet"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const(
	cmd_max_size = 128 * 1024
	cmd_verify_time = 30
	cmd_header_size = 4	// 3 bytes length of instruction + 1 byte flag
)


type IRpcTask interface{
	ParseMsg(data []byte) bool
	OnClose()
}



type RpcTask struct{
	closed 		int32
	verified	bool
	stopedChan	chan struct{}
	recvBuff	*gonet.ByteBuffer
	sendBuff	*gonet.ByteBuffer
	sendMutex	sync.Mutex
	Conn		net.Conn
	Drived		IRpcTask
	signal		chan struct{}

}

func (this *RpcTask) Signal(){
	select {
	case this.signal <- struct{}{}:
	default:
	}
}

func (this *RpcTask) RemoteAddr() string{
	if this.Conn == nil {
		return ""
	}
	return this.Conn.RemoteAddr().String()
}

func (this *RpcTask) LocalAddr() string {
	if this.Conn == nil {
		return ""
	}
	return this.Conn.LocalAddr().String()
}

func (this *RpcTask) Stop() bool {
	if this.IsClosed(){
		glog.Info("[Connection] Close Failed.", this.RemoteAddr())
		return false
	}
	select {
	case this.stopedChan <- struct{}{}:
	default:
		glog.Info("[Conntion] Close Failed.", this.RemoteAddr())
		return false
	}
	return true
}

func (this *RpcTask) Start() {
	if !atomic.CompareAndSwapInt32(&this.closed, -1, 0){
		return
	}
	job := &sync.WaitGroup{}
	job.Add(1)
	go this.sendloop(job)
	go this.recvloop()
	job.Wait()
	glog.Info("[Connection] Received connection.", this.RemoteAddr())
}

func (this *RpcTask)Close() {
	if !atomic.CompareAndSwapInt32(&this.closed, 0,1){
		return
	}
	glog.Info("[Connection] Closed.", this.RemoteAddr())
	this.Conn.Close()
	this.recvBuff.Reset()
	this.sendBuff.Reset()
	close(this.stopedChan)
	this.Drived.Onclose()
}

func NewRpcTask(conn net.Conn) *RpcTask{
	return &RpcTask{
		closed: 	-1,
		verified:	 false,
		Conn: 		conn,
		stopedChan: make(chan struct{}, 1),
		recvBuff: 	gonet.NewByteBuffer(),
		sendBuff: 	gonet.NewByteBuffer(),
		signal: 	make(chan struct{}, 1),
	}
}

func (this * RpcTask) SendBytes(buffer []byte) bool{
	if this.IsClosed(){
		return false
	}
	this.sendMutex.Lock()
	this.sendBuff.Append(buffer...)
	this.sendMutex.Unlock()
	this.Signal()
	return true
}

func (this *RpcTask) readAtLeast(buff *gonet.ByteBuffer, neednum int) error {
	buff.WrGrow(neednum)
	n, err := io.ReadAtLeast(this.Conn, buff.WrBuf(), neednum)
	buff.WrFlip(n)
	return err
}

func (this *RpcTask) Reset() bool {
	if atomic.LoadInt32(&this.closed) != 1 {
		return false
	}
	if !this.IsVerified(){
		return false
	}
	this.closed = -1
	this.verified = false
	this.stopedChan = make(chan struct{})
	glog.Info("[Connection] Reset connection.")
	return true
}

func (this *RpcTask)IsClosed() bool{
	return atomic.LoadInt32(&this.closed)!=0
}

func (this *RpcTask)Terminate(){
	this.Close()
}

func (this *RpcTask)Verify(){
	this.verified = true
}

func (this *RpcTask) IsVerified() bool{
	return this.verified
}

func (this *RpcTask) recvloop() {
	defer func(){
		this.Close()
		if err := recover(); err != nil {
			glog.Error("[Exception] ", err, "\n", string(debug.Stack()))
		}
	}()

	var(
		neednum 	int
		err			error
		totalsize	int
		datasize 	int
		msgbuff		[]byte
	)
	for {
		totalsize = this.recvBuff.RdSize()
		if totalsize <= cmd_header_size {
			neednum = cmd_header_size - totalsize
			err = this.readAtLeast(this.recvBuff, neednum)
			if err != nil {
				glog.Error("[Connection] Receive failed.")
				return
			}
			totalsize = this.recvBuff.RdSize()
		}

		msgbuff = this.recvBuff.RdBuf()
		datasize = int(msgbuff[0]) << 16 | int(msgbuff[1]) << 8 | int(msgbuff[2])
		if datasize > cmd_max_size{
			glog.Error("[Connection] Data oversize.", this.RemoteAddr(), ",", datasize)
			return
		} else if datasize < cmd_header_size {
			glog.Error("[Connection] Data too few.")
			return
		}

		if totalsize < datasize {
			neednum = datasize - totalsize
			err = this.readAtLeast(this.recvBuff, neednum)
			if err != nil {
				glog.Error("[Connection] Receive failed.", this.RemoteAddr(), ",", datasize)
				return
			}

			msgbuff = this.recvBuff.RdBuf()
		}
		this.Drived.ParseMsg(msgbuff[:datasize])
		this.recvBuff.RdFlip(datasize)

	}


}

func (this *RpcTask) sendloop(job *sync.WaitGroup){
	defer func(){
		this.Close()
		if err := recover(); err != nil {
			glog.Error("[Exception] ", err, "\n", string(debug.Stack()))
		}
	}()
	var (
		tmpByte 	= gonet.NewByteBuffer()
		timeout 	= time.NewTimer(time.Second * cmd_verify_time)
		writenum 	int
		err 		error
	)
	defer timeout.Stop()

	job.Done()

	for {
		select {
		case <-this.signal:
			for {
				this.sendMutex.Lock()
				if this.sendBuff.RdReady() {
					tmpByte.Append(this.sendBuff.RdBuf()[:this.sendBuff.RdSize()]...)
					this.sendBuff.Reset()
				}
				this.sendMutex.Unlock()

				if !tmpByte.RdReady() {
					break
				}

				writenum, err = this.Conn.Write(tmpByte.RdBuf()[:tmpByte.RdSize()])
				if err != nil {
					glog.Error("[Connetion] Send failed.", this.RemoteAddr(), ", ", err)
					return
				}
				tmpByte.RdFlip(writenum)
			}
		case <-this.stopedChan:
			return
		case <-timeout.C:
			if !this.IsVerified() {
				glog.Error("[Connection] Verified Timeout.", this.RemoteAddr())
				return
			}
		}
	}
}
