/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/incubator-answer-plugins/connector-wallet/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/google/go-github/v50/github"
	"github.com/segmentfault/pacman/log"

	"golang.org/x/oauth2"
)

type WalletConnector struct {
	// Config *ConnectorConfig
}

// type ConnectorConfig struct {
// 	ClientID     string `json:"client_id"`
// 	ClientSecret string `json:"client_secret"`
// }

func init() {
	plugin.Register(&WalletConnector{
		// Config: &ConnectorConfig{},
	})
}

// Implement the Base interface
func (g *WalletConnector) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "wallet_connector",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "i-Luicfer",
		Version:     "0.0.1",
		Link:        "https://github.com/apache/incubator-answer-plugins/tree/main/connector-wallet",
	}
}

// Implement the Connector plugin interface
func (g *WalletConnector) ConnectorLogoSVG() string {
	return `PHN2ZyB0PSIxNzE3ODM1NzkwNTM1IiBjbGFzcz0iaWNvbiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9IjMyNjU3IiB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCI+PHBhdGggZD0iTTgxMC42NjY2NjcgMjk4LjY2NjY2N2gtNDIuNjY2NjY3VjI1NmExMjggMTI4IDAgMCAwLTEyOC0xMjhIMjEzLjMzMzMzM2ExMjggMTI4IDAgMCAwLTEyOCAxMjh2NTEyYTEyOCAxMjggMCAwIDAgMTI4IDEyOGg1OTcuMzMzMzM0YTEyOCAxMjggMCAwIDAgMTI4LTEyOHYtMzQxLjMzMzMzM2ExMjggMTI4IDAgMCAwLTEyOC0xMjh6TTIxMy4zMzMzMzMgMjEzLjMzMzMzM2g0MjYuNjY2NjY3YTQyLjY2NjY2NyA0Mi42NjY2NjcgMCAwIDEgNDIuNjY2NjY3IDQyLjY2NjY2N3Y0Mi42NjY2NjdIMjEzLjMzMzMzM2E0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxIDAtODUuMzMzMzM0eiBtNjQwIDQyNi42NjY2NjdoLTQyLjY2NjY2NmE0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxIDAtODUuMzMzMzMzaDQyLjY2NjY2NnogbTAtMTcwLjY2NjY2N2gtNDIuNjY2NjY2YTEyOCAxMjggMCAwIDAgMCAyNTZoNDIuNjY2NjY2djQyLjY2NjY2N2E0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxLTQyLjY2NjY2NiA0Mi42NjY2NjdIMjEzLjMzMzMzM2E0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxLTQyLjY2NjY2Ni00Mi42NjY2NjdWMzc2Ljc0NjY2N0ExMjggMTI4IDAgMCAwIDIxMy4zMzMzMzMgMzg0aDU5Ny4zMzMzMzRhNDIuNjY2NjY3IDQyLjY2NjY2NyAwIDAgMSA0Mi42NjY2NjYgNDIuNjY2NjY3eiIgcC1pZD0iMzI2NTgiPjwvcGF0aD48L3N2Zz4=`
}

func (g *WalletConnector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *WalletConnector) ConnectorSlugName() string {
	return "wallet"
}

func (g *WalletConnector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	fmt.Printf("receiverURL: (%s) \n", receiverURL)
	redirectURL = "https://www.baidu.com/"
	fmt.Printf("redirectURL: (%s) \n", redirectURL)
	return ""
}

func (g *WalletConnector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	return userInfo, nil
}

// Implement the Translator interface
// 每个插件都有一个配置界面，这里用于在配置界面，显示配置信息
func (g *WalletConnector) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name: "client_id",
			Type: plugin.ConfigTypeInput,
			// Title:       plugin.MakeTranslator(i18n.ConfigClientIDTitle),
			// Description: plugin.MakeTranslator(i18n.ConfigClientIDDescription),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			// Value: g.Config.ClientID,
		},
		{
			Name: "client_secret",
			Type: plugin.ConfigTypeInput,
			// Title:       plugin.MakeTranslator(i18n.ConfigClientSecretTitle),
			// Description: plugin.MakeTranslator(i18n.ConfigClientSecretDescription),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			// Value: g.Config.ClientSecret,
		},
	}
}

// 接收配置字段
func (g *WalletConnector) ConfigReceiver(config []byte) error {
	// c := &ConnectorConfig{}
	// _ = json.Unmarshal(config, c)
	// g.Config = c
	return nil
}

// 绑定电子邮箱
func (g *WalletConnector) guaranteeEmail(email string, accessToken string) string {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	))
	client.Timeout = 15 * time.Second
	cli := github.NewClient(client)

	emails, _, err := cli.Users.ListEmails(context.Background(), &github.ListOptions{Page: 1})
	if err != nil {
		log.Error(err)
		return ""
	}
	for _, e := range emails {
		if e.GetPrimary() {
			return e.GetEmail()
		}
	}
	return email
}
