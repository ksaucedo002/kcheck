package kcheck

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/ksaucedo002/kcheck/pkg/functions"
	"github.com/ksaucedo002/kcheck/pkg/models"
)

var funcs functions.MapFuncs

const TAG = "chk"

var ErrorKCHECK = errors.New("error unexpected")

func init() {
	funcs = make(functions.MapFuncs)
	register()
}
func register() {
	funcs["num"] = functions.Number
	funcs["len"] = functions.Lenght
	funcs["max"] = functions.MaxLenght
	funcs["min"] = functions.MinLenght
}

// OmitFields lista de campos que no se tomaran en cuanta al realizar la verificación
type OmitFields []string

func (of *OmitFields) isBan(field string) bool {
	for _, v := range *of {
		if v == field {
			return true
		}
	}
	return false
}

func Valid(i interface{}) error {
	return ValidWithOmit(i, OmitFields{})
}
func ValidWithOmit(i interface{}, skips OmitFields) error {
	var rValue reflect.Value
	rType := reflect.TypeOf(i)
	if rType == nil {
		log.Println("ERROR: nil value was received")
		return ErrorKCHECK
	}
	switch rType.Kind() {
	case reflect.Struct:
		rValue = reflect.ValueOf(i)
	case reflect.Ptr:
		if rType.Elem().Kind() == reflect.Struct {
			rValue = reflect.ValueOf(i).Elem()
			rType = rType.Elem()
		} else {
			log.Printf("ERROR: a structure was type expected, invalid type `%v`\n", rType)
			return ErrorKCHECK
		}
	}
	for i := 0; i < rType.NumField(); i++ {
		rsf := rType.Field(i)
		rv := rValue.Field(i)
		if rsf.Type.Kind() == reflect.String {
			tagValues := rsf.Tag.Get(TAG)
			if tagValues == "" || skips.isBan(rsf.Name) {
				continue
			}
			atom := models.Atom{Name: SplitCamelCase(rsf.Name), Value: rv.String()}
			if err := ValidTarget(tagValues, atom); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidTarget(tags string, atom models.Atom) error {
	tags = StandardSpace(tags)
	keys := strings.Split(tags, " ")
	for _, key := range keys {
		if f, ok := funcs[key]; ok {
			if err := f(atom, ""); err != nil {
				return err
			}
		} else {
			valid, fkey, keyValues := SplitKeyValue(key)
			if valid {
				if ff, okk := funcs[fkey]; okk {
					if err := ff(atom, keyValues); err != nil {
						return err
					}
				} else {
					log.Printf("ERROR: tag value `%s` invalid in `%s` field\n", key, atom.Name)
					return ErrorKCHECK
				}
			} else {
				log.Printf("ERROR: tag value `%s` invalid in `%s` field\n", key, atom.Name)
				return ErrorKCHECK
			}
		}
	}
	return nil
}
