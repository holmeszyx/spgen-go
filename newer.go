package main

import (
	"fmt"
	"os"
)

const newerTemplate = `# "Config" is the reserved word which to mark the env

[Config]
package="[the package like 'a.b.c']"
dir= "./exportDir"
#author= "[Author name or empty]"
#nameKey="nm"
#typeKey="t"
#commentKey="cm"
#defaultKey="def"

# === The perference Group ===
# === [[Group]] ===

# === nm : name, cm: comment , t: type===

[[Group1]]
    nm = "ItemName1"
    t = "ItemType[int, long, float, string]"
	cm = "Comment"
	
[[Group1]]
    nm = "ItemName2"
    t = "string"
	cm = "Comment2"	
	
[[Group2]]
    nm = "ItemName3"
    t = "float"
    cm = "Comment3"	
`

func createNewConfigTempalte(file string) error {
	_, err := os.Stat(file)
	if err == nil || !os.IsNotExist(err) {
		return fmt.Errorf("the file \"%s\" is exists.", file)
	}
	of, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer of.Close()
	_, e := of.WriteString(newerTemplate)
	return e
}
