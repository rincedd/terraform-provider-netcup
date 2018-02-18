package api

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"text/template"
)

type (
	Client struct {
		LoginName string
		Password  string
	}

	VServerIpsResponseBody struct {
		XMLName xml.Name `xml:"Body"`
		IPs     []string `xml:"getVServerIPsResponse>return"`
	}

	VServerIpsResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    VServerIpsResponseBody
	}

	VServerStateResponseBody struct {
		XMLName xml.Name `xml:"Body"`
		State   string   `xml:"getVServerStateResponse>return"`
	}

	VServerStateResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    VServerStateResponseBody
	}

	VServerNicknameResponseBody struct {
		XMLName  xml.Name `xml:"Body"`
		Nickname string   `xml:"getVServerNicknameResponse>return"`
	}

	VServerNicknameResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    VServerNicknameResponseBody
	}
)

const netcupWSUrl = "https://www.servercontrolpanel.de:443/SCP/WSEndUser"
const vServerRequestTpl = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:end="http://enduser.service.web.vcp.netcup.de/">
	   <soap:Header/>
	   <soap:Body>
		  <end:{{.Operation}}>
			 <loginName>{{.LoginName}}</loginName>
			 <password>{{.Password}}</password>
			 <vserverName>{{.VServerName}}</vserverName>
			 <vservername>{{.VServerName}}</vservername>
		  </end:{{.Operation}}>
	   </soap:Body>
	</soap:Envelope>`

func (self *Client) getVServerRequestBody(operation string, vServerName string) (*bytes.Buffer, error) {
	tpl, err := template.New("vServerRequest").Parse(vServerRequestTpl)
	if err != nil {
		return nil, err
	}
	tplData := struct {
		Operation   string
		LoginName   string
		Password    string
		VServerName string
	}{Operation: operation, LoginName: self.LoginName, Password: self.Password, VServerName: vServerName}
	requestBody := bytes.Buffer{}
	err = tpl.Execute(&requestBody, tplData)
	if err != nil {
		return nil, err
	}
	return &requestBody, nil
}

func (self *Client) sendRequest(requestBody *bytes.Buffer, responseData interface{}) error {
	resp, err := http.Post(netcupWSUrl, "text/xml", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = xml.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return err
	}
	return nil
}

func (self *Client) GetVServerIPs(vServerName string) ([]string, error) {
	requestBody, err := self.getVServerRequestBody("getVServerIPs", vServerName)
	if err != nil {
		return nil, err
	}
	r := new(VServerIpsResponse)
	err = self.sendRequest(requestBody, r)

	return r.Body.IPs, nil
}

func (self *Client) GetVServerState(vServerName string) (string, error) {
	requestBody, err := self.getVServerRequestBody("getVServerState", vServerName)
	if err != nil {
		return "", err
	}
	r := new(VServerStateResponse)
	err = self.sendRequest(requestBody, r)
	if err != nil {
		return "", err
	}
	return r.Body.State, nil
}

func (self *Client) GetVServerNickname(vServerName string) (string, error) {
	requestBody, err := self.getVServerRequestBody("getVServerNickname", vServerName)
	if err != nil {
		return "", err
	}
	r := new(VServerNicknameResponse)
	err = self.sendRequest(requestBody, r)
	if err != nil {
		return "", err
	}
	return r.Body.Nickname, nil
}
