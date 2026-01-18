package emit

type Artifacts struct {
	TokensCSS []byte
}

func Write(_ Artifacts, _ string) error {
	return nil
}
