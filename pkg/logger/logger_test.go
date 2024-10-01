package logger

import "testing"

func TestInitialize(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				level: "Info",
			},
			wantErr: false,
		},
		{
			name: "second",
			args: args{
				level: "Debug",
			},
			wantErr: false,
		},
		{
			name: "third",
			args: args{
				level: "Fatal",
			},
			wantErr: false,
		},
		{
			name: "4th",
			args: args{
				level: "Error",
			},
			wantErr: false,
		},
		{
			name: "5th",
			args: args{
				level: "Unknown",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Initialize(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
