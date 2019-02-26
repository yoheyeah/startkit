package apis

type OPTIONS struct {
	API
	KeyValues map[string]string
}

func (o *OPTIONS) Run() (err error) {
	if len(o.KeyValues) > 0 {
		for k, v := range o.KeyValues {
			o.Ctx.Writer.Header().Set(k, v)
		}
	}
	return nil
}
