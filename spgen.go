package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

const (
	KEY_CONFIG = "Config"
)

// cli params
type cmder struct {
	configFile       string
	gener            string
	genertor         SpGenertor
	genNewConfigToml bool
}

func (c *cmder) match() error {
	if flag.NArg() > 0 {
		c.configFile = flag.Arg(0)
	}
	if c.configFile == "" {
		// set a default name
		c.configFile = "config.toml"
	}
	switch c.gener {
	case "std":
		c.genertor = new(StdGenertor)
	case "android:kt":
		c.genertor = new(KtGenerator)
	default:
		fmt.Println("unknown genertor:", c.gener, ", will use \"std\" instead of it")
		c.genertor = new(StdGenertor)
	}

	return nil
}

var cmd *cmder = &cmder{}

func init() {
	flag.StringVar(&(cmd.gener), "o", "std", "generator output plan. [std, android:kt]")
	flag.BoolVar(&(cmd.genNewConfigToml), "new", false, "create a new config template to a explicit file. (default: config.toml)")

	flag.Parse()

}

func main() {

	if err := cmd.match(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if cmd.genNewConfigToml {
		if err := createNewConfigTempalte(cmd.configFile); err != nil {
			fmt.Println(err)
		}
		return
	}

	configFile := cmd.configFile
	state, err := os.Stat(configFile)
	if (err != nil && os.IsNotExist(err)) || state.IsDir() {
		// no exist
		fmt.Println("Configuration define file no found. (", configFile, ")")
		os.Exit(1)
		return
	}
	sp, err := load(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	spConfig, groups, err := parse(sp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	// var genertor SpGenertor = new(KtGenerator) // new(StdGenertor)
	cmd.genertor.GenSp(spConfig, groups)
}

func load(file string) (sp *toml.Tree, err error) {
	sp, err = toml.LoadFile(file)
	return
}

func parse(sp *toml.Tree) (*SpConfig, []*SpGroup, error) {
	// get config
	spConfig := new(SpConfig)

	if c, ok := sp.Get(KEY_CONFIG).(*toml.Tree); ok {
		err := c.Unmarshal(spConfig)
		if err != nil {
			return nil, nil, err
		}
	}
	emptyDefault := func(origin, def string) string {
		if origin == "" {
			return def
		}
		return origin
	}
	spConfig.Author = emptyDefault(spConfig.Author, "[Spgen-go]")
	spConfig.KeyName = emptyDefault(spConfig.KeyName, "nm")
	spConfig.KeyType = emptyDefault(spConfig.KeyType, "t")
	spConfig.KeyComment = emptyDefault(spConfig.KeyComment, "cm")

	// get groups
	allKey := sp.Keys()
	indexOf := func(k string) int {
		for i, v := range allKey {
			if v == k {
				return i
			}
		}
		return -1
	}
	configIndex := indexOf(KEY_CONFIG)
	if configIndex != -1 {
		if configIndex != len(allKey)-1 {
			allKey = append(allKey[:configIndex], allKey[configIndex+1:]...)
		} else {
			allKey = allKey[:configIndex]
		}
	}

	groups := make([]*SpGroup, 0, len(allKey))

	for _, groupName := range allKey {
		// a group
		sg := new(SpGroup)
		sg.Name = groupName
		sg.Items = make([]*SpItem, 0, 16)

		if gTree, ok := sp.Get(groupName).([]*toml.Tree); ok {
			for _, itemTree := range gTree {
				spItem := parserItem(itemTree, spConfig)
				sg.Items = append(sg.Items, spItem)
			}
		}

		groups = append(groups, sg)
	}

	return spConfig, groups, nil
}

func parserItem(itemTree *toml.Tree, spConfig *SpConfig) *SpItem {
	item := new(SpItem)
	item.Name = itemTree.Get(spConfig.KeyName).(string)
	typeName := itemTree.Get(spConfig.KeyType).(string)
	typeName = strings.ToLower(typeName)
	item.Type = fromTypeName(typeName)
	item.Comment = itemTree.Get(spConfig.KeyComment).(string)
	return item
}

type SpConfig struct {
	Package    string `toml:"package,omitempty"`
	ExportDir  string `toml:"dir,omitempty"`
	Author     string `toml:"author,omitempty"`
	KeyName    string `toml:"nameKey,omitempty"`
	KeyType    string `toml:"typeKey,omitempty"`
	KeyComment string `toml:"commentKey,omitempty"`
}

type SpGroup struct {
	Name  string
	Items []*SpItem
}

type ItemType int

const (
	TypeNone   = ItemType(-1)
	TypeInt    = ItemType(iota)
	TypeLong   = ItemType(iota)
	TypeFloat  = ItemType(iota)
	TypeString = ItemType(iota)
)

func fromTypeName(typeName string) ItemType {
	switch typeName {
	case "int":
		return TypeInt
	case "long":
		return TypeLong
	case "string":
		return TypeString
	case "float":
		return TypeFloat
	default:
		return TypeNone
	}
}

type SpItem struct {
	Name    string
	Type    ItemType
	Comment string
}

func (s *SpItem) FuncName() string {
	return toFuncName(s.Name)
}

func toFuncName(name string) string {
	if name == "" {
		return name
	}
	name = strings.TrimSpace(name)
	var builder strings.Builder
	upper := false
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			if upper {
				upper = false
				if c >= 'a' && c <= 'z' {
					c = c - 32
				}
			}
			builder.WriteRune(c)
		} else {
			upper = true
		}
	}
	return builder.String()
}
