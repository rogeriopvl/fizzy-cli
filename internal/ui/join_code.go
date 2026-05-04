package ui

import (
	"fmt"
	"io"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func DisplayJoinCode(w io.Writer, jc *fizzy.JoinCode) error {
	fmt.Fprintf(w, "Code: %s\n", jc.Code)
	fmt.Fprintf(w, "Active: %v\n", jc.Active)
	fmt.Fprintf(w, "Usage: %d/%d\n", jc.UsageCount, jc.UsageLimit)
	fmt.Fprintf(w, "URL: %s\n", jc.URL)
	return nil
}
