package guid

func (i *impl) NewID() string {
	return i.uuidFn()
}
