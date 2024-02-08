package server

import (
	"github.com/cossim/kit/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// HTTPService Gin HTTP的服务实例想要注入此容器，必须实现该接口
type HTTPService interface {
	InternalDependency
	// Registry 该模块所需要注册的路由：添加传递的中间件
	// [因为gin不支持路由装饰，只能使用路由分组的方式区分特殊的路由(如：添加中间件)]
	Registry(r gin.IRoutes)
	// RegistryWithMiddle 该用于注册特殊的路由：不添加传递的中间件
	RegistryWithMiddle(r gin.IRoutes)

	//Ping() error
}

// GRPCService GRPC 服务的实例想要使用 manager 进行生命周期管理，必须实现该接口
type GRPCService interface {
	InternalDependency
	Registry(s *grpc.Server)
}

type InternalDependency interface {
	// Init 初始化
	Init(cfg *config.Config) error
	// Name 服务名称
	Name() string
	// Version 服务版本信息
	Version() string
}
