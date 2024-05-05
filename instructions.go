package i8080

func insArg2(op uint8) uint8 {
	return (op >> 3) & 0x6 // 0, 2, 4, 6, ie every 2nd 8-bit register
}

func insArg3(op uint8) uint8 {
	return op & 0x7
}

func insArg3b(op uint8) uint8 {
	return (op >> 3) & 0x7
}

func insArg8(c *CPU) (ret uint8) {
	ret = c.Read(c.PC)
	c.PC++
	return
}

func insArg16(c *CPU) (ret uint16) {
	ret = c.Read16(c.PC)
	c.PC += 2
	return
}

func insGetreg8(c *CPU, reg uint8) uint8 {
	switch reg {
	case M:
		return c.Read(c.HL())
	default:
		return c.Registers[reg]
	}
}

func insSetreg8(c *CPU, reg uint8, val uint8) {
	switch reg {
	case M:
		c.Write(c.HL(), val)
	default:
		c.Registers[reg] = val
	}
}

var ops = [256]func(uint8, *CPU) uint64{
	instrNOP, instrLXI, instrSTAX, instrINX, instrINR, instrDCR, instrMVI, instrRLC, instrNOP, instrDAD, instrLDAX, instrDCX, instrINR, instrDCR, instrMVI, instrRRC, // 0x
	instrNOP, instrLXI, instrSTAX, instrINX, instrINR, instrDCR, instrMVI, instrRAL, instrNOP, instrDAD, instrLDAX, instrDCX, instrINR, instrDCR, instrMVI, instrRAR, // 1x
	instrNOP, instrLXI, instrSHLD, instrINX, instrINR, instrDCR, instrMVI, instrDAA, instrNOP, instrDAD, instrLHLD, instrDCX, instrINR, instrDCR, instrMVI, instrCMA, // 2x
	instrNOP, instrLXI, instrSTA, instrINX, instrINR, instrDCR, instrMVI, instrSTC, instrNOP, instrDAD, instrLDA, instrDCX, instrINR, instrDCR, instrMVI, instrCMC, // 3x
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, // 4x
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, // 5x
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, // 6x
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrHLT, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, // 7x
	instrADD, instrADD, instrADD, instrADD, instrADD, instrADD, instrADD, instrADD, instrADC, instrADC, instrADC, instrADC, instrADC, instrADC, instrADC, instrADC, // 8x
	instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, // 9x
	instrANA, instrANA, instrANA, instrANA, instrANA, instrANA, instrANA, instrANA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, // Ax
	instrORA, instrORA, instrORA, instrORA, instrORA, instrORA, instrORA, instrORA, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, // Bx
	instrCondRET, instrPOP, instrCondJMP, instrJMP, instrCondCALL, instrPUSH, instrADI, instrRST, instrCondRET, instrRET, instrCondJMP, instrJMP, instrCondCALL, instrCALL, instrACI, instrRST, // Cx
	instrCondRET, instrPOP, instrCondJMP, instrOUT, instrCondCALL, instrPUSH, instrSUI, instrRST, instrCondRET, instrRET, instrCondJMP, instrIN, instrCondCALL /*instrBIOS*/, instrNOP, instrSBI, instrRST, // Dx
	instrCondRET, instrPOP, instrCondJMP, instrXTHL, instrCondCALL, instrPUSH, instrANI, instrRST, instrCondRET, instrPCHL, instrCondJMP, instrXCHG, instrCondCALL, instrCALL, instrXRI, instrRST, // Ex
	instrCondRET, instrPOP, instrCondJMP, instrDI, instrCondCALL, instrPUSH, instrORI, instrRST, instrCondRET, instrSPHL, instrCondJMP, instrEI, instrCondCALL, instrCALL, instrCPI, instrRST, // Fx
}
