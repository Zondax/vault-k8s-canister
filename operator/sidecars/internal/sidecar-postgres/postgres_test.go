package sidecarPostgres

import (
	"testing"
)

func Test_checkAndCreateDatabase(t *testing.T) {
	t.Skip()
	type args struct {
		dbName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ok test", args: args{dbName: "tester"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkAndCreateDatabase(tt.args.dbName); (err != nil) != tt.wantErr {
				t.Errorf("checkAndCreateDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrUpdateUserOnDatabase(t *testing.T) {
	t.Skip()
	type args struct {
		dbName       string
		userName     string
		userPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ok test", args: args{dbName: "testing", userName: "prateek", userPassword: "check"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateUserOnDatabase(tt.args.dbName, tt.args.userName, tt.args.userPassword); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateUserOnDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
