package limit_reader

type ReaderBoundsExceededError struct {
}

func (e ReaderBoundsExceededError) Error() string {
	return "exceeded bounds of reader"
}
