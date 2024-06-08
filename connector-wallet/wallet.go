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
	return `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`
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
