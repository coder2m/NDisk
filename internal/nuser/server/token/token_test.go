package token

import (
	"testing"
)

func TestAccessTokenTicket_Create(t *testing.T) {
	type fields struct {
		AccessToken  string
		RefreshToken string
	}
	type args struct {
		uid uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccessTokenTicket{
				AccessToken:  tt.fields.AccessToken,
				RefreshToken: tt.fields.RefreshToken,
			}
			if err := a.Encode(tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("%+v\n\n\n\n", a)
				t.Log(a.Decode())
			}
		})
	}

}
