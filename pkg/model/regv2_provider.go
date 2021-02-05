package model

import "encoding/json"

type ProviderInfo struct {
	RegisterKey
	Meta registerServerMeta // 注册的数据结构
}

// {"name":"","schema":"grpc","address":"10.117.22.137:18090","labels":{"enable":"true","env":"","group":"default","hostname":"askuy","region":"unknown","upTimestamp":"1570675391","vcsInfo":"","weight":"50","zone":"unknown"},"services":{"testproto:GreeterGrpc$IGreeter":{"namespace":"testproto","name":"GreeterGrpc$IGreeter","labels":{"dubbo":"2.0.2","name":"GreeterGrpc$IGreeter","namespace":"testproto","release":"2.7.0"},"methods":[]}}}
type registerServerMeta struct {
	originValue []byte
	Name        string `json:"name"`
	Scheme      string `json:"schema"`
	Address     string `json:"address"`
	Labels      struct {
		Enable      string `json:"enable"`
		Env         string `json:"env"`
		Group       string `json:"group"`
		Hostname    string `json:"hostname"`
		Region      string `json:"region"`
		UpTimestamp string `json:"startTs"` // 启动时间
		VcsInfo     string `json:"vcsInfo"`
		Weight      string `json:"weight"`
		Zone        string `json:"zone"`
	} `json:"labels"`
	Services map[string]registerOneService `json:"services"`
}

func (k *ProviderInfo) ParseValue(bs []byte) (err error) {
	if err = json.Unmarshal(bs, &k.Meta); err != nil {
		return
	}
	k.Meta.originValue = bs
	return nil
}

// todo methods
type registerOneService struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Labels    struct {
		Dubbo     string `json:"dubbo"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Release   string `json:"release"`
	} `json:"labels"`
	Methods []interface{} `json:"methods"` // rpcname
}

func (k *ProviderInfo) Region() string {
	return k.Meta.Labels.Region
}

func (k *ProviderInfo) Zone() string {
	return k.Meta.Labels.Zone
}

func (k *ProviderInfo) Weight() string {
	return k.Meta.Labels.Weight
}

func (k *ProviderInfo) Enable() string {
	return k.Meta.Labels.Enable
}

func (k *ProviderInfo) Group() string {
	return k.Meta.Labels.Group
}

func (k *ProviderInfo) Value() []byte {
	return k.Meta.originValue
}
