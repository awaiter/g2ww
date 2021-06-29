package router

import (
	"bytes"
	"fmt"
	"g2ww/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber"
)

type Hook struct {
	DashboardId int    `json:"dashboardId"`
	ImageUrl    string `json:"imageUrl"`
	// Alert message
	Message  string `json:"message"`
	OrgId    int    `json:"orgId"`
	PanelId  int    `json:"panelId"`
	RuleId   int    `json:"ruleId"`
	RuleName string `json:"ruleName"`
	RuleUrl  string `json:"ruleUrl"`
	State    string `json:"state"`
	// tags     string `json:"tags"`
	// Panel Title
	Title string `json:"title"`
}

var sent_count int = 0

const (
	OKMsg       string = "恢复"
	AlertingMsg string = "触发"
	OK          string = "OK"
	Alerting    string = "Alerting"
	ColorGreen  string = "info"
	ColorGray   string = "comment"
	ColorRed    string = "warning"
	WxPosturl   string = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	Baseurl     string = "http://localhost:3000"
)

var (
	// 触发告警环境 支持单个用户
	Env      string = config.Config.Grafana.Env
	Grafaurl string = config.Config.Grafana.Url
	Atuser   string = config.Config.Webhook.AtUser
)

func GwStat() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		stat_msg := "G2WW Server created by Nova Kwok is running! \nParsed & forwarded " + strconv.Itoa(sent_count) + " messages to WeChat Work! \n"
		c.Send(stat_msg)
		return
	}
}

func GwWorker() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		h := new(Hook)
		if err := c.BodyParser(h); err != nil {
			fmt.Println(err)
			c.Send("Error on JSON format")
			return
		}

		url := WxPosturl + c.Params("key")
		// grafana_url 如果本地开启公网配置，则无需修改
		grafana_url := strings.Replace(h.RuleUrl, Baseurl, Grafaurl, -1)
		alert_url := strings.Replace(grafana_url, "edit", "view", -1)

		var (
			Color      string
			Atuer_text string
			Env_text   string
		)

		// 判断Atuser是否为空
		if Atuser != "" {
			Atuer_text = "\n<@" + Atuser + ">"
		} else {
			Atuer_text = ""
			fmt.Println("没有需要被艾特的用户")
		}

		// 判断配置环境，如果为测试环境，则不触发艾特
		if Env == "test" {
			Atuer_text = ""
			Env_text = "测试"
		} else if Env == "production" {
			Env_text = "线上"
		} else {
			Env_text = ""
		}

		// color := ColorGreen
		// 判断消息状态，如果为OK，则返回默认颜色，并不再艾特用户
		if strings.Contains(h.Title, OK) {
			h.Title = strings.ReplaceAll(h.Title, OK, OKMsg)
			Color = ColorGreen
			Atuer_text = ""
		} else {
			h.Title = strings.ReplaceAll(h.Title, Alerting, AlertingMsg)
			Color = ColorRed
		}

		if h.Message == "" {
			h.Message = "Message is empty."
		}

		// 环境；告警颜色；规则名；消息；告警地址；艾特的用户
		msgmarkdownStr := fmt.Sprintf(`
                {
                        "msgtype": "markdown",
                        "markdown": {
                          "content": "#### Grafana %s告警 \n <font color=\"%s\">%s</font> \n <font color=\"comment\">%s</font>\n [查看详情](%s) %s"
			}
                }
                `, Env_text, Color, h.Title, h.Message, alert_url, Atuer_text)

		fmt.Println(msgmarkdownStr)
		jsonStr := []byte(msgmarkdownStr)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Send("Error sending to WeChat Work API")
			return
		}
		defer resp.Body.Close()
		c.Send(resp)
		sent_count++
	}
}
