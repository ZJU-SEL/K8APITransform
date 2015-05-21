server 端 把可以使用的 base image弄成 tar包的形式

由于镜像对于所有用户都是可见的 因此server端存储镜像不需要使用多租户的形式

client 端 pull 过来之后 直接导入成为镜像的形式

添加time out的限制

stop

查看应用运行状态

image类型看一下 下午把 list 完成 本地和server端统一好

在server 和 client 端维护一个链表

用户pull的时候 同时发一份给server 以及 registry 
