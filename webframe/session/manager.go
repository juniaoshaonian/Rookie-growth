package session

import "webframe"

type Manager struct {
	Store
	Propagator
	SessionName string
}

func (m *Manager) InitSession(ctx *webframe.Context, id string) (Session, error) {
	sess, err := m.Genrate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	err = m.Inject(id, ctx.Resp)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *Manager) GetSession(ctx *webframe.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	if sess, ok := ctx.UserValues[m.SessionName]; ok {
		return sess.(Session), nil
	}
	val, ok := ctx.UserValues[m.SessionName]
	if ok {
		return val.(Session), nil
	}
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.SessionName] = sess
	return sess, nil
}

func (m *Manager) ReflashSession(ctx *webframe.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Reflash(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	err = m.Inject(sess.ID(), ctx.Resp)

	return err
}

func (m *Manager) RemoveSession(ctx *webframe.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
