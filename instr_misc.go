package i8080

func instrNOP(op uint8, c *CPU) uint64 {
	return 4
}

func instrHLT(op uint8, c *CPU) uint64 {
	c.Halted = true
	return 7
}

func instrOUT(op uint8, c *CPU) uint64 {
	_ = insArg8(c)
	return 10
}

func instrIN(op uint8, c *CPU) uint64 {
	port := insArg8(c)
	c.Registers[A] = port // TODO
	return 10
}

func instrDI(op uint8, c *CPU) uint64 {
	// TODO
	return 4
}

func instrEI(op uint8, c *CPU) uint64 {
	// TODO
	return 4
}
