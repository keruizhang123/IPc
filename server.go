package ipconf

import (
	"github.com/KRZ/common/config"
	"github.com/KRZ/ipconf/domain"
	"github.com/KRZ/ipconf/source"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RunMain Start the web container
func RunMain(path string) {
	config.Init(path)
	source.Init() //The data source should be started first
	domain.Init() // Initialize the scheduling layer
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList) //get model port
	s.Spin()
}
