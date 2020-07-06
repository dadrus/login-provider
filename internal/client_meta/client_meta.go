package client_meta

import (
	"github.com/mitchellh/mapstructure"
	"login-provider/internal/utils"
)

type ScopeInfo struct {
	Scope       string
	Mandatory   bool
	Description string
}

type DataAccessArea struct {
	Description string
	ScopeInfos  []ScopeInfo
}

type ClientMetaInfo struct {
	AskConsent        bool              `json:"ask_consent" mapstructure:"ask_consent"`
	MandatoryScopes   []string          `json:"mandatory_scopes" mapstructure:"mandatory_scopes"`
	ScopeDescriptions map[string]string `json:"scope_descriptions" mapstructure:"scope_descriptions"`
}

func (cmi *ClientMetaInfo) Unmarshal(data interface{}) error {
	return mapstructure.Decode(data, cmi)
}

func (cmi *ClientMetaInfo) CreateScopeInfos(scopes []string) []ScopeInfo {
	var scopeInfos []ScopeInfo
	for _, scope := range scopes {
		if utils.Contains(cmi.MandatoryScopes, scope) {
			scopeInfos = append(scopeInfos, ScopeInfo{
				Scope:       scope,
				Mandatory:   true,
				Description: cmi.ScopeDescriptions[scope],
			})
		} else {
			scopeInfos = append(scopeInfos, ScopeInfo{
				Scope:       scope,
				Mandatory:   false,
				Description: cmi.ScopeDescriptions[scope],
			})
		}
	}
	return scopeInfos
}
