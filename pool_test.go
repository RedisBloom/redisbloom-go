package redis_bloom_go

import (
	"testing"
)

func TestNewMultiHostPool(t *testing.T) {
	type args struct {
		hosts    []string
		authPass *string
	}
	tests := []struct {
		name         string
		args         args
		wantPoolSize int
		wantConntNil bool
	}{
		{"same connection string", args{[]string{"localhost:6379", "localhost:6379"}, nil}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *MultiHostPool
			got = NewMultiHostPool(tt.args.hosts, tt.args.authPass)
			if len(got.hosts) != tt.wantPoolSize {
				t.Errorf("NewMultiHostPool() = %v, want %v", got, tt.wantPoolSize)
			}
			if gotConn := got.Get(); tt.wantConntNil == false && gotConn == nil {
				t.Errorf("NewMultiHostPool().Get() = %v, want %v", gotConn, tt.wantConntNil)
			}
		})
	}
}
