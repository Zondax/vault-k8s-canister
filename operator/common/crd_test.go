package common

import (
	"context"
	"github.com/zondax/tororu-operator/operator/common/v1"
	"testing"
)

func init() {
	CreateCommonKubernetesClient("~/.kube/config")
}

func TestGetTResFromName(t *testing.T) {
	t.Skip()
	type args struct {
		ctx       context.Context
		name      string
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.TororuResource
		wantErr bool
	}{
		{
			name: "ok test",
			args: args{
				ctx:       context.TODO(),
				name:      "postgres-user-juan",
				namespace: "default",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTResFromName(tt.args.ctx, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTResFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
