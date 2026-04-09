package grpc

func (i *impl) Routes() {
	if i.routes == nil {
		return
	}
	for _, route := range i.routes {
		route(i)
	}
}
