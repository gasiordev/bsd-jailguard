package main

type JailPortFwd struct {
	SrcIf   string `json:"src_if"`
	SrcPort string `json:"src_port"`
	DstJail string `json:"dst_jail"`
	DstPort string `json:"dst_port"`
	logger  func(int, string)
}

func (fwd *JailPortFwd) SetLogger(f func(int, string)) {
	fwd.logger = f
}

func (fwd *JailPortFwd) Add() error {
	return nil
}

func (fwd *JailPortFwd) Delete() error {
	return nil
}

func NewJailPortFwd(src_if string, src_port string, dst_jail string, dst_port string) *JailPortFwd {
	fwd := &JailPortFwd{SrcIf: src_if, SrcPort: src_port, DstJail: dst_jail, DstPort: dst_port}
	return fwd
}
