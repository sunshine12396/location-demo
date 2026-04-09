package stderr

func (i *impl) Message() string {
	return i.ErrorMessage
}
