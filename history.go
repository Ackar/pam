package main

import (
	"bufio"
	"fmt"
	"os"
)

type history struct {
	lines []string
}

func newHistory() *history {
	return &history{}
}

func (h *history) filename() string {
	return os.Getenv("HOME") + "/.pam_history"
}

func (h *history) load() []string {
	file, err := os.Open(h.filename())
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res
}

func (h *history) add(s string) {
	h.lines = append(h.lines, s)
}

func (h *history) save() {
	file, err := os.OpenFile(h.filename(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return
	}
	defer file.Close()

	for _, l := range h.lines {
		fmt.Fprintln(file, l)
	}
}
