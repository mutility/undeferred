package noshadow

// noshadow tests with ReferencedShadows=false, omitting several reports.

func nesting() (err error) {
	func() (err error) {
		defer takesErr(err) // want `defer captures current value of named result 'err'`
		return
	}()
	func() (noerr error) {
		defer takesErr(err) // want `defer captures current value of named result 'err'`
		return
	}()

	defer func() error {
		return err
	}()

	if err := mayErr(); err != nil { // omit `shadows named result 'err' referenced in later defer`
		defer takesErr(err)

		defer func() error {
			return err // omit `defer references shadow of named result 'err'`
		}()

		defer func() error {
			if err := mayErr(); err != nil { // in-defer shadows are okay
				takesErr(err)
				defer takesErr(err) // even if deferred
				return err
			}
			// outside the inner shadow, the outer shadow is yet again flagged
			return err // omit `defer references shadow of named result 'err'`
		}()

		defer func(err error) error { // this is an in-defer shadow
			defer takesErr(err)
			return err
		}(err)

		defer func() (err error) { // this is too, but it's a new named return
			defer takesErr(err) // want `defer captures current value of named result 'err'`
			return err
		}()
	}

	defer func() {
		err := mayErr()
		takesErr(err)
	}()

	return
}

var globalErr error

func mayErr() error  { return globalErr }
func takesErr(error) {}

var _ = func() {
	nesting()
}
