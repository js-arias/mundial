// Copyright (c) 2022 J. Salvador Arias <jsalarias@gmail.com>.
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Partido simula un partido particular.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/rand"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/stat/distuv"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

// ELO. Fuente: <https://www.eloratings.net/2022_World_Cup>.
var elo = map[string]int{
	"Alemania":       1919,
	"Arabia Saudita": 1692,
	"Argentina":      2086,
	"Australia":      1702,
	"Bélgica":        2020,
	"Brasil":         2169,
	"Canadá":         1763,
	"Camerún":        1610,
	"Catar":          1642 + 100, // extra por ser local
	"Costa Rica":     1723,
	"Corea del Sur":  1786,
	"Croacia":        1914,
	"Dinamarca":      1952,
	"Ecuador":        1871,
	"España":         2068,
	"Estados Unidos": 1797,
	"Francia":        2022,
	"Gales":          1791,
	"Ghana":          1567,
	"Inglaterra":     1957,
	"Irán":           1760,
	"Japón":          1831,
	"Marruecos":      1779,
	"México":         1809,
	"Países Bajos":   2050,
	"Polonia":        1814,
	"Portugal":       2006,
	"Senegal":        1677,
	"Serbia":         1898,
	"Suiza":          1902,
	"Túnez":          1726,
	"Uruguay":        1936,
}

// Probabilidad de victoria
var probsMax = []float64{
	0.500,
	0.547,
	0.594,
	0.638,
	0.682,
	0.725,
	0.764,
	0.802,
	0.837,
	0.869,
	0.897,
	0.922,
	0.944,
}

// Expectativa de goles
var goles = []float64{
	1.3,
	1.4,
	1.5,
	1.6,
	1.7,
	1.8,
	1.9,
	2.0,
	2.1,
	2.2,
	2.3,
	2.4,
	2.5,
}

const promedioGoles = 2.6

// Partidos retorna el número de goles
// entre dos equipos
// dado sus valores de ELO.
func partido(e1, e2, min int) (g1, g2 int) {
	if e1 < e2 {
		g2, g1 = partido(e2, e1, min)
		return g1, g2
	}

	dif := float64(e1-e2) / 400
	exp := 1 / (math.Pow(10, -dif) + 1)

	i, _ := slices.BinarySearch(probsMax, exp)
	if i >= len(goles) {
		i = len(goles) - 1
	}

	t := float64(90-min) / 90
	exp1 := distuv.Poisson{Lambda: goles[i] * t}
	exp2 := distuv.Poisson{Lambda: (promedioGoles - goles[i]) * t}
	g1 = int(exp1.Rand())
	g2 = int(exp2.Rand())
	return g1, g2
}

// Extra es el tiempo suplementario.
func extra(e1, e2, min int) (g1, g2 int) {
	if e1 < e2 {
		g2, g1 = extra(e2, e1, min)
		return g1, g2
	}

	dif := float64(e1-e2) / 400
	exp := 1 / (math.Pow(10, -dif) + 1)

	i, _ := slices.BinarySearch(probsMax, exp)
	if i >= len(goles) {
		i = len(goles) - 1
	}
	t := float64(30-min) / 30
	exp1 := distuv.Poisson{Lambda: t * goles[i] / 3}
	exp2 := distuv.Poisson{Lambda: t * (promedioGoles - goles[i]) / 3}

	g1 = int(exp1.Rand())
	g2 = int(exp2.Rand())
	return g1, g2
}

