# groutine源码解析

```
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

type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo

	_panic       *_panic // innermost panic - offset known to liblink
	_defer       *_defer // innermost defer
	m            *m      // current m; offset known to arm liblink
	sched        gobuf
    atomicstatus uint32  // g的状态Gidle,Grunnable,Grunning,Gsyscall,Gwaiting,Gdead
    schedlink    guintptr // 下一个g, g链表

	preempt       bool // preemption signal, duplicates stackguard0 = stackpreempt

    lockedm        muintptr // 锁定的M,g中断恢复指定M执行
	gopc           uintptr         // pc of go statement that created this goroutine
	startpc        uintptr         // pc of goroutine function
}

type m struct {
	g0      *g     // goroutine with scheduling stack

	curg          *g       // current running goroutine
	p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	id            int64
	spinning      bool // m is out of work and is actively looking for work
	
	park          note
	alllink       *m // on allm
	schedlink     muintptr

	freelink      *m // on sched.freem

}

type p struct {
	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr

	// runnext, if non-nil, is a runnable G that was ready'd by
	// the current G and should be run next instead of what's in
	// runq if there's time remaining in the running G's time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready'd
	// goroutines to the end of the run queue.
	runnext guintptr

	// Available G's (status == Gdead)
	gFree struct {
		gList
		n int32
	}


	// Per-P GC state
	gcBgMarkWorker       guintptr // (atomic)

	// gcw is this P's GC work buffer cache. The work buffer is
	// filled by write barriers, drained by mutator assists, and
	// disposed on certain GC state transitions.
	gcw gcWork

}

type schedt struct {

	lock mutex

	// When increasing nmidle, nmidlelocked, nmsys, or nmfreed, be
	// sure to call checkdead().

	midle        muintptr // idle m's waiting for work
	nmidle       int32    // number of idle m's waiting for work
	nmidlelocked int32    // number of locked m's waiting for work
	mnext        int64    // number of m's that have been created and next M ID
	maxmcount    int32    // maximum number of m's allowed (or die) 10000
	nmsys        int32    // number of system m's not counted for deadlock
	nmfreed      int64    // cumulative number of freed m's

	pidle      puintptr // idle p's
	npidle     uint32
	nmspinning uint32 // See "Worker thread parking/unparking" comment in proc.go.

	// Global runnable queue.
	runq     gQueue
	runqsize int32

	// Global cache of dead G's.
	gFree struct {
		lock    mutex
		stack   gList // Gs with stacks
		noStack gList // Gs without stacks
		n       int32
	}


	// freem is the list of m's waiting to be freed when their
	// m.exited is set. Linked through m.freelink.
	freem *m

}

```
##### 1. newproc
	
	```
	func newproc(siz int32, fn *funcval) {
		argp := add(unsafe.Pointer(&fn), sys.PtrSize)
		gp := getg()
		pc := getcallerpc()
		systemstack(func() {
			#***真正的执行者***
			newg := newproc1(fn, argp, siz, gp, pc) 
			_p_ := getg().m.p.ptr()
			runqput(_p_, newg, true)
	
			if mainStarted {
				wakep()
			}
		})
	}
	
	```


##### 2. newproc1
	
