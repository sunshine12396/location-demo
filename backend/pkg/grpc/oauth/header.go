package oauth

import (
	"encoding/base64"
	"encoding/json"
)

type Header interface {
	Base64Encode() (string, error)
}

type headerImpl struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Cty string `json:"cty"`
}

func newHeader() Header {
	return &headerImpl{
		Alg: AlgorithmHS256, // is default AlgorithmHS256
		Typ: TypeBearerJWT,
		Cty: ContentTypeJSON,
	}
}

func (h *headerImpl) Base64Encode() (string, error) {
	headerBytes, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(headerBytes), nil
}
