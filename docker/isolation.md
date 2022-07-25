# 如何实现Docker和宿主机的隔离的

1. 进程隔离 Namespace

2. 文件隔离 （Control Group(cGroup)、mount 挂载方式）

3. 资源限制 Control Group(cGroup)

# Docker 有哪些概念

1. dockerfile

2. image

3. container

docker build dockerfile then generate image, docker run image then generate container

#Docker 有哪些底层镜像

1. alphin 5.7M
2. ubuntu 大
3. centos 大

# Docker pull, Docker Push