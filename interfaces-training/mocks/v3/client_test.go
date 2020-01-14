package stripe

import (
	"fmt"
	"testing"
)

type StripeClientMock struct {
	c   *Charge
	err error
}

func (scm *StripeClientMock) Charge(amount int, source, desc string) (*Charge, error) {
	return scm.c, scm.err
}

func TestApp_Run(t *testing.T) {
	tests := []struct {
		name    string
		client  StripeClient
		wantErr bool
	}{
		{
			name: "test1",
			client: &StripeClientMock{
				c: &Charge{},
			},
			wantErr: false,
		},
		{
			name: "test2",
			client: &StripeClientMock{
				err: fmt.Errorf("mock error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				sc: tt.client,
			}
			if err := a.Run(); (err != nil) != tt.wantErr {
				t.Errorf("App.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
