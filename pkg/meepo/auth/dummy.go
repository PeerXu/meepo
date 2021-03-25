package auth

type DummyEngine struct{}

func (*DummyEngine) Sign(payload Context) (Context, error) {
	return map[string]interface{}{
		CONTEXT_NAME: "dummy",
	}, nil
}

func (*DummyEngine) Verify(payload, signature Context) error {
	return nil
}

func NewDummyEngine(...NewEngineOption) (Engine, error) {
	return &DummyEngine{}, nil
}

func init() {
	RegisterNewEngineFunc("dummy", NewDummyEngine)
}
