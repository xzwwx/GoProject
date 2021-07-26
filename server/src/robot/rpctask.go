package main
import (
	glog "glog-master"
	"gonet"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const(
	cmd_verify_time = 30
)
type RpcTask struct{
	closed 		int32
	virified	bool
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
		virified:	 false,
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

func (this *RpcTask)IsClosed() bool{
	return atomic.LoadInt32(&this.closed)!=0
}

func (this *RpcTask)Terminate(){
	this.Close()
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
		case <- this.signal:
			for {
				this.sendMutex.Lock()
				if this.sendBuff.RdReady(){
					tmpByte.Append(this.sendBuff.RdBuf()[:this.sendBuff.RdSize()]...)
					this.sendBuff.Reset()
				}
				this.sendMutex.Unlock()
			}
			if !tmpByte
		}
	}
}