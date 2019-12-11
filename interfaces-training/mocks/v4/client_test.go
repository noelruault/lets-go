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
		{
			name: "test_failed",
			fields: fields{
				httpCli: &mockHttpClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status: "bad_request",
							// StatusCode: http.StatusBadRequest,
							Body: ioutil.NopCloser(bytes.NewBufferString(`
							{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
							"description":"Charge for demo purposes.",
							"status":"failed"}`)),
						}, nil
					},
				},
			},
			args:    args{amount: 0, source: "", desc: ""},
			want:    &Charge{ID: "ch_1DEjEH2eZvKYlo2CxOmkZL4D", Amount: 2000, Description: "Charge for demo purposes.", Status: "failed"},
			wantErr: false,
		},
		{
			name: "test_success",
			fields: fields{
				httpCli: &mockHttpClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status: "success",
							// StatusCode: http.StatusOK,
							Body: ioutil.NopCloser(bytes.NewBufferString(`
							{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
							"description":"Charge for demo purposes.",
							"status":"success"}`)),
						}, nil
					},
				},
			},
			args:    args{amount: 0, source: "", desc: ""},
			want:    &Charge{ID: "ch_1DEjEH2eZvKYlo2CxOmkZL4D", Amount: 2000, Description: "Charge for demo purposes.", Status: "success"},
			wantErr: false,
		},
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
