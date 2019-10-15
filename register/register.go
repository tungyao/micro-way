package register

import (
	"fmt"
	"os"
	"os/exec"
)

// 设计 每一个 服务 都是一个单独channel
// chan *Service

var Containers = new(Container) // 所有 service 容器 都存放在这里 有没有问题 ,我不知道

func LoadGlobalService(rulers map[string]Ruler) { // 从全局加载文件
	names := make([]string, 0)
	for k, v := range rulers {
		names = append(names, k)
		Containers.Number = Containers.Number + 1
		Containers.Rulers = append(Containers.Rulers, v)
	}
	fmt.Print("\n------------------------------------\n|\tContainer have *", Containers.Number, "* Services\t\t|\n------------------------------------")
	for k, v := range names {
		fmt.Print("\n|\tService ", k, "\t\t\t", v, "\t\t|\n------------------------------------")
	}
	fmt.Println("\n\t\t\tIT'S STARTING ....")
}
func LoadSingleService(serviceName string, ruler Ruler) bool { // 手动 或者 新发现服务 注册服务 调用
	for _, v := range Containers.Rulers {
		if v.Name == serviceName {
			fmt.Println("=>\tNew Service has be registered =>", v.Name)
			return false
		}
	}
	fmt.Println("=>\tNew Service is registering =>", serviceName)
	Containers.mux.Lock()
	defer Containers.mux.Unlock()
	Containers.Number = Containers.Number + 1
	Containers.Rulers = append(Containers.Rulers, ruler)
	FlushScreen()
	return true
}
func GetStatusSingleService(serviceName string) { // 获取 单个 服务 状态 , 用户 可以 调用

}
func SetStatusSingleService(serviceName string) { // 设置单个 服务 状态 , 用户 禁止调用

}
func FlushScreen() {
	cmd := exec.Command("cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
