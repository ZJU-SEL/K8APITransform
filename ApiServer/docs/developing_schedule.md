***每次push的时候更新下文档中对应内容***



####Api server 基本功能点
调用k8s的rest api ， 根据前端需求对api进行合理的组织以及二次封装

####基本路由结构
参见 api.odt 文档


####createenv
1 从前台接收传过来的json字段，作为基本的env环境env环境（作为创建service参数的一部分） {TomcatV string , JdkV    string , NodeNum string, Name    string , Used    int} 

2 对传递过来的参数进行验证 主要用到models/appenv中的一些工具函数。

3 验证通过后将env信息存在etcd中

4 返回details信息，details为嵌套的结构，在detail.context 和 details.children中可以继续嵌套details, ?(具体的势力信息与k8的结构如何对应？)

####问题
多用户的情况下env的存储

####Getuploadwars
扫描 application/用户名 文件夹下的内容，对文件进行处理，返回当前系统中存储的war包的名称

####问题
本地存储瓶颈问题

####Upload
1 提供用户上传上来的war包的接口，存储在application/username/appname_deploy文件夹下

####问题
控制上传的文件的大小

####deploy
1 创建deploy request实例{ EnvName  string，WarName   string，AppVersion string，IsGreyUpdating string}

2 用户传递过来的json信息赋值给deploy request ， 生成 deploy request 实例。

3 验证传递过来的参数

4 根据参数中的envname从etcd中取得env信息

5 如果不进行灰度升级（参数为0）向服务端发送get请求，得到已经创建的serveice list，(服务端的ip是在main函数的时候时候传进来的，查询的时候service的label指定为env=envname)，将这个env下的service与rc全都删除掉(名字是一样的)。

6 对上传的war包进行wartoimage的操作：根据Dockerfile模板生成新的操作，在本地build好，生成新的镜像（镜像名称为 registryname/appname.war）,push到k8master的registry中。

7 按照service生成的格式，生成 AppCreateRequest 实例，创建 service,rc，每个service 有两个label对其进行具体控制， env 以及 name

8 更新在etcd中存储的对应的env，使用数目+1

9 向服务端发送get请求，返回创建好的service的信息，还是采用嵌套的形式返回。 


####问题

如果当前的应用正在进行部署 再次点击了部署按钮的话 应该返回正在进行部署的信息

获取username的过程应该单独提炼出来





