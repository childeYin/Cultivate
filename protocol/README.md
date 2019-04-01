# Protocol

##  HTTP
1. HTTP 请求头，响应头
2. CSRF，XSS如何防御
    - 同源测略，白名单
    - 过滤请求数据
3. HTTP GET vs POST,
4. HTTP Connection: keep-alivee with Content-Length, 
    - Connection: keep-alivee with Content-Length 连用，在复用的时候，可以知道自己接收完了信息
    - Transfer-Encoding: chunked 不需要Content-Length 因为chunk的最后一个包大小是0，必须是0 
5. [RFC 7540 :  HTTP/2 ](https://tools.ietf.org/html/rfc7540)
6. Session vs Cookie
7. [url输入之后发生了什么](http://fex.baidu.com/blog/2014/05/what-happen/)
    - 推荐一本书 《网络是怎么连接的》
8. [http status code rfc2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec6.html#sec6.1.1)

## HTTP2
1. http request header

        GET / HTTP/1.1
        Host: example.com
        Connection: Upgrade, HTTP2-Settings
        Upgrade: h2c
        HTTP2-Settings: <base64url encoding of HTTP/2 SETTINGS payload>
        
2. http response header
        
        HTTP/1.1 101 Switching Protocols
        Connection: Upgrade
        Upgrade: h2c
        
        [ HTTP/2 connection ... ]

### tip ：
1. 如果服务端不支持http/2,忽略Upgrade字段，直接返回http/1.1


## HTTPS  
1. [HTTPS](./https.md)


## Websocket
1. [RFC 6455:  WebSocket Protopol](https://tools.ietf.org/html/rfc6455)
2. [WebSocket详解](http://www.52im.net/forum.php?mod=viewthread&tid=331&ctid=15)
3. request header

        GET ws://test.com/ HTTP/1.1
        Connection: Upgrade
        Upgrade: websocket
        Origin: http://example.com
        Sec-WebSocket-Version: 13
        Sec-WebSocket-Key: d4egt7snxxxxxx2WcaMQlA==
        Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits
4. response header

        HTTP/1.1 101 Switching Protocols
        Connection: Upgrade
        Upgrade: websocket
        Sec-WebSocket-Accept: gczJQPmQ4Ixxxxxx6pZO8U7UbZs=
## TCP
1. TCP 三次握手，四次挥手
    - sync攻击
2. backlog数目，sync queue, accept queue
    
3. TCP粘包，以及如何处理
    - 明确规定每个包的大小
    - 每个包规定结束的字符，可以区分包
    
4. scoket参数TCP_NODELAY,SO_RCVTIMEO和SO_SNDTIMEO,
    - TCP_NODELAY 开始Nagle算法，减少需要传输的数据包，来优化网络
    - SO_RCVTIMEO, SO_SNDTIMEO接收和发送的超时时间
    
5. 流量控制，拥塞控制
    - 滑动窗口
    
## Reading

* [keyless ssl原理](https://andblog.cn/?p=852)
* 网络各层和经典的协议
![网络层](../imgs/protocal_level.png)
* 网络各层的数据最小单位
![单位](../imgs/osi_unit.png)
* 网络通讯协议图
![tcp/ip](../imgs/tcp_ip.png)