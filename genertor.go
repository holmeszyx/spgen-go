package main

import (
	"fmt"
	"strings"
)

type SpGenertor interface {
	GenSp(spConfig *SpConfig, groups []*SpGroup)

	MapTypeSymbol(itemType ItemType) string
}

type GenEntity struct {
	Config   *SpConfig
	Group    *SpGroup
	FileName string
	FilePath string
	Date     string
}

type StdGenertor struct {
}

func (s *StdGenertor) GenSp(spConfig *SpConfig, groups []*SpGroup) {
	fmt.Println("Config:")
	fmt.Println("   package:", spConfig.Package)
	fmt.Println("   toDir:", spConfig.ExportDir)
	for _, group := range groups {
		fmt.Println("Group:", group.Name, "total:", len(group.Items))
		for _, item := range group.Items {
			fmt.Println("    item", s.itemString(item))
		}
	}
}

func (s *StdGenertor) itemString(item *SpItem) string {
	var builder strings.Builder
	builder.WriteRune('[')
	builder.WriteRune(' ')
	builder.WriteString(item.FuncName())

	builder.WriteRune('(')
	builder.WriteString(item.Name)
	builder.WriteRune(')')

	builder.WriteRune(' ')
	builder.WriteString(s.MapTypeSymbol(item.Type))
	builder.WriteRune(' ')

	builder.WriteString(" : ")
	builder.WriteString(item.Comment)

	builder.WriteRune(' ')
	builder.WriteRune(']')
	return builder.String()
}

func (s *StdGenertor) MapTypeSymbol(itemType ItemType) string {
	switch itemType {
	case TypeNone:
		return "UnkownType"
	case TypeInt:
		return "Int"
	case TypeFloat:
		return "Float"
	case TypeLong:
		return "Long"
	case TypeString:
		return "String"
	default:
		return "WTF"
	}
}
