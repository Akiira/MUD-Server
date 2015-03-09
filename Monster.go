// Monster
package main
import "math/rand"
type Monster struct {
	HP int
	Name string
	Defense int
}

func (m *Monster) getAttackRoll() int {
	return rand.Int() % 6
}