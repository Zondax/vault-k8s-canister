package common

import (
	"context"
	"testing"
)

func TestCreateCRDOnboardRequest(t *testing.T) {
	type args struct {
		crdId string
		ttl   uint32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ok", args: args{crdId: "testing", ttl: 100}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateCRDOnboardRequest(tt.args.crdId, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("CreateCRDOnboardRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createROAccessRequest(t *testing.T) {
	type args struct {
		ctx    context.Context
		crdIds []string
		podId  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:    context.TODO(),
				crdIds: []string{"postgres-user-prateek"},
				podId:  "testing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createROAccessRequest(tt.args.ctx, tt.args.crdIds, tt.args.podId); (err != nil) != tt.wantErr {
				t.Errorf("createROAccessRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
