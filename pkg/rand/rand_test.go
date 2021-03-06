package xrand

import (
	"github.com/coder2z/ndisk/pkg/constant"
	"testing"
)

func TestCreateRandomString(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{len: constant.VerificationCodeLength},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateRandomNumber(tt.args.len); got != tt.want {
				t.Errorf("CreateRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
