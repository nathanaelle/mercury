package	main

import	(
	"fmt"
	"errors"
	"os"
	"reflect"
)


const	(
	Err_Drv		int = iota
	Err_Missing_Key
	Err_Key_Wrong_Type
)


type MercuryErr struct {
	EType	int

	Driver	string
	Key	string
	File	string
	ExpType	string

	Err	error
}


func (e MercuryErr)Error()string  {
	err	:= ""
	switch e.EType {
		case	Err_Drv:
			err	= fmt.Sprintf("Unknown driver ! driver: %s", e.Driver )

		case	Err_Missing_Key:
			err	= fmt.Sprintf("key required in conf file ! file: %s key: %s", e.File, e.Key )

		case	Err_Key_Wrong_Type:
			err	= fmt.Sprintf("type required for key in conf file ! file: %s key: %s type: %s", e.File, e.Key, e.ExpType )
	}
	return	err + e.Err.Error()
}


func unknown_driver(drv string) *MercuryErr {
	return &MercuryErr { Err_Drv, drv, "", "", "", errors.New("") }
}


func error_missing_key(file string, key string) *MercuryErr {
	return &MercuryErr { Err_Missing_Key, "", key, file, "", errors.New("") }
}


func error_key_type(file string, key string, t string) *MercuryErr {
	return &MercuryErr { Err_Key_Wrong_Type, "", key, file, t, errors.New("") }
}


func exterminate(err error)  {
	var s reflect.Value

	if err == nil {
		return
	}

	s_t	:= reflect.ValueOf(err)

	for  s_t.Kind() == reflect.Ptr {
		s_t = s_t.Elem()
	}

	switch s_t.Kind() {
		case reflect.Interface:	s = s_t.Elem()
		default:		s = s_t
	}

	typeOfT := s.Type()
	pkg	:= typeOfT.PkgPath() + "/" + typeOfT.Name()

	fmt.Printf("\n------------------------------------\nKind : %d %d\n%s\n\n", s_t.Kind(), s.Kind(), err.Error())

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.CanInterface() {
			fmt.Printf("%s %d: %s %s = %v\n", pkg, i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		} else {
			fmt.Printf("%s %d: %s %s = %s\n", pkg, i, typeOfT.Field(i).Name, f.Type(), f.String())
		}
	}

	os.Exit(500)

	//syscall.Kill(syscall.Getpid(),syscall.SIGTERM)
}
