package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		want     []byte
	}{
		{OpConstant, []int{65534}, []byte{OpConstant.Byte(), 255, 254}},
	}

	for _, tt := range tests {
		insn := Make(tt.op, tt.operands...)

		if len(insn) != len(tt.want) {
			t.Errorf("instruction has wrong length; want=%d, got=%d", len(tt.want), len(insn))
		}

		for i, b := range tt.want {
			if insn[i] != b {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, insn[i])
			}
		}
	}
}
