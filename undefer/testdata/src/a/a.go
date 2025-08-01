package a

func toSpecial() (err error) {
	defer println(1)
	defer println(err)          // want `defer captures current value of named result 'err'`
	defer println(1, err)       // want `defer captures current value of named result 'err'`
	defer println(1, err, &err) // want `defer captures current value of named result 'err'`
	defer println(1, &err)

	e := mayErr()
	defer println(e)
	defer println(1, e)
	return
}

func toExact() (err error) {
	defer takesErr(nil)
	defer takesErr(err) // want `defer captures current value of named result 'err'`
	defer takesValErr(1, nil)
	defer takesValErr(1, err) // want `defer captures current value of named result 'err'`

	e := mayErr()
	defer takesErr(nil)
	defer takesErr(e)
	defer takesValErr(1, nil)
	defer takesValErr(1, e)
	return
}

func toVarExact() (err error) {
	defer takesVarErr(nil)
	defer takesVarErr(err)      // want `defer captures current value of named result 'err'`
	defer takesVarErr(err, nil) // want `defer captures current value of named result 'err'`
	defer takesVarErr(nil, err) // want `defer captures current value of named result 'err'`
	defer takesValVarErr(1, nil)
	defer takesValVarErr(1, err)      // want `defer captures current value of named result 'err'`
	defer takesValVarErr(1, err, nil) // want `defer captures current value of named result 'err'`
	defer takesValVarErr(1, nil, err) // want `defer captures current value of named result 'err'`

	e := mayErr()
	defer takesVarErr(nil)
	defer takesVarErr(e)
	defer takesVarErr(e, nil)
	defer takesVarErr(nil, e)
	defer takesValVarErr(1, nil)
	defer takesValVarErr(1, e)
	defer takesValVarErr(1, e, nil)
	defer takesValVarErr(1, nil, e)
	return
}

func toAny() (err error) {
	defer takesAny(nil)
	defer takesAny(err) // want `defer captures current value of named result 'err'`
	defer takesAny(&err)
	defer takesValAny(1, nil)
	defer takesValAny(1, err) // want `defer captures current value of named result 'err'`

	e := mayErr()
	defer takesAny(nil)
	defer takesAny(e)
	defer takesValAny(1, nil)
	defer takesValAny(1, e)
	return
}

func toVarAny() (err error) {
	defer takesVarAny(nil)
	defer takesVarAny(err) // want `defer captures current value of named result 'err'`
	defer takesVarAny(&err)
	defer takesVarAny(err, nil) // want `defer captures current value of named result 'err'`
	defer takesVarAny(nil, err) // want `defer captures current value of named result 'err'`
	defer takesValVarAny(1, nil)
	defer takesValVarAny(1, err)       // want `defer captures current value of named result 'err'`
	defer takesValVarAny(1, err, nil)  // want `defer captures current value of named result 'err'`
	defer takesValVarAny(1, nil, err)  // want `defer captures current value of named result 'err'`
	defer takesValVarAny(1, &err, err) // want `defer captures current value of named result 'err'`
	defer takesValVarAny(1, err, &err) // want `defer captures current value of named result 'err'`
	defer takesValVarAny(1, &err, &err)

	e := mayErr()
	defer takesVarAny(nil)
	defer takesVarAny(e)
	defer takesVarAny(e, nil)
	defer takesVarAny(nil, e)
	defer takesValVarAny(1, nil)
	defer takesValVarAny(1, e)
	defer takesValVarAny(1, e, nil)
	defer takesValVarAny(1, nil, e)
	return
}

func toAnyNested() (err error) {
	func() (err error) {
		defer takesAny(err) // want `defer captures current value of named result 'err'`
		return
	}()
	func() (noerr error) {
		defer takesAny(err) // want `defer captures current value of named result 'err'`
		return
	}()
	return
}

var globalErr error

func mayErr() error        { return globalErr }
func takesErr(error)       {}
func takesVarErr(...error) {}
func takesAny(any)         {}
func takesVarAny(...any)   {}

func takesValErr(int, error)       {}
func takesValVarErr(int, ...error) {}
func takesValAny(int, any)         {}
func takesValVarAny(int, ...any)   {}

var _ = func() {
	toSpecial()
	toExact()
	toVarExact()
	toAny()
	toVarAny()
	toAnyNested()
}
