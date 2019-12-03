package stripe

import (
	"reflect"
	"testing"
)

func TestClient_Charge(t *testing.T) {
	type fields struct {
		Key     string
		baseURL string
		httpCli httpClient
	}
	type args struct {
		amount int
		source string
		desc   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Charge
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Key:     tt.fields.Key,
				baseURL: tt.fields.baseURL,
				httpCli: tt.fields.httpCli,
			}
			got, err := c.Charge(tt.args.amount, tt.args.source, tt.args.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Charge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Charge() = %v, want %v", got, tt.want)
			}
		})
	}
}
