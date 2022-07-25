# Go的内存分配

1. 内存分配：栈区（stack）堆区（heap）
	- 函数调用的参数、返回值以及局部变量大都会被分配到栈上，这部分内存会由编译器进行管理

##### 2. 内存管理
	* 用户程序（Mutator）
	* 分配器（Asllocator）
	* 收集器（Collector）
	![图片来源见参考](https://img.draveness.me/2020-02-29-15829868066411-mutator-allocator-collector.png)

##### 3. 分配方法

	* 线性分配器（Sequential Allocator Bump Allocator）
		* 快速移动
		* 内存不可重复利用
		* 标记压缩（Mark-Compact）、复制回收（Copying GC）和分代回收（Generational GC）

	* 空闲链表分配器(Free-list Allocator)
		* 可重用内存
		* 分配原则
			- 首次适应
			- 循环首次适应
			- 最优适应
			- 隔离适应

		![隔离适应](https://img.draveness.me/2020-02-29-15829868066452-segregated-list.png)

##### 4. 分级分配
	
	* 核心理念是使用多级缓存将对象根据大小分类，并按照类别实施不同的分配策略。
	* 对象大小	
		* 微对象（0，16B）
		* 小对象 [16B,32k]
		* 大对象（32k, +∞）
	* 多级缓存
		* 线程缓存（thread cache）
		* 中心缓存(central cache)
		* 页缓存(page heap)

	![多级缓存](https://img.draveness.me/2020-02-29-15829868066457-multi-level-cache.png)

##### 5. 虚拟内存
	* 线性内存
		* span (512M)
		* bitmap (16G)
		* arena (512G  page 8K)
	* 稀疏内存
		* heapArena （64M）
		```
		type heapArena struct {
			bitmap       [heapArenaBitmapBytes]byte
			spans        [pagesPerArena]*mspan
			pageInUse    [pagesPerArena / 8]uint8
			pageMarks    [pagesPerArena / 8]uint8
			pageSpecials [pagesPerArena / 8]uint8
			checkmarks   *checkmarksMap
			zeroedBase   uintptr
		}
		```
##### 6.地址空间
	* None 内存没有被保留或者映射
	* Reserved 运行时持有该地址空间，但是访问内存报错
	* Prepared 内存被保留，可以快速转换为Ready
	* Ready 可以被安全访问

![memory region states & transitions](https://img.draveness.me/2020-02-29-15829868066474-memory-regions-states-and-transitions.png)


#####7. 内存管理组件
	
	* 内存管理单元
	* 线程缓存
	* 中心缓存
	* 页堆

![Go memory layout ](https://img.draveness.me/2020-02-29-15829868066479-go-memory-layout.png)

	*每一个处理器都会分配一个线程缓存 runtime.mcache 用于处理微对象和小对象的分配，它们会持有内存管理单元 runtime.mspan

	* runtime.mspan 不存在空闲对象的时候，从 runtime.mheap 持有的 134 个中心缓存 runtime.mcentral 中获取新的内存单元,中心缓存属于全局的堆结构体 runtime.mheap，它会从操作系统中申请内存。

	* amd64 Linux, runtime.mheap 会持有 4,194,304 runtime.heapArena，每个 runtime.heapArena 都会管理 64MB 的内存，单个 Go 语言程序的内存上限也就是 256TB。


	* runtime.mspan 是 Go 语言内存管理的基本单元
	```
	type mspan struct {
		next *mspan
		prev *mspan
		...
		startAddr uintptr // 起始地址
		npages    uintptr // 页数
		freeindex uintptr

		allocBits  *gcBits
		gcmarkBits *gcBits
		allocCache uint64
	}
	```


