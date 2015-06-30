###Apiserver

根据实际的需求对于k8s的api进行调用和组装，相当于一个工作流引擎。



###ContainerCli

一个命令行工具，向Apiserver发送命令，可以进行服务的CRUD操作



###Fti

将war包打包成镜像上传的时候所用到的一些基本工具函数


###K8Apitool

对k8s的API进行了一些封装，基本上不再使用这个包中的内容，对于k8s进行调用的时候直接调用 pkg/client 的内容来发送rest API
