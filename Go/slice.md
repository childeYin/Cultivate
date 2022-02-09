#slice

##### 1. cmd/compile/internal/types/type.go

	// NewSlice returns the slice Type with element type elem.
	func NewSlice(elem *Type) *Type {
		if t := elem.Cache.slice; t != nil {
			if t.Elem() != elem {
				Fatalf("elem mismatch")
			}
			return t
		}

		t := New(TSLICE)
		t.Extra = Slice{Elem: elem}
		elem.Cache.slice = t
		return t
	}

##### 2. slice 结构体 

* reflect/value.go

	// SliceHeader is the runtime representation of a slice.
	// It cannot be used safely or portably and its representation may
	// change in a later release.
	// Moreover, the Data field is not sufficient to guarantee the data
	// it references will not be garbage collected, so programs must keep
	// a separate, correctly typed pointer to the underlying data.
	type SliceHeader struct {
		Data uintptr
		Len  int
		Cap  int
	}

#####3. 初始化slice 

	* runtime/slice.go
	* arr[0:3],slice[0:3]
	* []int{1,2,3}
	* make([]int,10)

	func makeslice(et *_type, len, cap int) unsafe.Pointer {
		mem, overflow := math.MulUintptr(et.size, uintptr(cap))
		if overflow || mem > maxAlloc || len < 0 || len > cap {
			// NOTE: Produce a 'len out of range' error instead of a
			// 'cap out of range' error when someone does make([]T, bignumber).
			// 'cap out of range' is true too, but since the cap is only being
			// supplied implicitly, saying len is clearer.
			// See golang.org/issue/4085.
			mem, overflow := math.MulUintptr(et.size, uintptr(len))
			if overflow || mem > maxAlloc || len < 0 {
				panicmakeslicelen() //	panic(errorString("makeslice: len out of range"))

			}
			panicmakeslicecap() //	panic(errorString("makeslice: cap out of range"))

		}

		return mallocgc(mem, et, true)
	}

	tips:
		需要注意的是使用下标初始化切片不会拷贝原数组或者原切片中的数据，它只会创建一个指向原数组的切片结构体，所以修改新切片的数据也会修改原切片。

		slice1 := []int{1, 2, 3, 4, 5}
		slice2 := slice1[0:2]
		slice2 = append(slice2, 0, 5)
		fmt.Println(slice2) // 1,2,0,5
		fmt.Println(slice1) // 1,2,0,5,5

##### 4. slice 内存占用
	
	* 内存空间=切片中元素大小×切片容量

	* tips
		* 内存空间的大小发生了溢出
		* 申请的内存大于最大可分配的内存
		* 传入的长度小于 0 或者长度大于容量

##### 5. slice 追加元素

* compile/internal/gc/ssa.go 

	func (s *state) append(n *Node, inplace bool) *ssa.Value {
		// If inplace is false, process as expression "append(s, e1, e2, e3)":
		// 不覆盖原有的
 		ptr, len, cap := s
		 newlen := len + 3
		 if newlen > cap {
		     ptr, len, cap = growslice(s, newlen)
		     newlen = len + 3  //recalculate to avoid a spill
		 }
		 // with write barriers, if needed:
		 *(ptr+len) = e1
		 *(ptr+len+1) = e2
		 *(ptr+len+2) = e3
		 return makeslice(ptr, newlen, cap)
		
		//
		// If inplace is true, process as statement "s = append(s, e1, e2, e3)":
		// 覆盖原有的
 		a := &s
		 ptr, len, cap := s
		 newlen := len + 3
		 if uint(newlen) > uint(cap) {
		    newptr, len, newcap = growslice(ptr, len, cap, newlen)
		    vardef(a)        //if necessary, advise liveness we are writing a new a
		    *a.cap = newcap  //write before ptr to avoid a spill
		    *a.ptr = newptr  //with write barrier
		 }
		 newlen = len + 3  //recalculate to avoid a spill
		 *a.len = newlen
		//  with write barriers, if needed:
		 *(ptr+len) = e1
		 *(ptr+len+1) = e2
		 *(ptr+len+2) = e3		

	}

* runtime/slice.go

	func growslice(et *_type, old slice, cap int) slice {
		newcap := old.cap
		doublecap := newcap + newcap
		//确定切片容量
		if cap > doublecap {
			newcap = cap
		} else {
			if old.len < 1024 {
				newcap = doublecap
			} else {
				// Check 0 < newcap to detect overflow
				// and prevent an infinite loop.
				for 0 < newcap && newcap < cap {
					newcap += newcap / 4
				}
				// Set newcap to the requested cap when
				// the newcap calculation overflowed.
				if newcap <= 0 {
					newcap = cap
				}
			}
		}
		//内存对齐
		var overflow bool
		var lenmem, newlenmem, capmem uintptr
		// Specialize for common values of et.size.
		// For 1 we don't need any division/multiplication.
		// For sys.PtrSize, compiler will optimize division/multiplication into a shift by a constant.
		// For powers of 2, use a variable shift.
		switch {
		case et.size == 1:
			lenmem = uintptr(old.len)
			newlenmem = uintptr(cap)
			capmem = roundupsize(uintptr(newcap))
			overflow = uintptr(newcap) > maxAlloc
			newcap = int(capmem)
		case et.size == sys.PtrSize:
			lenmem = uintptr(old.len) * sys.PtrSize
			newlenmem = uintptr(cap) * sys.PtrSize
			capmem = roundupsize(uintptr(newcap) * sys.PtrSize)
			overflow = uintptr(newcap) > maxAlloc/sys.PtrSize
			newcap = int(capmem / sys.PtrSize)
			。。。。
		}

		。。。。

		//返回了新的切片
		return slice{p, old.len, newcap}
	}

* 运行时根据切片的当前容量选择不同的策略进行扩容：

	* 如果期望容量大于当前容量的两倍就会使用期望容量
	* 如果当前切片的长度小于 1024 就会将容量翻倍
	* 如果当前切片的长度大于 1024 就会每次至少增加 25% 的容量，直到新容量大于期望容量




##### 6.参考文章：

	1. [Go slice扩容深度分析](https://juejin.cn/post/6844903812331732999)
	2. [3.2 切片](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array-and-slice/)