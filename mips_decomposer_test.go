package gapstone

import "testing"
import "bytes"
import "fmt"
import "io/ioutil"

func mipsInsnDetail(insn Instruction, engine *Engine, buf *bytes.Buffer) {

	if len(insn.Mips.Operands) > 0 {
		fmt.Fprintf(buf, "\top_count: %v\n", len(insn.Mips.Operands))
	}

	for i, op := range insn.Mips.Operands {
		switch op.Type {
		case MIPS_OP_REG:
			fmt.Fprintf(buf, "\t\toperands[%v].type: REG = %v\n", i, engine.RegName(op.Reg))
		case MIPS_OP_IMM:
			fmt.Fprintf(buf, "\t\toperands[%v].type: IMM = 0x%x\n", i, (uint64(op.Imm)))
		case MIPS_OP_MEM:
			fmt.Fprintf(buf, "\t\toperands[%v].type: MEM\n", i)
			if op.Mem.Base != MIPS_REG_INVALID {
				fmt.Fprintf(buf, "\t\t\toperands[%v].mem.base: REG = %s\n",
					i, engine.RegName(op.Mem.Base))
			}
			if op.Mem.Disp != 0 {
				fmt.Fprintf(buf, "\t\t\toperands[%v].mem.disp: 0x%x\n", i, uint64(op.Mem.Disp))
			}
		}

	}

	fmt.Fprintf(buf, "\n")
}

func TestMips(t *testing.T) {

	final := new(bytes.Buffer)
	spec_file := "mips.SPEC"

	for i, platform := range mips_tests {

		engine, err := New(platform.arch, platform.mode)
		if err != nil {
			t.Errorf("Failed to initialize engine %v", err)
			return
		}
		if i == 0 {
			maj, min := engine.Version()
			t.Logf("Arch: Mips. Capstone Version: %v.%v", maj, min)
		}
		defer engine.Close()
		insns, err := engine.Disasm([]byte(platform.code), offset, 0)
		if err == nil {
			fmt.Fprintf(final, "****************\n")
			fmt.Fprintf(final, "Platform: %s\n", platform.comment)
			fmt.Fprintf(final, "Code:")
			dumpHex([]byte(platform.code), final)
			fmt.Fprintf(final, "Disasm:\n")
			for _, insn := range insns {
				fmt.Fprintf(final, "0x%x:\t%s\t%s\n", insn.Address, insn.Mnemonic, insn.OpStr)
				mipsInsnDetail(insn, &engine, final)
			}
			fmt.Fprintf(final, "0x%x:\n", insns[len(insns)-1].Address+insns[len(insns)-1].Size)
			fmt.Fprintf(final, "\n")
		} else {
			t.Errorf("Disassembly error: %v\n", err)
		}

	}

	spec, err := ioutil.ReadFile(spec_file)
	if err != nil {
		t.Errorf("Cannot read spec file %v: %v", spec_file, err)
	}
	if fs := final.String(); string(spec) != fs {
		//fmt.Println(fs)
		t.Errorf("Output failed to match spec!")
	} else {
		t.Logf("Clean diff with %v.\n", spec_file)
	}

}