```

func newproc1(fn *funcval, argp unsafe.Pointer, narg int32, callergp *g, callerpc uintptr) *g {
	_g_ := getg()

	if fn == nil {
		_g_.m.throwing = -1 // do not dump full stacks
		throw("go of nil func value")
	}
	#***g.m.locks++ ***
	acquirem() // disable preemption because it can be holding p in a local var
	siz := narg
	siz = (siz + 7) &^ 7

	// We could allocate a larger initial stack if necessary.
	// Not worth it: this is almost always an error.
	// 4*sizeof(uintreg): extra space added below
	// sizeof(uintreg): caller's LR (arm) or return address (x86, in gostartcall).
	#***groutine 限制大小***
	if siz >= _StackMin-4*sys.RegSize-sys.RegSize {
		throw("newproc: function arguments too large for new goroutine")
	}

	_p_ := _g_.m.p.ptr()
	newg := gfget(_p_)
	if newg == nil {
		newg = malg(_StackMin)
		casgstatus(newg, _Gidle, _Gdead)
		allgadd(newg) // publishes with a g->status of Gdead so GC scanner doesn't look at uninitialized stack.
	}
	if newg.stack.hi == 0 {
		throw("newproc1: newg missing stack")
	}

	if readgstatus(newg) != _Gdead {
		throw("newproc1: new g is not Gdead")
	}

	totalSize := 4*sys.RegSize + uintptr(siz) + sys.MinFrameSize // extra space in case of reads slightly beyond frame
	totalSize += -totalSize & (sys.SpAlign - 1)                  // align to spAlign
	sp := newg.stack.hi - totalSize
	spArg := sp
	if usesLR {
		// caller's LR
		*(*uintptr)(unsafe.Pointer(sp)) = 0
		prepGoExitFrame(sp)
		spArg += sys.MinFrameSize
	}
	if narg > 0 {
		memmove(unsafe.Pointer(spArg), argp, uintptr(narg))
		// This is a stack-to-stack copy. If write barriers
		// are enabled and the source stack is grey (the
		// destination is always black), then perform a
		// barrier copy. We do this *after* the memmove
		// because the destination stack may have garbage on
		// it.
		if writeBarrier.needed && !_g_.m.curg.gcscandone {
			f := findfunc(fn.fn)
			stkmap := (*stackmap)(funcdata(f, _FUNCDATA_ArgsPointerMaps))
			if stkmap.nbit > 0 {
				// We're in the prologue, so it's always stack map index 0.
				bv := stackmapdata(stkmap, 0)
				bulkBarrierBitmap(spArg, spArg, uintptr(bv.n)*sys.PtrSize, 0, bv.bytedata)
			}
		}
	}

	memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
	newg.sched.sp = sp
	newg.stktopsp = sp
	newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
	newg.sched.g = guintptr(unsafe.Pointer(newg))
	gostartcallfn(&newg.sched, fn)
	newg.gopc = callerpc
	newg.ancestors = saveAncestors(callergp)
	newg.startpc = fn.fn
	if _g_.m.curg != nil {
		newg.labels = _g_.m.curg.labels
	}
	if isSystemGoroutine(newg, false) {
		atomic.Xadd(&sched.ngsys, +1)
	}
	casgstatus(newg, _Gdead, _Grunnable)

	if _p_.goidcache == _p_.goidcacheend {
		// Sched.goidgen is the last allocated id,
		// this batch must be [sched.goidgen+1, sched.goidgen+GoidCacheBatch].
		// At startup sched.goidgen=0, so main goroutine receives goid=1.
		_p_.goidcache = atomic.Xadd64(&sched.goidgen, _GoidCacheBatch)
		_p_.goidcache -= _GoidCacheBatch - 1
		_p_.goidcacheend = _p_.goidcache + _GoidCacheBatch
	}
	newg.goid = int64(_p_.goidcache)
	_p_.goidcache++
	if raceenabled {
		newg.racectx = racegostart(callerpc)
	}
	if trace.enabled {
		traceGoCreate(newg, newg.startpc)
	}
	releasem(_g_.m)

	return newg
}
```
##### 3. 详细说明

1.GPM (P最大为256，M最大为10000，)
	* G: 表示 Goroutine，每个 Goroutine 对应一个 G 结构体，G 存储 Goroutine 的运行堆栈、状态以及任务函数，可重用。G 并非执行体，每个 G 需要绑定到 P 才能被调度执行。

	* P: Processor，表示逻辑处理器， 对 G 来说，P 相当于 CPU 核，G 只有绑定到 P(在 P 的 local runq 中)才能被调度。对 M 来说，P 提供了相关的执行环境(Context)，如内存分配状态(mcache)，任务队列(G)等，P 的数量决定了系统内最大可并行的 G 的数量（前提：物理 CPU 核数 >= P 的数量），P 的数量由用户设置的 GOMAXPROCS 决定，但是不论 GOMAXPROCS 设置为多大，P 的数量最大为 256。

	* M: Machine，OS 线程抽象，代表着真正执行计算的资源，在绑定有效的 P 后，进入 schedule 循环；而 schedule 循环的机制大致是从 Global 队列、P 的 Local 队列以及 wait 队列中获取 G，切换到 G 的执行栈上并执行 G 的函数，调用 goexit 做清理工作并回到 M，如此反复。M 并不保留 G 状态，这是 G 可以跨 M 调度的基础，M 的数量是不定的，由 Go Runtime 调整，为了防止创建过多 OS 线程导致系统调度不过来，目前默认最大限制为 10000 个。

