package graphql

type Token interface {
	String() string
}

const ()

type nameToken string

func (n nameToken) String() string {
	return string(n)
}
