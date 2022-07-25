# 内存池

1. Zend Memory Manager(ZMM) ，如果ZMM内存可用，直接分配给php程序，如果不够用再有ZMM像系统申请
	- 减少系统调用次数，并优化内存空间的使用效率

2. 内存池提供了
	- 规划 allocate
	- 使用 access
	- 归还 free
3. ZMM思想：
	向系统申请大块内存，再按照固定的几种规格分割成较小的内存块，由内存池统一管理。当调用方申请内存的时候，从池子中匹配已经预分配的合适的大小的内存块返回。

4. 
	- Huge(chunk): 申请内存大于2M，直接调用系统分配，分配若干个chunk
	- Large(page): 申请内存大于3K(3/4 page_size)，小于2044K(511 page_size)，分配若干个page
	- Small(slot): 申请内存小于等于3K(3/4 page_size)

5. 