package kcheck

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

var funcs MapFuncs

var TAG = "chk"

var ErrorKCHECK = errors.New("error unexpected")

func init() {
	funcs = make(MapFuncs)
	register()
}
func register() {
	funcs["nonil"] = noNilFunc
	funcs["nosp"] = noSpacesStartAndEnd
	funcs["stxt"] = sTextFunc
	funcs["email"] = emailFunc
	funcs["num"] = numberFunc
	funcs["decimal"] = decimalFunc
	funcs["len"] = lenghtFunc
	funcs["max"] = maxLenghtFunc
	funcs["min"] = minLenghtFunc
	funcs["rgx"] = regularExpression
}

// OmitFields lista de campos que no se tomaran en cuanta al realizar la verificación
type OmitFields []string

/*
	AddFunc permite registrar una nueva función personalizada, la cual será asociada
	al tagKey indicado, si le takKey ya existe, este será remplazado, por ejemplo si
	se usa el tagKey `num` este remplaza al existente. La función ValidFunc recibe
	como primer parámetro un objeto con los datos del campo a verificar, incluye el
	nombre y el valor, y como segundo parámetro recibe el valor después del `=` del
	tagKey, por ejemplo el tag es `len` y este recibe un valor `len=10` el 10 es enviado
	como segundo parámetro en formato string.
	Nota: importante que el registro de nuevas funciona no se haga en tiempo de ejecución,
	esto podría generar problemas de acceso por parte de las gorutines

*/
func AddFunc(tagKey string, f ValidFunc) {
	funcs[tagKey] = f
}
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
			atom := Atom{Name: SplitCamelCase(rsf.Name), Value: rv.String()}
			if err := ValidTarget(tagValues, atom); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidTarget(tags string, atom Atom) error {
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
