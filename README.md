# 项目文档
### 服务器启动方法
> 在根目录运行make server,默认监听8888端口

### 使用第三方库
1. gnet
- ***使用原因：相较于原生net库，gnet基于事件驱动，在性能和可拓展性上要优于原生net库***

### 可能的拓展方案
1. 将存储换成redis和mysql集群，使用etcd进行服务器注册，即可对服务器进行横向拓展