package sherlock

// int abs
func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}
