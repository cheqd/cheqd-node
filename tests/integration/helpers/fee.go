package helpers

func GenerateFees(amount, feePayer string) []string {
	return []string{
		"--fees", amount,
		"--fee-payer", feePayer,
		"--broadcast-mode", "block",
	}
}