func main() {
	var eloFlag bool
	var tiempoSup bool
	var frecFlag bool
	var minuto int
	var simulaciones int
	var veroRes string
	flag.BoolVar(&eloFlag, "elo", false, "usa los valor de elo indicados")
	flag.BoolVar(&tiempoSup, "sup", false, "usa tiempo suplementario para resolver empates")
	flag.BoolVar(&frecFlag, "frec", false, "imprime la frecuencia de los resultados")
	flag.IntVar(&minuto, "min", 0, "tiempo de juego")
	flag.IntVar(&simulaciones, "sims", 1_000_000, "simulaciones")
	flag.StringVar(&veroRes, "vero", "", "verosimilitud de un resultado")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "esperando nombre de los equipos")
	}

	var e1, e2 int
	if eloFlag {
		var err error
		e1, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "usando --elo: argumento %q: %v\n", args[0], err)
			os.Exit(1)
		}
		e2, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "usando --elo: argumento %q: %v\n", args[1], err)
			os.Exit(1)
		}
	} else {
		var ok bool
		e1, ok = elo[args[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "país %q no reconocido\n", args[0])
			os.Exit(1)
		}
		e2, ok = elo[args[1]]
		if !ok {
			fmt.Fprintf(os.Stderr, "país %q no reconocido\n", args[1])
			os.Exit(1)
		}
	}

	var m1, m2 int
	if minuto > 0 {
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "esperando marcador de forma \"2-1\"\n")
			os.Exit(1)
		}

		if minuto > 90 && !tiempoSup {
			fmt.Fprintf(os.Stderr, "minuto %d sin suplementarios. Usar --sup\n", minuto)
			os.Exit(1)
		}
		if minuto > 120 {
			fmt.Fprintf(os.Stderr, "minuto %d invalido", minuto)
			os.Exit(1)
		}

		var err error
		vs := strings.Split(args[2], "-")
		if len(vs) != 2 {
			fmt.Fprintf(os.Stderr, "formato de marcador no reconocido: %q", args[2])
		}
		m1, err = strconv.Atoi(vs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error al leer marcador %q: %v", args[2], err)
		}
		m2, err = strconv.Atoi(vs[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error al leer marcador %q: %v", args[2], err)
		}
	}

	var v1, v2 int
	if veroRes != "" {
		var err error
		vs := strings.Split(veroRes, "-")
		if len(vs) != 2 {
			fmt.Fprintf(os.Stderr, "formato de marcador no reconocido: %q", args[2])
		}
		v1, err = strconv.Atoi(vs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error al leer marcador %q: %v", args[2], err)
		}
		v2, err = strconv.Atoi(vs[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error al leer marcador %q: %v", args[2], err)
		}
	}

	frecs := make(map[string]float64)
	var v, e int
	var mas, menos int

	for i := 0; i < simulaciones; i++ {
		g1, g2 := m1, m2
		if minuto < 90 {
			x1, x2 := partido(e1, e2, minuto)
			g1 += x1
			g2 += x2
		}

		if tiempoSup {
			mSup := 0
			if minuto > 90 {
				mSup = minuto - 90
			}
			if minuto > 90 || g1 == g2 {
				x1, x2 := extra(e1, e2, mSup)
				g1 += x1
				g2 += x2
			}
		}

		if veroRes != "" {
			if g1 == v1 && g2 == v2 {
				v++
			}
			continue
		}

		if g1 > g2 {
			v++
		}
		if g1 == g2 {
			e++
		}
		mas += g1
		menos += g2

		marcador := fmt.Sprintf("%d-%d", g1, g2)
		frecs[marcador]++
	}

	sims := float64(simulaciones)
	if veroRes != "" {
		vero := float64(v) / sims
		fmt.Printf("%s - %s: verosimilitud del resultado %q: %.6f log: %.6f\n", args[0], args[1], veroRes, vero, math.Log(vero))
		return
	}

	fmt.Printf("%s:\n\tvictorias = %.1f %%\n", args[0], float64(v*100)/sims)
	fmt.Printf("\tempates   = %.1f %%\n", float64(e*100)/sims)
	fmt.Printf("\tderrotas  = %.1f %%\n", (1-float64(v+e)/sims)*100)
	fmt.Printf("\tgoles     = %.1f-%.1f\n", float64(mas)/sims, float64(menos)/sims)

	if !frecFlag {
		return
	}
	marcadores := make([]string, 0, len(frecs))
	for m := range frecs {
		marcadores = append(marcadores, m)
	}

	slices.SortFunc(marcadores, func(a, b string) bool {
		if frecs[a] != frecs[b] {
			return frecs[a] > frecs[b]
		}
		return a < b
	})
	var sum float64
	for _, m := range marcadores {
		f := frecs[m] / sims
		fmt.Printf("%s\t%.1f %%\n", m, f*100)
		sum += f
		if sum > 0.95 {
			break
		}
	}
}
