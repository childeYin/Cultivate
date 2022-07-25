# map 

1. 关于哈希碰撞

	哈希查找表一般会存在“碰撞”的问题，就是说不同的 key 被哈希到了同一个 bucket。

	一般有两种应对方法：链表法和开放地址法。

	链表法将一个 bucket 实现成一个链表，落在同一个 bucket 中的 key 都会插入这个链表。(Golang)

	开放地址法则是碰撞发生后，通过一定的规律，在数组的后面挑选“空位”，用来放置新的 key。

	Go 语言中，通过哈希查找表实现 map，用链表法解决哈希冲突。

2. // A header for a Go map.
type hmap struct {
    // 元素个数，调用 len(map) 时，直接返回此值
    count     int
    flags     uint8 // 首先会检查 h.flags 标志，如果发现写标位是 1，直接panic, 标识有其他的协程再操作
    // buckets 的对数 log_2
    B         uint8
    // overflow 的 bucket 近似数
    noverflow uint16
    // 计算 key 的哈希的时候会传入哈希函数
    hash0     uint32
    // 指向 buckets 数组，大小为 2^B
    // 如果元素个数为0，就为 nil
    buckets    unsafe.Pointer
    // 扩容的时候，buckets 长度会是 oldbuckets 的两倍
    oldbuckets unsafe.Pointer
    // 指示扩容进度，小于此地址的 buckets 迁移完成
    nevacuate  uintptr
    extra *mapextra // optional fields
}

type bmap struct {
    tophash [bucketCnt]uint8
}

type bmap struct {
    topbits  [8]uint8
    keys     [8]keytype
    values   [8]valuetype
    pad      uintptr
    overflow uintptr
}

3. map 

map count 存储了map的长度，可以用len(map)查看长度，flags 标识有其他的协程再操作，防止资源竞争，map非线程安全。

赋值流程：（删除，扩容大部分同理，扩容多了移动）
	1. 先判断了map是否为nil
	2. 判断map flags 是否有其他资源再操作写入
	3. 设置map flags
	4. h.buckets 是否有值，没有的话 创建newobject(t.bucket)
	5. 循环 again
		1. 获取bucket
		2. 判断map是否需要growing growWork
		3. 循环 使用hash计算，寻址 放到合适的槽位
		4. 结束的时候再次判断是否有碰撞，flags 置为0
		5. 返回值所在的地址


4. map 扩容触发在 mapassign 中，我们之前注释过了，主要是两点:

是不是已经到了 load factor 的临界点，即元素个数 >= 桶个数 * 6.5，这时候说明大部分的桶可能都快满了，如果插入新元素，有大概率需要挂在 overflow 的桶上。
overflow 的桶是不是太多了，当 bucket 总数 < 2 ^ 15 时，如果 overflow 的 bucket 总数 >= bucket 的总数，那么我们认为 overflow 的桶太多了。当 bucket 总数 >= 2 ^ 15 时，那我们直接和 2 ^ 15 比较，overflow 的 bucket >= 2 ^ 15 时，即认为溢出桶太多了。为啥会导致这种情况呢？是因为我们对 map 一边插入，一边删除，会导致其中很多桶出现空洞，这样使得 bucket 使用率不高，值存储得比较稀疏。在查找时效率会下降。
两种情况官方采用了不同的解决方法:

针对 1，将 B + 1，进而 hmap 的 bucket 数组扩容一倍；
针对 2，通过移动 bucket 内容，使其倾向于紧密排列从而提高 bucket 利用率。