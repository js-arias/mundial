// Copyright (c) 2022 J. Salvador Arias <jsalarias@gmail.com>.
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Probs calcula el resultado esperado
// basado en una distribución fija de goles
// entre dos equipos
// y una distribución de Poisson.
package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

var promedioGol float64
var simulaciones int

func main() {
	flag.Float64Var(&promedioGol, "gol", 2.6, "goles en promedio")
	flag.IntVar(&simulaciones, "sims", 1_000_000, "número de simulaciones")
	flag.Parse()

	// los valores son de 0.1 en 0.1
	// así que calculamos el mínimo valor para ese rango
	min := math.Round(10*promedioGol/2) / 10

	fmt.Printf("# %d simulaciones\n# promedio de gol: %.3f\n", simulaciones, promedioGol)
	for i := min; i <= promedioGol; i += 0.1 {
		p := victorias(i, promedioGol-i)
		fmt.Printf("%.1f: %.3f\n", i, p)
	}
}

func victorias(e1, e2 float64) float64 {
	p1 := distuv.Poisson{Lambda: e1}
	p2 := distuv.Poisson{Lambda: e2}

	var v float64
	for s := 0; s < simulaciones; s++ {
		g1 := int(p1.Rand())
		g2 := int(p2.Rand())

		// un empate cuenta como "media" victoria
		if g1 == g2 {
			v += 0.5
			continue
		}
		if g1 > g2 {
			v++
		}
	}
	return v / float64(simulaciones)
}
