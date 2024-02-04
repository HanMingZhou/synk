package controller

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddressesController
func AddressesController(c *gin.Context) {
	addrs, err := net.InterfaceAddrs() // 获取所有的ip地址
	if err != nil {
		log.Fatal(err)
	}
	var result []string
	for _, address := range addrs { // 遍历所有的ip
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}
