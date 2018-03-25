package net

import (
	"bytes"
	"encoding/xml"
	"fmt"
	saml "github.com/russellhaering/gosaml2"
	"net/http"
	"net/url"
	"sparrow/base"
	"sparrow/oauth"
	"sparrow/provider"
	"sparrow/utils"
	"strings"
	"text/template"
	"time"
)

const respXml = `
<?xml version="1.0" encoding="UTF-8"?>
<saml2p:Response Destination="{{.DestinationUrl}}" ID="{{.RespId}}" InResponseTo="{{.ReqId}}" IssueInstant="{{.CurTime}}" Version="2.0" xmlns:saml2p="urn:oasis:names:tc:SAML:2.0:protocol">
	<saml2:Issuer xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">{{.IssuerUrl}}</saml2:Issuer>
	<saml2p:Status>
		<saml2p:StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Success"/>
	</saml2p:Status>
	<saml2:Assertion ID="{{.AssertionId}}" IssueInstant="{{.CurTime}}" Version="2.0" xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">
		<saml2:Issuer>{{.IssuerUrl}}</saml2:Issuer>
		<saml2:Subject>
			<saml2:NameID Format="urn:oasis:names:tc:SAML:2.0:nameid-format:persistent">{{.NameId}}</saml2:NameID>
			<saml2:SubjectConfirmation Method="urn:oasis:names:tc:SAML:2.0:cm:bearer">
				<saml2:SubjectConfirmationData InResponseTo="{{.ReqId}}" NotOnOrAfter="{{.NotOnOrAfter}}" Recipient="{{.DestinationUrl}}"/>
			</saml2:SubjectConfirmation>
		</saml2:Subject>
		<saml2:Conditions NotBefore="{{.CurTime}}" NotOnOrAfter="{{.NotOnOrAfter}}"/>
		<saml2:AuthnStatement AuthnInstant="{{.CurTime}}" SessionIndex="{{.SessionIndexId}}">
			<saml2:AuthnContext>
				<saml2:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport</saml2:AuthnContextClassRef>
			</saml2:AuthnContext>
		</saml2:AuthnStatement>
		<saml2:AttributeStatement>
		{{range $_, $v := .Attributes}}
			{{$v}}
		{{end}}
		</saml2:AttributeStatement>
	</saml2:Assertion>
</saml2p:Response>`

const attributeXml = `
<saml2:Attribute Name="{{.Name}}" NameFormat="{{.Format}}">
 <saml2:AttributeValue>{{.Value}}</saml2:AttributeValue>
</saml2:Attribute>`

type samlResponse struct {
	DestinationUrl string
	RespId         string
	ReqId          string
	IssuerUrl      string
	CurTime        string
	AssertionId    string
	NameId         string // only persistent format is supported
	NotOnOrAfter   string
	SessionIndexId string
	Attributes     map[string]string
	ResponseText   string
	RelayStateVal  string
}

var respTemplate *template.Template
var attributeTemplate *template.Template

func init() {
	respTemplate = &template.Template{}
	attributeTemplate = &template.Template{}
	template.Must(respTemplate.Parse(respXml))
	template.Must(attributeTemplate.Parse(attributeXml))
}

