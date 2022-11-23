package helpers

func GenerateFees(amount string) []string {
	return []string{
		"--fees", amount,
		"--broadcast-mode", "block",
	}
}
