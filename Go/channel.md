# Channel

##### 1.CSP 模型

* Communicating sequential processes 通信顺序进程
* Do not communicate by sharing memory; instead, share memory by communicating
* Channel本质上还是一个队列，遵循FIFO（First In-First Out）原则，

##### 2. 基本结构

```

type hchan struct {
	qcount   uint           // total data in the queue 队列总得元素个数
	dataqsiz uint           // size of the circular queue 
	buf      unsafe.Pointer // points to an array of dataqsiz elements  环形队列
	elemsize uint16
	closed   uint32 //是否已经关闭
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
}

type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	parent   *sudog // semaRoot binary tree
	waitlink *sudog // g.waiting list or semaRoot
	waittail *sudog // semaRoot
	c        *hchan // channel
}

```

##### 3. 实现channel

* chan1 := make(chan type, cap) //cap可以没有

* 接收消息 val := <- chan1  val,ok := <- chan1 

* 写入消息 chan1 <- type Val


