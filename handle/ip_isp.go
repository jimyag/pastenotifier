package handle

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/jimyag/log"
)

var (
	templateUrl             = "http://ip-api.com/json/%s?fields=status,message,country,city,isp,reverse,query&lang=zh-CN"
	_, private24BitBlock, _ = net.ParseCIDR("10.0.0.0/8")
	_, private20BitBlock, _ = net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ = net.ParseCIDR("192.168.0.0/16")
)

type IpIsp struct{}

func (i *IpIsp) Handle(content string) (title string, message string, err error) {
	if content == "" {
		return "", "", nil
	}
	ip := net.ParseIP(content)
	if ip == nil {
		return "", "", nil
	}
	if i.isPrivateIP(ip) {
		return "", "", nil
	}
	respData, err := i.getInfo(ip)
	if err != nil || respData == nil {

		return "", "", err
	}
	if respData.Status != "success" {
		log.Error().Str("message", respData.Message).Msg("ipisp error")
		return "", "", nil
	}
	return content, respData.format(), nil
}

func (h *IpIsp) isPrivateIP(ip net.IP) bool {
	return private24BitBlock.Contains(ip) || private20BitBlock.Contains(ip) || private16BitBlock.Contains(ip)
}

func (h *IpIsp) getInfo(ip net.IP) (*responseData, error) {
	u := fmt.Sprintf(templateUrl, ip.String())
	resp, err := http.DefaultClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData := responseData{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	return &respData, nil
}

// responseData https://ip-api.com/docs/api:json
type responseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Country string `json:"country"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
	Reverse string `json:"reverse"`
	Query   string `json:"query"`
}

func (data *responseData) format() string {
	if data.Message != "" {
		return data.Message
	}

	str := ""
	if data.Country != "" {
		str = str + "country: " + data.Country + "\n"
	}
	if data.City != "" {
		str = str + "city: " + data.City + "\n"
	}
	if data.ISP != "" {
		str = str + "isp: " + data.ISP + "\n"
	}
	if data.Reverse != "" {
		str = str + "reverse: " + data.Reverse + "\n"
	}
	return str
}
