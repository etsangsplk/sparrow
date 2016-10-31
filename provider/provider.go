package provider

import (
	"bytes"
	"crypto"
	"fmt"
	_ "github.com/dgrijalva/jwt-go"
	logger "github.com/juju/loggo"
	"io/ioutil"
	"path/filepath"
	"sparrow/base"
	"sparrow/conf"
	_ "sparrow/rbac"
	"sparrow/schema"
	"sparrow/silo"
	"strings"
)

type Provider struct {
	Schemas   map[string]*schema.Schema       // a map of Schema ID to Schema
	RsTypes   map[string]*schema.ResourceType // a map of Name to ResourceTye
	RtPathMap map[string]*schema.ResourceType // a map of EndPoint to ResourceTye
	config    *conf.Config
	sl        *silo.Silo
	layout    *Layout
	Name      string // the domain name
	PubKey    crypto.PublicKey
	PrivKey   crypto.PrivateKey
}

var log logger.Logger

func init() {
	log = logger.GetLogger("sparrow.scim.provider")
}

func NewProvider(layout *Layout) (prv *Provider, err error) {
	schemas, err := base.LoadSchemas(layout.SchemaDir)
	if err != nil {
		return nil, err
	}

	prv = &Provider{}
	prv.Schemas = schemas

	prv.RsTypes, prv.RtPathMap, err = base.LoadResTypes(layout.ResTypesDir, prv.Schemas)
	if err != nil {
		return nil, err
	}

	prv.config, err = conf.ParseConfig(filepath.Join(layout.ConfDir, "domain.json"))
	if err != nil {
		return nil, err
	}

	dataFilePath := filepath.Join(layout.DataDir, layout.name)

	prv.sl, err = silo.Open(dataFilePath, prv.config, prv.RsTypes, prv.Schemas)

	if err != nil {
		return nil, err
	}

	prv.layout = layout
	prv.Name = layout.name
	prv.sl.Engine.Domain = layout.name

	err = prv.createDefaultResources()

	return prv, err
}

func (prv *Provider) createDefaultResources() error {
	// TODO how to identify admin groups uniquely when the name is not unique??
	adminGroup := `{"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
				    "displayName": "Administrators",
				    "permissions": [{"value": "READ"}, {"value": "CREATE"}, {"value": "UPDATE"}, {"value": "DELETE"}]
				   }`

	buf := bytes.NewBufferString(adminGroup)

	grpRes, err := base.ParseResource(prv.RsTypes, prv.Schemas, buf)
	if err != nil {
		return err
	}

	grpRes, err = prv.sl.Insert(grpRes)
	if err != nil {
		return err
	}

	log.Infof("Successfully inserted default admin group")

	adminUser := `{"schemas":["urn:ietf:params:scim:schemas:core:2.0:User"],
                   "userName":"admin",
                   "displayName":"Administrator",
                   "password":"secret",
                   "emails":[
                       {
                         "value":"admin@%s",
                         "type":"work",
                         "primary":true
                       }
                     ]
                   }`

	adminUser = fmt.Sprintf(adminUser, prv.Name) // set domain name in email
	buf = bytes.NewBufferString(adminUser)
	userRes, err := base.ParseResource(prv.RsTypes, prv.Schemas, buf)
	if err != nil {
		return err
	}

	subAtMap := make(map[string]interface{})
	subAtMap["value"] = grpRes.GetId()
	userRes.AddCA("groups", subAtMap)

	_, err = prv.sl.Insert(userRes)
	if err != nil {
		return err
	}

	log.Infof("Successfully inserted default administrator user")

	return nil
}

func (prv *Provider) GetSchemaJsonArray() string {
	json := "["

	for _, v := range prv.Schemas {
		json += v.Text + ","
	}

	json = strings.TrimSuffix(json, ",")

	return json + "]"
}

func (prv *Provider) GetSchema(id string) (string, error) {
	sc := prv.Schemas[id]

	if sc == nil {
		return "", fmt.Errorf("No schema present with the ID %s", id)
	}

	return sc.Text, nil
}

func (prv *Provider) GetResTypeJsonArray() string {
	json := "["

	for _, v := range prv.RsTypes {
		json += v.Text + ","
	}

	json = strings.TrimSuffix(json, ",")

	return json + "]"
}

func (prv *Provider) GetResourceType(name string) (string, error) {
	rt := prv.RsTypes[name]

	if rt == nil {
		return "", fmt.Errorf("No resource type present with the ID %s", name)
	}

	return rt.Text, nil
}

func (prv *Provider) GetConfigJson() (data []byte, err error) {
	f := filepath.Join(prv.layout.ConfDir, "config.json")
	return ioutil.ReadFile(f)
}

func (prv *Provider) CreateResource(crCtx *base.CreateContext) (res *base.Resource, err error) {
	return prv.sl.Insert(crCtx.InRes)
}

func (prv *Provider) DeleteResource(delCtx *base.DeleteContext) error {
	return prv.sl.Delete(delCtx.Rid, delCtx.Rt)
}

func (prv *Provider) GetResource(getCtx *base.GetContext) (res *base.Resource, err error) {
	return prv.sl.Get(getCtx.Rid, getCtx.Rt)
}

func (prv *Provider) Search(sc *base.SearchContext, outPipe chan *base.Resource) {
	prv.sl.Search(sc, outPipe)
}

func (prv *Provider) Replace(replaceCtx *base.ReplaceContext) (res *base.Resource, err error) {
	return prv.sl.Replace(replaceCtx.InRes)
}

func (prv *Provider) Patch(patchCtx *base.PatchContext) (res *base.Resource, err error) {
	return prv.sl.Patch(patchCtx.Rid, patchCtx.Pr, patchCtx.Rt)
}

func (prv *Provider) Authenticate(ar *base.AuthRequest) (authToken string, err error) {
	user, err := prv.sl.Authenticate(ar.Username, ar.Password)

	if err != nil {
		msg := "Invalid username or password"
		log.Debugf(msg)
		return "", base.NewForbiddenError(msg)
	}

	session := prv.sl.Engine.NewRbacSession(user)
	authToken = session.ToJwt(prv.PrivKey)

	return authToken, nil
}