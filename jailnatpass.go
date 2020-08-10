package main

type JailNATPass struct {
	JailName string `json:"jail_name"`
	GwIf     string `json:"gw_if"`
	logger   func(int, string)
}

func (fwd *JailNATPass) SetLogger(f func(int, string)) {
	fwd.logger = f
}

func (fwd *JailNATPass) Create() error {
	return nil
}

func (fwd *JailNATPass) Remove() error {
	return nil
}

func NewJailNATPass(gw_if string, jail_name string) *JailNATPass {
	nat := &JailNATPass{GwIf: gw_if, JailName: jail_name}
	return nat
}
