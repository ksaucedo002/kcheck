# KCheck 
KCheck es una librería que permite validar campos del tipo `string` en una estructura.

## Guide
### Installation
```bach
    go get github.com/ksaucedo002/kcheck
```

Example
```go 
    type Persona struct {
        Dni             string `chk:"num len=8"`
        Nombre          string `chk:"max=20"`
        ApellidoPaterno string `chk:"min=10 max=80"`
        Edad            uint
    }

    func main() {
        p := Persona{Dni: "98736712", Nombre: "Kevin", ApellidoPaterno: "Saucedo", Edad: 23}
        //return err if `p` is invalid
        if err := kcheck.Valid(p); err != nil {
            log.Println(err)
        }
        //valid with omits, Field Nombre is omited
        if err := kcheck.ValidWithOmit(p, kcheck.OmitFields{"Nombre"}); err != nil {
                log.Println(err)
        }
    }
```

| tag keys  | Descripcion|
| ---- | ---------- |
|`nonil`| la cadena no puede ser nula o vacia, y los espacios vacíos no cuentas|
| `nosp` | verifica que no hay espacios vacíos al inicio o al final de la cadena|
| `stxt`| stext, safe text; solo permite texto sin espacios al inicio ni al final, las palabras no pueden estar separadas por más de dos espacios, y solo están permitidos los caracteres que no están en esta lista:`!\"#$%&'()*+,./:;<=>?@[\\]^_}{~\|`, puede usarse para recibir nombre, apellidos entre otros
| `email` | correo electronico |
| `url` | web url, http y https //aun no disponible|
| `num` | verifica si todos los caracteres son numéricos |
| `decimal` | solo numeros decimales; `00.00`, `1.00` si son validos |
| `len=:numbre` | la longitud de la cadena ingresada, debe ser igual a `:number` |
| `max=:numbre` | la longitud maxima debe ser menor igual `:number` |
| `min=:numbre` | la longitud minima debe ser mayor igual `:number` |
| `rgx=:expression` |  permite pasar un expresión regular, ejem tagKey: `rgx=(^[0-9]{8}$)`|
--------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------
Tag example: `chk:"nonil num len=8"` `//la cadena no puede estar vacía, solo puede haber 8 caracteres numéricos`