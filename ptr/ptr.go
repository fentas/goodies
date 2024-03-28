package ptr

func To[T any](p T) *T {
	return &p
}

func Safe[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}
