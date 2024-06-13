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

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/incubator-answer-plugins/connector-github/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/google/go-github/v50/github"
	"github.com/segmentfault/pacman/log"
	"golang.org/x/oauth2"
	oauth2GitHub "golang.org/x/oauth2/github"
)

type Connector struct {
	Config *ConnectorConfig
}

type ConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func init() {
	plugin.Register(&Connector{
		Config: &ConnectorConfig{},
	})
}

func (g *Connector) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "github_connector",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "answerdev",
		Version:     "1.2.6",
		Link:        "https://github.com/apache/incubator-answer-plugins/tree/main/connector-github",
	}
}

func (g *Connector) ConnectorLogoSVG() string {
	return `PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0Ij48cGF0aCBkPSJNMTIgMGMtNi42MjYgMC0xMiA1LjM3My0xMiAxMiAwIDUuMzAyIDMuNDM4IDkuOCA4LjIwNyAxMS4zODcuNTk5LjExMS43OTMtLjI2MS43OTMtLjU3N3YtMi4yMzRjLTMuMzM4LjcyNi00LjAzMy0xLjQxNi00LjAzMy0xLjQxNi0uNTQ2LTEuMzg3LTEuMzMzLTEuNzU2LTEuMzMzLTEuNzU2LTEuMDg5LS43NDUuMDgzLS43MjkuMDgzLS43MjkgMS4yMDUuMDg0IDEuODM5IDEuMjM3IDEuODM5IDEuMjM3IDEuMDcgMS44MzQgMi44MDcgMS4zMDQgMy40OTIuOTk3LjEwNy0uNzc1LjQxOC0xLjMwNS43NjItMS42MDQtMi42NjUtLjMwNS01LjQ2Ny0xLjMzNC01LjQ2Ny01LjkzMSAwLTEuMzExLjQ2OS0yLjM4MSAxLjIzNi0zLjIyMS0uMTI0LS4zMDMtLjUzNS0xLjUyNC4xMTctMy4xNzYgMCAwIDEuMDA4LS4zMjIgMy4zMDEgMS4yMy45NTctLjI2NiAxLjk4My0uMzk5IDMuMDAzLS40MDQgMS4wMi4wMDUgMi4wNDcuMTM4IDMuMDA2LjQwNCAyLjI5MS0xLjU1MiAzLjI5Ny0xLjIzIDMuMjk3LTEuMjMuNjUzIDEuNjUzLjI0MiAyLjg3NC4xMTggMy4xNzYuNzcuODQgMS4yMzUgMS45MTEgMS4yMzUgMy4yMjEgMCA0LjYwOS0yLjgwNyA1LjYyNC01LjQ3OSA1LjkyMS40My4zNzIuODIzIDEuMTAyLjgyMyAyLjIyMnYzLjI5M2MwIC4zMTkuMTkyLjY5NC44MDEuNTc2IDQuNzY1LTEuNTg5IDguMTk5LTYuMDg2IDguMTk5LTExLjM4NiAwLTYuNjI3LTUuMzczLTEyLTEyLTEyeiIvPjwvc3ZnPg==`
}

func (g *Connector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *Connector) ConnectorSlugName() string {
	return "github"
}

// 点击登录时，组装参数
// param receiverURL github跳回到我们的地址
// return redirectURL 我们跳到github的地址
// http://localhost/answer/api/v1/connector/redirect/github
// https://github.com/login/oauth/authorize?client_id=Ov23liLXpj4g5UkCr4Y1&redirect_uri=
func (g *Connector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	pp("ConnectorSender", receiverURL)
	oauth2Config := &oauth2.Config{
		ClientID:     g.Config.ClientID,
		ClientSecret: g.Config.ClientSecret,
		Endpoint:     oauth2GitHub.Endpoint,
		RedirectURL:  receiverURL,
		Scopes:       []string{"user:email"},
	}
	redirectURL = oauth2Config.AuthCodeURL("state")
	pp("ConnectorSender", redirectURL)
	return redirectURL
}

// 点击登录后，github跳转回时，调用的函数
// http://localhost/answer/api/v1/connector/redirect/github
func (g *Connector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	code := ctx.Query("code")
	pp("ConnectorReceiver", receiverURL, code) // code = ba828f7273312a21970f

	// Exchange code for token
	oauth2Config := &oauth2.Config{
		ClientID:     g.Config.ClientID,
		ClientSecret: g.Config.ClientSecret,
		Endpoint:     oauth2GitHub.Endpoint,
	}
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return userInfo, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	// Exchange token for user info
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	))
	client.Timeout = 15 * time.Second
	cli := github.NewClient(client)
	resp, _, err := cli.Users.Get(context.Background(), "")
	if err != nil {
		return userInfo, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	metaInfo, _ := json.Marshal(resp)
	userInfo = plugin.ExternalLoginUserInfo{
		ExternalID:  fmt.Sprintf("%d", resp.GetID()),
		DisplayName: resp.GetName(),
		Username:    resp.GetLogin(),
		Email:       resp.GetEmail(),
		MetaInfo:    string(metaInfo),
		Avatar:      resp.GetAvatarURL(),
	}

	// guarantee email was verified
	userInfo.Email = g.guaranteeEmail(userInfo.Email, token.AccessToken)

	pj("ConnectorReceiver", userInfo)
	return userInfo, nil
}

func (g *Connector) guaranteeEmail(email string, accessToken string) string {
	fmt.Println("guaranteeEmail:", email, accessToken)
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
			fmt.Println("guaranteeEmail:", email)
			return e.GetEmail()
		}
	}
	fmt.Println("guaranteeEmail:", email)
	return email
}

// 列出插件
func (g *Connector) ConfigFields() []plugin.ConfigField {
	pp("ConfigFields")
	return []plugin.ConfigField{
		{
			Name:        "client_id",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigClientIDTitle),
			Description: plugin.MakeTranslator(i18n.ConfigClientIDDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientID,
		},
		{
			Name:        "client_secret",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigClientSecretTitle),
			Description: plugin.MakeTranslator(i18n.ConfigClientSecretDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientSecret,
		},
	}
}

// 保存配置
func (g *Connector) ConfigReceiver(config []byte) error {
	pp("ConfigReceiver", string(config))
	c := &ConnectorConfig{}
	_ = json.Unmarshal(config, c)
	g.Config = c
	return nil
}

func pp(params ...interface{}) {
	for index, param := range params {
		fmt.Printf("index (%d) param (%v) \n", index, param)
	}
}

func pv(fn string, param interface{}) {
	fmt.Printf("index (%s) param (%+v) \n", fn, param)
}

func pj(fn string, param interface{}) {
	result, _ := json.Marshal(param)
	fmt.Printf("index (%s) param (%s) \n", fn, result)
}
