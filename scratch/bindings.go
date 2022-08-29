// scratch defines bindings to builtin scratch functions
package scratch

func Say(text string)
func SayFor(text string, time float64)

func Cast[T any](v any) T {
	return v.(T)
}

func Wait(time float64)

func Clear[T any](arr []T) {}
