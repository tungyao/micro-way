package register

import (
	"fmt"
)

// 设计 每一个 服务 都是一个单独channel
// chan *Service
var GlobalContainer = new(Container) // 所有 service 容器 都存放在这里 有没有问题 ,我不知道
var GlobalPosition map[string]int = make(map[string]int)
var GlobalMdString map[string]string = make(map[string]string)

func LoadGlobalService(rulers map[string]*Ruler) { // 从全局加载文件  计算md5
	names := make([]string, 0)
	for k, v := range rulers {
		names = append(names, k)
		GlobalContainer.Number = GlobalContainer.Number + 1
		GlobalMdString[k] = MD(k)
		GlobalPosition[MD(k)] = GlobalContainer.Number
		GlobalContainer.Rulers = append(GlobalContainer.Rulers, v)
	}
	fmt.Print("\n------------------------------------\n|\tContainer have *", GlobalContainer.Number, "* Services\t\t|\n------------------------------------")
	for k, v := range names {
		fmt.Print("\n|\tService ", k, "\t\t\t", v, "\t\t|\n------------------------------------")
	}
	fmt.Println("\n\t\t\tIT'S STARTING ....")
}
func LoadSingleService(serviceName string, ruler Ruler) bool { // 手动 或者 新发现服务 注册服务 调用
	for _, v := range GlobalContainer.Rulers {
		if v.Name == serviceName {
			fmt.Println("=>\tNew Service has be registered =>", v.Name)
			return false
		}
	}
	fmt.Println("=>\tNew Service is registering =>", serviceName)
	GlobalContainer.mux.Lock()
	GlobalContainer.Number = GlobalContainer.Number + 1
	GlobalContainer.Rulers = append(GlobalContainer.Rulers, &ruler)
	GlobalContainer.mux.Unlock()
	return true
}

// issues 频繁得调用 不知道会不会出现问题
func GetStatusSingleService(serviceName string) (bool, int, *Service) { // 获取 单个 服务 状态 , 用户 可以 调用  返回值 isDie , status
	gd := GlobalMdString[serviceName]
	if GlobalPosition[gd] != 0 {
		return GlobalContainer.
				Rulers[GlobalPosition[gd]-1].
				IsDie,
			GlobalContainer.
				Rulers[GlobalPosition[gd]-1].
				Status,
			GlobalContainer.
				Rulers[GlobalPosition[gd]-1].Service

	}
	return true, -2, nil
}
func SetStatusSingleService(serviceName string) { // 设置单个 服务 状态 , 用户 禁止调用  设置服务状态 返回时延

}

// Monitoring service status
func MonitorService() {

}
