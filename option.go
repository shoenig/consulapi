package consulapi

type Optional struct {
	// http request options
	scheme string
	host   string

	// common api options
	ctx   Ctx
	token string

	// api coordinate options
	datacenter string
	namespace  string
	partition  string
}

func optional(opts ...Option) *Optional {
	o := new(Optional)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type Option func(*Optional)

func WithContext(ctx Ctx) Option {
	return func(o *Optional) {
		o.ctx = ctx
	}
}

func WithToken(token string) Option {
	return func(o *Optional) {
		o.token = token
	}
}

func WithDatacenter(dc string) Option {
	return func(o *Optional) {
		o.datacenter = dc
	}
}

func WithNamespace(namespace string) Option {
	return func(o *Optional) {
		o.namespace = namespace
	}
}

func WithPartition(partition string) Option {
	return func(o *Optional) {
		o.partition = partition
	}
}

func WithSchme(scheme string) Option {
	return func(o *Optional) {
		o.scheme = scheme
	}
}

func WithHost(host string) Option {
	return func(o *Optional) {
		o.host = host
	}
}