2. work-stealing
	
	* 每个 P 维护一个 G 的本地队列；
	* 当一个 G 被创建出来，或者变为可执行状态时，就把他放到 P 的可执行队列中；
	* 当一个 G 在 M 里执行结束后，P 会从队列中把该 G 取出；如果此时 P 的队列为空，即没有其他 G 可以执行， M 就随机选择另外一个 P，从其可执行的 G 队列中取走一半。


初始化G:

	graph TB;
	    A["go func() {}  我们写的goroutine"] --> B
	    B["newproc 获取func和参数"] -- "切换到g0,使用g0栈空间" --> C
	    C[newproc1] --> D[gfget 从当前P获取空闲G]
	    D --> E{"P空&全局不空"}
	    E -- Y --> F[全局移32个到P本地]
	    E -- N --> G
	    F --> G[本地取出空闲G,初始化栈空间]
	    G --> H{"获取空闲G成功"}
	    H -- N --> I[创建G,初始化栈空间, 加入全局G数组]
	    H -- Y --> J
	    I --> J["参数复制到栈,清除堆,pc:func,sp:goexit1"]
	    J --> K["状态设为runable,设置goid"]
	    K -- runqput 加入当前P的runable队列 --> L[用g替换runnext]
	    L --> M{ 本地runable队列满 }
	    M -- Y --> N[本地runbale队列移一半到全局] 
	    M -- N --> O[加入本地runable队列]
	    N --> M
	    O --> P{有空闲P&没有自旋的M }
	    P -- Y --> Q["wakep()"]

初始化P:

	通过启动时候的schedinit调用procresize生成对应个数的P。因为可以通过runtime.GOMAXPROCS来动态修改P的个数，所以在procresize中会对P数组进行调整，或新增P或减少P。被减少的P会将自身的runable、runnext、gfee移到全局去。

	如果当前P不在多余的P中，则状态为running
	如果当前P在多余的P中，则将当前M和P解绑，再将M和P数组的第一P绑定，并设为running
	除了当前P外；所有P都设为idle，如果P中没有runnable,则将P加入全局空闲P,否则获取全局空闲M和P绑定。
    
初始化M:
	graph TB;
		A["有空闲P&没有自旋的M"] --> B["wakep()"]
		B -- startm --> B2[全局获取空闲P] 
		B2--> C[全局获取空闲M]
		C --> D{获取成功}
		D -- Y --> E[M和P绑定]
		D -- N --> F["创建M"]
		F --> E
		E --> G[唤醒M]
		G --> H[mstart / mstart1]
		H --> I[schedule]
		I --> J[在P本地或全局获取runbale G]
		J --> K{获取成功}
		K -- Y --> L
		K -- N --> M["反复在本地、全局、网络、其他P中获取runable G"]
		M --> N{获取成功}
		N -- Y --> L["G和M相互绑定，G设为running"]
		L --> O["汇编执行G的pc:func。执行完RET弹出sp:goexit1"]
		O --> P[goexit1 / goexit0]
		P --> Q["将G状态设为dead,清除G的各种信息，G和M相互解绑"]
		Q --> R["G放入本地空闲链表。如果本地空闲个数大于64个，则移一半到全局去"]
		R --> I
		N -- N --> S["MP解绑，P加入全局空闲P,M加入全局空闲M"]
		S --> T[M进入睡眠]

启动小结:
	
	graph TB;
		A["创建G0、M0, g0.m = m0"] --> C
	    C[初始化命令行和OS]-->D
	    D["schedinit:设置M最大数量、P个数、栈和内存初始化"] --> E
	    E["newproc:为main.main创建一个主goroutine"] --> F
	    F["mstart:运行主goroutine --> 运行main.main"]
      

