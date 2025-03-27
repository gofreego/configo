package parser

import (
	"context"
	"testing"
)

type TestConfig struct {
}

type TestConfigChild struct {
	IsTest bool `name:"isTest" type:"boolean" description:"name of the config child" required:"false"`
}

type TestConfig1 struct {
	Name  string          `name:"name" type:"string" description:"name of the config" required:"true"`
	Child TestConfigChild `name:"child" type:"parent" description:"child test config"`
}

func Test_marshal(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		cfg any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "simple config with 1 field",
			args: args{
				ctx: ctx,
				cfg: &TestConfig{},
			},
			want:    "null",
			wantErr: false,
		},
		{
			name: "config with child config",
			args: args{
				ctx: ctx,
				cfg: &TestConfig1{
					Name: "test",
					Child: TestConfigChild{
						IsTest: true,
					},
				},
			},
			want:    `[{"name":"name","type":"string","description":"name of the config","required":true,"value":"test","children":null},{"name":"child","type":"parent","description":"child test config","required":false,"value":null,"children":[{"name":"isTest","type":"boolean","description":"name of the config child","required":false,"value":true,"children":null}]}]`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshal(t *testing.T) {
	type args struct {
		ctx   context.Context
		value string
		cfg   any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple config with 1 field",
			args: args{
				ctx:   context.Background(),
				value: "null",
				cfg:   &TestConfig{},
			},
			wantErr: false,
		},
		{
			name: "config with child config",
			args: args{
				ctx:   context.Background(),
				value: `[{"name":"name","type":"string","description":"name of the config","required":true,"value":"test","children":null},{"name":"child","type":"parent","description":"child test config","required":false,"value":null,"children":[{"name":"isTest","type":"boolean","description":"name of the config child","required":false,"value":true,"children":null}]}]`,
				cfg:   &TestConfig1{},
			},
			wantErr: false,
		},
		{
			name: "config with child config with value",
			args: args{
				ctx:   context.Background(),
				value: `[{"name":"name","type":"string","description":"name of the config","required":true,"value":"test","children":null},{"name":"child","type":"parent","description":"child test config","required":false,"value":null,"children":[{"name":"isTest","type":"boolean","description":"name of the config child","required":false,"value":true,"children":null}]}]`,
				cfg: &TestConfig1{
					Name: "test",
					Child: TestConfigChild{
						IsTest: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.ctx, tt.args.value, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
