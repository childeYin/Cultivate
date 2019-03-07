# Distrubuted System

* 全局ID怎么生成
* log如何追踪
* 分布式事务
* 理解什么是事务
* 为什么会有分布式事务(微服务)
* 分布式事务在做的时候是考虑数据的一致性还是考虑可用性较多
    - 在互联网领域的绝大多数的场景中，都需要牺牲强一致性来换取系统的高可用性，系统往往只需要保证“最终一致性”，只要这个最终时间是在用户可以接受的范围内即可。
* [最终一致性 + 事务补偿](https://qinnnyul.github.io/2018/09/01/distributed-tx-solutions/)
* 分布式事务的算法 raft
* [详解分布式协调服务 ZooKeeper](https://draveness.me/zookeeper-chubby)
* [分布式事务的实现原理](https://draveness.me/distributed-transaction-principle)