# Go

- [exerciseGo 练习 - Go程序设计语言](https://github.com/childeYin/exerciseGo.git)

1. func 首字母大写 外部可引用，否则外部不可用，只能内部调用
2. 声明变量，赋值，如果只声明了，未赋值则会自动采用零值机制。记int 0, string "", bool false
3. 指针 局部变量的地址是安全的
4. go GC mark && sweep, 
    - mark： mark the object's status
    - sweep: if object's status is unreached, they will be recycled
    - GCPercent [CPU & Memory]
    - MaxHeap []
5. go %取模，符号和被除数保持一致
6. Go 语言中的字符串是一个只读的字节数组切片
     
             type StringHeader struct {
                Data uintptr
                Len  int
            }  
            struct String {
                    byte*   str;
                    intgo   len;
            };
            
     * 字符串作为只读的类型，我们并不会直接向字符串直接追加元素改变其本身的内存空间，所有追加的操作都是通过拷贝来完成的。
     * 新的字符串其实就是一片新的内存空间，与原来的字符串没有任何关联。指针地址并不会改
     * gc的暂停时间小于1ms，指阻塞程序小于1ms
     * 字符串类型设置为不可变是由一些好处的。 我想到的有，可以更好的重用字符串。 还有多线程操作字符串不用加锁。--南哥

7. Array连续的内存空间，每个元素有唯一一个索引(或者叫下标)来访问。   
8. slice底层，相同类型元素的可变长度的序列
    
        struct	Slice
        {				// must not move anything
            byte*	array;		// actual data
            uintgo	len;		// number of elements
            uintgo	cap;		// allocated number of elements
        };
9. slice不指定长度，slice 只能和nil做比较 其他的类型都不可以直接比较，slice的零值是nil, 值为
nil的slice 没有对应的底层数组。
10. slice 是否为空用len(slice) == 0
11. slice 的append处理
    - 检测slice容量是否够大，容的下新的元素，
    
         - 容的下： 引用原始的底层数组，将新添加的放到新位置
         - 容不下： 创建一个足够的新的底层数组，将append的值都分别放到指定位置

