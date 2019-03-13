#  Cache

1. 什么是NoSQL，NewSQL
2. Redis常用数据类型,以及底部存储结构
3. Redis和Memecache MongoDB对比，它为啥高性能 
4. Redis数据落地怎么做，aof和rdb优缺点
5. Redis集群的方案 
6. 分布式锁，为什么需要分布式锁, 有什么实现方式？
    - 数据的一致性，同时成功，同时失败
7. Redis中大key怎么查找
    - bigkeys
8. Redis淘汰算法
    - noeviction：当内存使用达到阈值的时候，所有引起申请内存的命令会报错。
    - allkeys-lru：在主键空间中，优先移除最近未使用的key。
    - volatile-lru：在设置了过期时间的键空间中，优先移除最近未使用的key。
    - allkeys-random：在主键空间中，随机移除某个key。
    - volatile-random：在设置了过期时间的键空间中，随机移除某个key。
    - volatile-ttl：在设置了过期时间的键空间中，具有更早过期时间的key优先移除。
9. Redis, 如果List中数据量不定, 有的时候特别大，采用何种方式优化
    - 阻塞读取