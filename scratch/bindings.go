// scratch defines bindings to builtin scratch functions
package scratch

func Cast[T any](v any) T {
	return v.(T)
}

// Looks
func Say(text string)
func SayFor(text string, time float64)
func Think(text string)
func ThinkFor(text string, time float64)
func Show()
func Hide()

// Control
func Wait(time float64)

// Lists
func Append[T any](arr []T, v T)          {}
func Remove[T any](arr []T, ind int)      {}
func Insert[T any](arr []T, ind int, v T) {}
func Clear[T any](arr []T)                {}
func Find[T any](arr []T, v T)            {}
func Contains[T any](arr []T, v T) bool   { return false }
