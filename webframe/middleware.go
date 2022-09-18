package webframe

type Middleware func(next HanleFunc) HanleFunc
