package utils

import (
	"ZinxV0.2/src/zinx/ziface"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

/*
	存储一切有关Zinx框架的全局参数,供其他模块使用
	一些参数可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	//server
	TcpServer ziface.IServer `yaml:"TcpServer"` //当前Zinx全局的Server对象
	Host      string         `yaml:"Host"`      //当前服务器主机监听的IP
	TcpPort   int            `yaml:"TcpPort"`   //当前服务器主机监听的端口号
	Name      string         `yaml:"Name"`      //当前服务器的名称

	//Zinx
	Version        string `yaml:"Version"`        //当前Zinx的版本号
	MaxConn        int    `yaml:"MaxConn"`        //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 `yaml:"MaxPackageSize"` //当前Zinx框架数据包的最大值
}

//定义一个全局对外的Globalobj
var GlobalObject *GlobalObj

//提供init方法,初始化GlobalObject
func init() {

	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//从conf/zinx.json去加载一些用户自定义的参数

	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	//file, _ := exec.LookPath(os.Args[0])
	//path, _ := filepath.Abs(file)
	//index := strings.LastIndex(path, string(os.PathSeparator))
	//path = path[:index]
	//fmt.Println("打印路径", path)
	//

	pwd, _ := os.Getwd()
	fmt.Println("当前的操作路径为:", pwd)
	data, err := ioutil.ReadFile("conf/zinx.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
	//将json文件解析gstruct
	//err = json.Unmarshal(data, &GlobalObject)
	//if err != nil {
	//	panic(err)
	//}
}
