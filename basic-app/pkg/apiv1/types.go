package apiv1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)


type Controller struct {
	ApiVersion      string `json:"apiVersion"`
	Kind            string `json:"kind"`
	meta.ObjectMeta `json:"metadata"`
}

type SyncRequest struct {
	Parent Controller `json:"parent"`
}