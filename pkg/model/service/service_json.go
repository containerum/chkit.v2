package service

import (
	"encoding/json"

	"github.com/containerum/chkit/pkg/model"
	kubeModel "github.com/containerum/kube-client/pkg/model"
)

var (
	_ model.JSONrenderer = Service{}
	_ json.Marshaler     = Service{}
)

func (serv Service) RenderJSON() (string, error) {
	data, err := serv.MarshalJSON()
	return string(data), err
}

func (serv Service) MarshalJSON() ([]byte, error) {
	serv.ToKube()
	data, err := json.MarshalIndent(serv.ToKube(), "", model.Indent)
	return data, err
}

func (serv *Service) UnmarshalJSON(data []byte) error {
	var kubeSvc kubeModel.Service
	if err := json.Unmarshal(data, &kubeSvc); err != nil {
		return err
	}
	*serv = ServiceFromKube(kubeSvc)
	return nil
}
