package main

import (
//"errors"
)

type HistoryEntry struct {
	Entry   string `json:"entry"`
	Created string `json:"created"`
}

func NewHistoryEntry(c string, e string) *HistoryEntry {
	he := &HistoryEntry{Entry: e, Created: c}
	return he
}
