package fetch

func closer(f func() error, in *error) {
	cerr := f()
	if *in == nil {
		*in = cerr
	}
}
