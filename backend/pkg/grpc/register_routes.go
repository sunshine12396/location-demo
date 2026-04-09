package grpc

func (i *impl) RegisterRoutes(routes ...func(Server)) {
	if routes == nil {
		return
	}
	i.routes = routes
}
