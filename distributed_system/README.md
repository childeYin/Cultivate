# Distrubuted System

* 全局ID怎么生成(参考Vesta)
    - 持久型(表自增字段，Sequence生成)
    - 时间型(机器号，业务号，时间，单节点内自增ID)
    - 完成定期校对模式的一个关键
* log如何追踪
* 什么是分布式事务,为什么会有分布式事务
* 分布式事务在做的时候是考虑数据的一致性还是考虑可用性较多，如何解决数据一致性
    - 在互联网领域的绝大多数的场景中，都需要牺牲强一致性来换取系统的高可用性，系统往往只需要保证“最终一致性”，只要这个最终时间是在用户可以接受的范围内即可。
    - ACID(酸)，BASE(碱) 酸碱平衡
    - 最终一致的模式(查询模式，补偿模式，异步确保模式，定期校对模式，可靠消息模式，缓存一致性模式)
* [最终一致性 + 事务补偿](https://qinnnyul.github.io/2018/09/01/distributed-tx-solutions/)
* 分布式事务的算法 raft
* [详解分布式协调服务 ZooKeeper](https://draveness.me/zookeeper-chubby)
* [分布式事务的实现原理](https://draveness.me/distributed-transaction-principle)


    
    
    