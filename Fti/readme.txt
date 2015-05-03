将deployments目录下的文件打包成 tar.gz文件放在tardir目录下
可能要检验一下 目标文件是否已经存在
参考了这个 注意文件头信息的提取方式 不要手动填入

之后要一个docker环境 把deployments.tar.gz 放入到一个base image中
这个base image 里面有tomcat的镜像
之后可能需要通过编写dockerfile 把tar包注入进来 生成新的镜像
在镜像中要把 tar.gz给解压缩开
放在固定的那个文件夹下