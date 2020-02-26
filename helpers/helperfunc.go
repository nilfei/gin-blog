package helpers

import (
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Builtins = template.FuncMap{
	"showStatus": ShowStatus,
	"showtime":   ShowTime,
	"ip2long":    Ip2long,
	"long2ip":    Long2ip,
}

func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func Long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}

func ShowStatus(ts interface{}, on string, on2 string) (tsf string) {
	switch ts {
	case ts.(int):
		if ts == 0 {
			tsf = on
		} else if ts == 1 {
			tsf = on2
		}
	case ts.(string):
		if ts == "0" {
			tsf = on
		} else if ts == "1" {
			tsf = on2
		}
	}
	return tsf
}

func ShowTime(t time.Time, format string) string {
	return t.Format(format)
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func Abort(c *gin.Context, message string) {
	c.HTML(http.StatusInternalServerError, "/home/err/500.html", gin.H{
		"message": message,
	})
}

func GetMapFilterQuery(query url.Values) map[string]interface{} {
	queryMap := make(map[string]interface{})
	for key, value := range query {
		if len(value) > 0 {
			isFilter := strings.ContainsAny(key, "filter_")
			if isFilter != false {
				index := strings.Index(key, "_")
				if value[0] != "" {
					queryMap[key[index+1:len(key)]] = value[0]
				}
			}
		}

	}
	return queryMap
}
