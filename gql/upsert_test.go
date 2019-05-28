package gql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMutationTxnBlock1(t *testing.T) {
	query := `
	query {
		me(func: eq(age, 34)) {
			uid
			friend {
				uid
				age
			}
		}
	}
`
	_, _, err := ParseMutation(query)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Invalid block: [query]")
}

func TestMutationTxnBlock2(t *testing.T) {
	query := `
	upsert {
		query {
			me(func: eq(age, 34)) {
				uid
				friend {
					uid
					age
				}
			}
		}
	}
}
`
	_, _, err := ParseMutation(query)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Too many right curl")
}

// Is this okay?
//  - Doesn't contain mutation op inside upsert block
//  - uid and age are in the same line
func TestMutationTxnBlock3(t *testing.T) {
	query := `
	upsert {
		query {
			me(func: eq(age, 34)) {
				uid
				friend {
					uid age
				}
			}
		}
	}
`
	_, _, err := ParseMutation(query)
	require.Nil(t, err)
}

func TestMutationTxnBlock4(t *testing.T) {
	query := `
	upsert {
		query {
			me(func: eq(age, 34)) {
				uid
				friend {
					uid
					age
				}
			}
		}

		mutation {
			set {
				"_:user1" <age> "45" .
			}
		}
	}
`
	_, _, err := ParseMutation(query)
	require.Nil(t, err)
}

func TestMutationTxnBlock5(t *testing.T) {
	query := `
	upsert {
		mutation {
			set {
				"_:user1" <age> "45" .
			}
		}

		query {
			me(func: eq(age, 34)) {
				uid
				friend {
					uid
					age
				}
			}
		}
	}
`
	_, _, err := ParseMutation(query)
	require.Nil(t, err)
}

// Is this okay?
func TestMutationTxnBlock6(t *testing.T) {
	query := `
	upsert {
		mutation {
			set {
				"_:user1" <age> "45" .
			}
		}

		query {
			me(func: eq(age, 34)) {
				uid
				friend {
					uid
					age
				}
			}
		}

		query {
			me2(func: eq(age, 34)) {
				uid
				friend {
					uid
					age
				}
			}
		}
	}
`
	_, _, err := ParseMutation(query)
	require.Nil(t, err)
}

// This is definitely not okay!
func TestMutationTxnBlock7(t *testing.T) {
	query := `upsert {}`
	_, _, err := ParseMutation(query)
	require.Nil(t, err)
}

func TestMutationTxnBlock8(t *testing.T) {
	query := `upsert {`
	_, _, err := ParseMutation(query)
	require.Contains(t, err.Error(), "Unclosed upsert block")
}

// Is this okay?
func TestMutationTxnBlock9(t *testing.T) {
	query := `
	upsert {
		mutation {
			set {
				"_:user1" <age> "45" .
			}
		}

		query {
			me(func: eq(age, 34)) {
				...fragmentA
				friend {
					...fragmentA
					age
				}
			}
		}

		fragment fragmentA {
			uid
		}
	}
`
	_, _, err := ParseMutation(query)
	require.Nil(t, err)
}
