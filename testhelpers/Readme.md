# Test helpers

Consists of functions helping in writing tests

## assertion.ErrorWithMessage

Simplifies writing table tests.

It implements testify's `assert.ErrorAssertionFunc` 
and checks if error is not nil and has expected message, not bothering about underlying error types.

You can use it like this

```go
func Test_MyFunc(t *testing.T) {
 type args struct {
  id int
 }

 tests := []struct {
  name    string
  args    args
  want    int64
  wantErr assert.ErrorAssertionFunc
 }{
  {
   name: "happy path",
   args: args{
    id: 1,
   },
   want:    1632,
   wantErr: assert.NoError,
  },
  {
   name: "unexpected id - error",
   args: args{
    id: -1,
   },
   want:    0,
   wantErr: assertion.ErrorWithMessage("-1 is not expected"),
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   got, err := MyFunc(tt.args.id)

   require.True(t, tt.wantErr(t, err))
   assert.Equal(t, tt.want, got)
  })
 }
}
```

## GenerateSlicePermutationsForTests

Generates all possible combinations of a struct permutating all slices.

It could help while working with uber's gomock library,
the [gomock.AnyOf](https://github.com/uber-go/mock/blob/main/gomock/matchers.go#L360) function in particular.

So, if you have a structure like

```go
type top3 struct{ Field3 []int }

type My struct {
		Field1 string
		Field2 []string
		Top3 top3
	}

myStruct := My{
	Field1: "field 1",
	Field2: []string{"one", "two", "three"},
	Top3: top3{
		Field3: []int{1, 2, 3},
	},
}
```

you can use it like

```go
myMock.EXPECT().
	myMethod(gomock.Any(GenerateSlicePermutationsForTests(myStruct)...))
```

and it will generate a slice of your structs with all possible combinations of inner slices for you.
