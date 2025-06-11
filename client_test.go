package tbot

import (
	"reflect"
	"testing"
)

func TestClient_Me(t *testing.T) {
	type fields struct {
		token   string
		baseURL string
		url     string
	}
	var tests []struct {
		name    string
		fields  fields
		want    *User
		wantErr bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				token:   tt.fields.token,
				baseURL: tt.fields.baseURL,
				url:     tt.fields.url,
			}
			got, err := c.Me()
			if (err != nil) != tt.wantErr {
				t.Errorf("Me() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Me() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		token   string
		baseUrl string
	}
	var tests []struct {
		name string
		args args
		want *Client
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.token, tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
