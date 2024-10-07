package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/osbuild/images/pkg/blueprint"
	"github.com/spf13/cobra"
)

func maybeDeref(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Pointer:
		return t.Elem()
	default:
		return t
	}
}

func findTag(keys []string, t reflect.Type) {
	if len(keys) == 0 {
		return
	}
	t = maybeDeref(t)
	for idx := 0; idx < t.NumField(); idx++ {
		field := t.Field(idx)
		if tomlTag, ok := field.Tag.Lookup("toml"); ok {
			tomlTag = strings.TrimSuffix(tomlTag, ",omitempty")
			if tomlTag == keys[0] {
				fmt.Printf("Field: %+v\n", field)
				findTag(keys[1:], field.Type)
			}
		}
	}

}

func findBlueprintTag(keys []string) {
	bp := blueprint.Blueprint{}
	bpt := reflect.TypeOf(bp)

	findTag(keys, bpt)
}

func cmdSet(cmd *cobra.Command, args []string) error {
	cfgKeys := strings.Split(args[0], ".")
	findBlueprintTag(cfgKeys)
	return nil
}
