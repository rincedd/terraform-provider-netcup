package api

import (
	"bytes"
	"encoding/xml"
	"log"
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
)

const netcupWSUrl = "https://www.servercontrolpanel.de:443/SCP/WSEndUser"
const vServerRequestTpl = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:end="http://enduser.service.web.vcp.netcup.de/">
	   <soap:Header/>
	   <soap:Body>
		  <end:{{.Operation}}>
			 <loginName>{{.LoginName}}</loginName>
			 <password>{{.Password}}</password>
			 <vserverName>{{.VServerName}}</vserverName>
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

func (self *Client) sendRequest(requestBody *bytes.Buffer) (*http.Response, error) {
	resp, err := http.Post(netcupWSUrl, "text/xml", requestBody)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (self *Client) GetVServerIPs(vServerName string) ([]string, error) {
	requestBody, err := self.getVServerRequestBody("getVServerIPs", vServerName)
	if err != nil {
		return nil, err
	}
	log.Printf("Sending %s", requestBody.String())
	resp, err := self.sendRequest(requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := new(VServerIpsResponse)
	err = xml.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	log.Printf("Got %s", r)
	return r.Body.IPs, nil
}

func (self *Client) GetVServerState(vServerName string) (string, error) {
	requestBody, err := self.getVServerRequestBody("getVServerState", vServerName)
	if err != nil {
		return "", err
	}
	log.Printf("Sending %s", requestBody.String())
	resp, err := self.sendRequest(requestBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	r := new(VServerStateResponse)
	err = xml.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return "", err
	}
	log.Printf("Got %s", r)
	return r.Body.State, nil
}
