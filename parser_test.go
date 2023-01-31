package graphql_test

import (
	"testing"

	"github.com/KarpelesLab/graphql"
)

func TestParser(t *testing.T) {
	v := `    query IntrospectionQuery {
      __schema {
        
        queryType { name }
        mutationType { name }
        
        types {
          ...FullType
        }
        directives {
          name
          description
          
          locations
          args {
            ...InputValue
          }
        }
      }
    }

    fragment FullType on __Type {
      kind
      name
      description
      
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }

    fragment InputValue on __InputValue {
      name
      description
      type { ...TypeRef }
      defaultValue
      
      
    }

    fragment TypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
`

	doc, err := graphql.Parse(v)
	if err != nil {
		t.Errorf("parse error: %s", err)
	}

	_ = doc
	//res, err := json.MarshalIndent(doc, "", "  ")
	//log.Printf("GOT DOCUMENT:\n%s", res)
}
