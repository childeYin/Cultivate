# k8s

1. 隔离技术 Namesapce 
	- 创建线程 clone(main_function,stack_size,CLONE_NEWPID | SIGCHLD,NULL)
	- PID Namespace, 
	- Mount, UTs, IPC, Network,User Namespace 
	- 容器只能看到当前Namespce所限定的资源，文件，设备，状态或者配置等

2. 容器资源限制 Cgroups (Linux Control Group, CPU,内存，磁盘，网络，进程的优先级等) /sys/fs/cgroups
	- mount -t cgroup
	- cfs_period
	- cfg_quota
	- blkio 为块设备I/O限制 一般用于磁盘等设备
	- cpuset 为进程分配单独的cup核 和对应的内存节点
	- memory 为进程内存限制

3. Mount Namespce 
	- 它对容器的进程视图的改变 是伴随着挂载操作才生效的
	- chroot (change root file system)
	- rootfs (/bin, /etc, /proc 是一个操作系统所包含的文件呢、配置、和目录 并不包含操作系统内核)
	 
4. 虚拟机 Vs 容器
	1. 虚拟机通过硬件虚拟化，和操作系统os交互，i/o等会带来性能的损耗；容器化后的应用只是一个普通进程，不会因为虚拟化技术产生性能损耗
	2. 敏捷，高性能是容器比虚拟机好的地方
	3. 容器隔离不彻底

5. 容器里面执行top,free等是从/proc目录中的相关文件读取的，容器并没有对/proc, /sys等文件系统做隔离，因此独到的信息是宿主机的信息
	- lxcfs 是fuse文件系统， 维护容器内部的/proc文件

6. 容器膨胀
	- writeout （当你删除一个文件时，它只在读/写层被标记为已删除，但它仍然存在于下面的一个或多个层中。而且，删除文件不会减少镜像的大小。）
	- 

7. 如何再容器里修改ubuntu镜像的内容呢
	- copy on write(容器的文件系统需要被修改的时候，整个文件从只读层被复制到容器的读写层。每个容器拥有自己的读写层，所以共享文件的修改再任何其他容器内是不可见的)
	- 	

8. unionFS 联合文件系统
	- OverlayFS, 
	- AUFS, 
	- Btrfs, 
	- VFS, 
	- ZFS 
	- Device Mapper

5. 总结：
	容器是一个”单进程“模型

	- 启用Linux Namespace 配置
	- 设置指定的Cgroups
	- 切换进程的根目录(Docker 先用pivot_root, 没有才会选择chroot)

6. Pod
	- 凡是调度、网络、存储，以及安全相关的属性，基本上是 Pod 级别的。
	- lifecycle(postStart, preStop)
	- 生命周期
		- Pending (yaml文件已经提交给了k8s, api对象已经被创建并保存在了etcd中)
		- Running (容器都已经创建成功，并且至少有一个再运行)
		- Succeeded (pod所有的容器都正常运行，并且已经退出)
		- Failed (Pod至少有一个容器以不正常的状态退出)
		- Unknown (异常状态，pod状态不能持续被kublet汇报给kube-apiserver, 可能通信出了问题)

7. Kubernetes 支持的 Projected Volume 一共有四种
	- Secret (保存base64加密的数据)
	- ConfigMap (普通配置)
	- Downward API (Downward API 能够获取到的信息，一定是 Pod 里的容器进程启动 之前就能够确定下来的信息。)
	- ServiceAccountToken

8. Pod 
	- ServiceAccount
	- 容器健康检查 (Probe 探针)
		- livenessProbe(存活探针 检测容器是否正常运行 exec, initialDelaySeconds:5 每5秒执行一次， http, tcp )
		- readinessProbe(就绪探针 判定容器是否启动完成)
	- 恢复机制(restartPolicy)
		- 当前节点上，恢复pod(deployment 可以帮pod迁移到其他的节点上)
		- Always
		- OnFailure
		- Never
	- 只要pod的restartPolicy 指定的策略允许异常重启，那么Pod将会保持Running 状态，并将容器进行重启
	- 对于多个容器的pod，只有所有的容器进行异常状态后，Pod才会进入Failed状态

9. deployment 
	- deployment 控制 ReplicaSet 的数目以及版本， ReplicaSet控制Pod的个数
	- 水平扩展、收缩、
		- ReplicaSet 
		- kubectl scale deployment nginx-deployment --replicas=4
		- kubectl get deployment 
	- 滚动更新
		- 将一个集群中正在运行的多个 Pod 版本，交替地逐一升级的过程，就是“滚动更
新”
	- 随机字符串叫作 pod-template-hash，保证pod不混淆
	- RollingUpdateStrategy


10. StatefulSet 
	- 编排有状态的应用
	- 保证了 Pod 网络标识的稳定性

11. DaemonSet
	- 守护进程

12. Job CronJob
	- 












