package optparse

import (
	"reflect"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// Registry for special non-primitive types.
type getter func(d discord.SlashCommandInteractionData, name string) (any, bool)

var typeGetters = map[reflect.Type]getter{
	reflect.TypeFor[discord.ResolvedMember](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		return d.OptMember(name)
	},
	reflect.TypeFor[discord.Member](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		v, ok := d.OptMember(name)
		return v.Member, ok
	},
	reflect.TypeFor[discord.User](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		return d.OptUser(name)
	},

	reflect.TypeFor[discord.ResolvedChannel](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		return d.OptChannel(name)
	},

	reflect.TypeFor[discord.Role](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		return d.OptRole(name)
	},
	reflect.TypeFor[discord.Attachment](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		return d.OptAttachment(name)
	},
	reflect.TypeFor[snowflake.ID](): func(d discord.SlashCommandInteractionData, name string) (any, bool) {
		return d.OptSnowflake(name)
	},
}

func ParseCommandArguments(data discord.SlashCommandInteractionData, ref any) {
	args := ref

	t := reflect.TypeOf(args)
	ps := reflect.ValueOf(args)
	s := ps.Elem()

	if s.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < s.NumField(); i++ {
		tfield := t.Field(i)
		field := s.Field(i)

		if !field.IsValid() || !field.CanSet() {
			continue
		}

		tag := tfield.Tag.Get("disgo")
		if tag == "" {
			// Fallback to field name if you want:
			tag = strings.ToLower(tfield.Name)
		} else if tag == "-" {
			continue
		}

		baseType, ptrDepth := derefType(tfield.Type)

		// 1) Try a registered getter for special types
		if g, ok := typeGetters[baseType]; ok {
			if v, ok := g(data, tag); ok {
				setValue(field, v, ptrDepth)
			}
			continue
		}

		switch baseType.Kind() {
		case reflect.String:
			// Adjust to your library’s API if different
			v := data.String(tag)
			setValue(field, v, ptrDepth)

		case reflect.Bool:
			v := data.Bool(tag)
			setValue(field, v, ptrDepth)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// If your API exposes Number(tag) float64, convert appropriately.
			// Prefer Int(tag) if available.
			vi := int64(data.Int(tag)) // replace with your API
			setInt(field, vi, ptrDepth)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			vu := uint64(data.Int(tag)) // or convert from Number(tag) if needed
			setUint(field, vu, ptrDepth)

		case reflect.Float32, reflect.Float64:
			vf := data.Float(tag) // or Float(tag) if available
			setFloat(field, vf, ptrDepth)

		default:
			// Not supported type.
		}
	}
}
