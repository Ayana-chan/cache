## 分布式缓存
### 启动方法
&emsp;&emsp;启动一个节点： 
```run go run_nodeserver.go -端口号 -lru容量```

&emsp;&emsp;启动一个Router：
```run go run_router.go 多个节点地址（以空格分隔） -端口号 -每个真实节点对应的虚拟节点数```

### 基本信息
&emsp;&emsp;每个节点用LRU算法来管理数据，Router通过一致性哈希算法来管理节点。LRU算法通过HashMap+双向链表实现，一致性哈希通过数组+排序实现。当前仅实现了查询与设值。

### HTTP API
&emsp;&emsp;GET {Router地址}/cache/{key值}&emsp;&emsp;来进行查询

&emsp;&emsp;POST {Router地址}/cache/{key值}&emsp;&emsp;来进行设值，值通过http body发送

### 接下来可能还会补充的功能
&emsp;&emsp;在线增删节点

&emsp;&emsp;~~Router感知节点宕机（现在如果一个节点宕机整个系统就崩溃了）~~ (2021年12月30日完成)

&emsp;&emsp;主从节点

&emsp;&emsp;持久化

&emsp;&emsp;... ...

### 2021年12月30日更新
&emsp;&emsp;添加了Router感知节点宕机的功能，当查询或设值的节点已经无法使用，Router会删除该节点。在设值时会自动反复重新尝试设值，若所有节点均失效，则Router关闭。