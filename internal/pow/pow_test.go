package pow

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
)

func Test_PowCalcValidate(t *testing.T) {
	type args struct {
		targetBits uint16
		keySize    int32
		data       []byte
		solution   uint64
	}

	tests := []struct {
		name       string
		args       args
		shouldFail bool
	}{
		{
			name: "Test1",
			args: args{
				targetBits: 20,
				keySize:    40,
				data:       generatePhrase(40),
				solution:   0,
			},
			shouldFail: false,
		},
		{
			name: "Test2",
			args: args{
				targetBits: 24,
				keySize:    40,
				data:       []byte{20, 3, 225, 240, 199, 46, 63, 138, 19, 100, 110, 250, 8, 245, 69, 199, 64, 189, 173, 206, 13, 215, 22, 198, 83, 55, 98, 101, 106, 129, 59, 104, 45, 64, 8, 29, 187, 156, 127, 163},
				solution:   1234065,
			},
			shouldFail: false,
		},
		{
			name: "Test3",
			args: args{
				targetBits: 24,
				keySize:    40,
				data:       []byte{20, 3, 225, 240, 199, 46, 63, 138, 19, 100, 110, 250, 8, 245, 69, 199, 64, 189, 173, 206, 13, 215, 22, 198, 83, 55, 98, 101, 106, 129, 59, 104, 45, 64, 8, 29, 187, 156, 127, 163},
				solution:   1234066,
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := Calc(context.Background(), tt.args.data, tt.args.targetBits)
			if tt.args.solution != 0 {
				result = tt.args.solution
			}
			isValid := Validate(tt.args.data, tt.args.targetBits, result)
			if isValid && tt.shouldFail {
				t.Errorf("Test failed on %s", tt.name)
			}
		})
	}
}

var table = []struct {
	complexity uint16
}{
	{complexity: 2},
	{complexity: 8},
	{complexity: 16},
	{complexity: 20},
	{complexity: 24},
	{complexity: 25},
}

func BenchmarkCalculateProof(b *testing.B) {
	data := generatePhrase(40)

	for _, v := range table {
		b.Run(fmt.Sprintf("target_bits_%d", v.complexity), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Calc(context.Background(), data, v.complexity)
			}
		})
	}
}

func BenchmarkCheckProof(b *testing.B) {
	data := generatePhrase(40)
	var targetBits uint16 = 25

	for _, v := range table {
		ans := Calc(context.Background(), data, targetBits)
		b.ResetTimer()
		b.Run(fmt.Sprintf("target_bits_%d", v.complexity), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Validate(data, targetBits, ans)
			}
		})
	}
}

func generatePhrase(len int32) (data []byte) {
	data = make([]byte, len)
	rand.Read(data)
	return
}
