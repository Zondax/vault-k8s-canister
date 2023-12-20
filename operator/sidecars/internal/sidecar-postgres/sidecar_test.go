package sidecarPostgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/zondax/tororu-operator/operator/common"
)

func init() {
	common.CreateCommonKubernetesClient("~/.kube/config")
}

func Test_tResInfo_rotateAndUpdate(t *testing.T) {
	t.Skip()

	dClient, _ := common.GetKubernetesClients()
	tests := []struct {
		name   string
		fields tResInfo
	}{
		{
			name: "ok test",
			fields: tResInfo{
				commChan:         make(chan string),
				dynamicClient:    dClient,
				Name:             "postgres-user-juan",
				RotationDuration: time.Duration(3200) * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commsChan := tt.fields.commChan
			tr := &tResInfo{
				Name:             tt.fields.Name,
				RotationDuration: tt.fields.RotationDuration,
				commChan:         tt.fields.commChan,
				dynamicClient:    tt.fields.dynamicClient,
			}
			go tr.rotateAndUpdateForever()
			for {
				msg := <-commsChan
				fmt.Println(msg)
			}
		})
	}
}
