# redis 
1. 脑裂问题
2. cluster 扩容
3. 注意事项
	1. key不要过大
	2. 设置缓存时间
	3. 

4. 底层数据结构
	- 简单动态字符串 - sds
	- 压缩列表 - ZipList
	- 快表 - QuickList
	- 字典/哈希表 - Dict
	- 整数集 - IntSet
	- 跳表 - ZSkipList
5.  SDS 不是以空字符串来判断是否结束，而是以 len 属性表示的长度来判断字符串是否结束。