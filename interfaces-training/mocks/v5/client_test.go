package stripe

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type mockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{}, nil
}

func TestClient_Charge(t *testing.T) {
	type fields struct {
		Key     string
		baseURL string
	}
	type args struct {
		cli    httpClient
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
		{
			name: "test400",
			args: args{amount: 0, source: "", desc: "",
				cli: &mockHttpClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						// do whatever you want
						return &http.Response{
							Status:     "bad_request",
							StatusCode: http.StatusBadRequest,
							Body: ioutil.NopCloser(bytes.NewBufferString(`
							{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
							"description":"Charge for demo purposes.",
							"status":"failed"}`)),
						}, nil
					},
				}},
			want:    &Charge{ID: "ch_1DEjEH2eZvKYlo2CxOmkZL4D", Amount: 2000, Description: "Charge for demo purposes.", Status: "failed"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Key:     tt.fields.Key,
				baseURL: tt.fields.baseURL,
			}
			got, err := c.Charge(tt.args.cli, tt.args.amount, tt.args.source, tt.args.desc)
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
