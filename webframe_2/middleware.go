package webframe_2

type Middleware func(next HandleFunc) HandleFunc
