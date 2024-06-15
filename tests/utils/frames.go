package utils

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/muesli/termenv"
)

var escapes = []string{
	termenv.CSI + termenv.HideCursorSeq,
	termenv.CSI + termenv.ShowCursorSeq,

	termenv.CSI + termenv.EnableBracketedPasteSeq,
	termenv.CSI + termenv.DisableBracketedPasteSeq,
	termenv.CSI + termenv.StartBracketedPasteSeq,
	termenv.CSI + termenv.EndBracketedPasteSeq,

	termenv.CSI + termenv.EnableMousePressSeq,
	termenv.CSI + termenv.DisableMousePressSeq,
	termenv.CSI + termenv.EnableMouseSeq,
	termenv.CSI + termenv.EnableMouseHiliteSeq,
	termenv.CSI + termenv.DisableMouseHiliteSeq,
	termenv.CSI + termenv.EnableMouseCellMotionSeq,
	termenv.CSI + termenv.DisableMouseCellMotionSeq,
	termenv.CSI + termenv.EnableMouseAllMotionSeq,
	termenv.CSI + termenv.DisableMouseAllMotionSeq,
	termenv.CSI + termenv.EnableMouseExtendedModeSeq,
	termenv.CSI + termenv.DisableMouseExtendedModeSeq,
	termenv.CSI + termenv.EnableMousePixelsModeSeq,
	termenv.CSI + termenv.DisableMousePixelsModeSeq,
	termenv.CSI + termenv.SaveCursorPositionSeq,
	termenv.CSI + termenv.RestoreCursorPositionSeq,

	termenv.CSI + termenv.CursorUpSeq,
	termenv.CSI + termenv.CursorDownSeq,
	termenv.CSI + termenv.CursorForwardSeq,
	termenv.CSI + termenv.CursorBackSeq,
	termenv.CSI + termenv.CursorNextLineSeq,
	termenv.CSI + termenv.CursorPreviousLineSeq,
	termenv.CSI + termenv.CursorHorizontalSeq,
	termenv.CSI + termenv.CursorPositionSeq,
	termenv.CSI + termenv.ScrollUpSeq,
	termenv.CSI + termenv.ScrollDownSeq,
	termenv.CSI + termenv.ChangeScrollingRegionSeq,
	termenv.CSI + termenv.InsertLineSeq,
	termenv.CSI + termenv.DeleteLineSeq,
}

var splits = []string{
	termenv.CSI + termenv.EraseLineRightSeq,
	termenv.CSI + termenv.EraseLineLeftSeq,
	termenv.CSI + termenv.EraseEntireLineSeq,
	termenv.CSI + termenv.RestoreScreenSeq,
	termenv.CSI + termenv.SaveScreenSeq,
	termenv.CSI + termenv.AltScreenSeq,
	termenv.CSI + termenv.ExitAltScreenSeq,
	termenv.CSI + termenv.EraseLineSeq,
	termenv.CSI + termenv.EraseDisplaySeq,
}

var patternReplacePairs = []string{
	"%d", "\\d+",
	"[", "\\[",
	"]", "\\]",
	"?", "\\?",
}

// Frames returns text-only terminal text frames by whole termenv output with special symbols.
func Frames(output []byte) []string {
	if len(output) == 0 {
		return nil
	}

	chunks := make([]string, 0, len(escapes))
	for _, escape := range escapes {
		chunk := strings.NewReplacer(patternReplacePairs...).Replace(escape)
		chunks = append(chunks, chunk)
	}
	pattern := "(" + strings.Join(chunks, "|") + ")"
	replaced := regexp.MustCompile(pattern).ReplaceAllString(string(output), "")

	chunks = make([]string, 0, len(escapes))
	for _, split := range splits {
		chunk := strings.NewReplacer(patternReplacePairs...).Replace(split)
		chunks = append(chunks, chunk)
	}
	pattern = "(" + strings.Join(chunks, "|") + ")"
	rawFrames := regexp.MustCompile(pattern).Split(replaced, -1)

	results := make([]string, 0, 1)
	for _, rawFrame := range rawFrames {
		if rawFrame != "" {
			results = append(results, rawFrame)
		}
	}
	return results
}

func ToLF(text []byte) []byte {
	return bytes.ReplaceAll(text, []byte("\r"), []byte(""))
}
