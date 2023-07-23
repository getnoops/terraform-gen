package ast

type Source struct {
	// Name is the filename of the source
	Name string
	// Input is the actual contents of the source file
	Input string
}

type Position struct {
	Start  int     // The starting position, in runes, of this token in the input.
	End    int     // The end position, in runes, of this token in the input.
	Line   int     // The line number at the start of this item.
	Column int     // The column number at the start of this item.
	Src    *Source // The source document this token belongs to
}