// STEP 1 - check the presence of session otherwise redirect to login
func handleSamlReq(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	if session != nil {
		// valid session exists serve the SAMLResponse
		log.Debugf("Valid session exists, sending the final response")
		sendSamlResponse(w, r, session, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("Failed to parse the request form %s", err.Error())
		log.Debugf("%s", err.Error())
		sendSamlError(w, err, http.StatusBadRequest)
		return
	}

	af := &authFlow{}
	af.From = FROM_SAML

	setAuthFlow(af, w)

	paramMap, err := parseParamMap(r)
	if err != nil {
		sendSamlError(w, err, http.StatusBadRequest)
		return
	}

	// do a redirect to /login with all the parameters
	redirect("/login", w, r, paramMap)
}

func parseParamMap(r *http.Request) (paramMap map[string]string, err error) {
	paramMap = make(map[string]string)

	for k, v := range r.Form {
		if len(v) > 1 {
			err = fmt.Errorf("Invalid request the parameter %s is included more than once", k)
			return nil, err
		}

		paramMap[k] = v[0]
	}

	return paramMap, nil
}

func sendSamlResponse(w http.ResponseWriter, r *http.Request, session *base.RbacSession, af *authFlow) {
	err := r.ParseForm()
	if err != nil {
		log.Debugf("Failed to parse the form, sending error to the user agent")
		sendSamlError(w, err, http.StatusBadRequest)
		return
	}

	var samlAuthnReq saml.AuthNRequest
	var data []byte

	samlReq := r.Form.Get("SAMLRequest")
	if r.Method == http.MethodGet {
		samlReq, err = url.QueryUnescape(samlReq)
		if err == nil {
			data = utils.B64Decode(samlReq)
		}
	} else {
		data = []byte(samlReq)
	}

	err = xml.Unmarshal(data, &samlAuthnReq)
	if err != nil {
		err = fmt.Errorf("Failed to parse the SAML authentication request %s", err.Error())
		log.Debugf("%s", err.Error())
		sendSamlError(w, err, http.StatusBadRequest)
		return
	}

	//TODO verify signature of received SAML request

	log.Debugf("Received SAMLRequest is valid, searching for client")

	pr, _ := getPrFromParam(r)
	var cl *oauth.Client
	if pr != nil {
		cl = pr.GetClient(samlAuthnReq.Issuer)
	}

	if cl == nil {
		err = fmt.Errorf("Application with issuer ID %s not found", samlAuthnReq.Issuer)
		log.Warningf("%s", err.Error())
		sendSamlError(w, err, http.StatusNotFound)
		return
	}

	genSamlResponse(w, r, pr, session, cl, samlAuthnReq)
}

func genSamlResponse(w http.ResponseWriter, r *http.Request, pr *provider.Provider, session *base.RbacSession, cl *oauth.Client, authnReq saml.AuthNRequest) {
	user, err := pr.GetUserById(session.Sub)
	if err != nil {
		err = fmt.Errorf("Error while generating SAML response for the request ID %s", authnReq.ID)
		log.Warningf("%s", err.Error())
		sendSamlError(w, err, http.StatusInternalServerError)
		return
	}

	sp := samlResponse{}
	sp.AssertionId = "_" + utils.GenUUID()

	curTime := time.Now().UTC()

	sp.CurTime = curTime.Format(time.RFC3339)
	sp.NotOnOrAfter = curTime.Add(time.Duration(cl.Saml.AssertionValidity) * time.Second).Format(time.RFC3339)
	sp.DestinationUrl = cl.Saml.ACSUrl
	sp.IssuerUrl = issuerUrl + "/saml/idp"
	sp.NameId = session.Sub
	sp.ReqId = authnReq.ID
	sp.RespId = "_" + utils.GenUUID()
	sp.SessionIndexId = "_" + utils.GenUUID()

	sp.Attributes = make(map[string]string)
	var buf bytes.Buffer
	for k, v := range cl.Saml.Attributes {
		if len(v.StaticVal) > 0 {
			if len(v.StaticMultiValDelim) == 0 {
				v.Value = v.StaticMultiValDelim
				attributeTemplate.Execute(&buf, v)
				sp.Attributes[k] = buf.String()
				buf.Reset()
			} else {
				splitValues := strings.Split(v.StaticVal, v.StaticMultiValDelim)
				for i, s := range splitValues {
					v.Value = s
					attributeTemplate.Execute(&buf, v)
					sp.Attributes[k+string(i)] = buf.String()
					buf.Reset()
				}
			}
		} else {
			val := v.GetValueFrom(user)
			if val != nil {
				v.Value = val
				attributeTemplate.Execute(&buf, v)
				sp.Attributes[k] = buf.String()
				buf.Reset()
			}
		}
	}

	respTemplate.Execute(&buf, sp)
	sp.RelayStateVal = r.Form.Get("RelayState")
	sp.ResponseText = utils.B64Encode(buf.Bytes())
	log.Debugf("SAML response: %s", sp.ResponseText)
	templates["saml_response_html"].Execute(w, sp)
}

func sendSamlError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
