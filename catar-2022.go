// Copyright (c) 2022 J. Salvador Arias <jsalarias@gmail.com>.
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Catar-2022 es un simulador del mundial de fútbol
// basado en el indicador ELO.
package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"golang.org/x/exp/rand"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/stat/distuv"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

// ELO antes de comenzar el mundial
// fuente: <https://www.eloratings.net/2022_World_Cup>.
var eloBase = map[string]int{
	"Alemania":       1963,
	"Arabia Saudita": 1635,
	"Argentina":      2143,
	"Australia":      1719,
	"Bélgica":        2007,
	"Brasil":         2169,
	"Canadá":         1776,
	"Camerún":        1610,
	"Catar":          1680 + 100, // extra por ser local
	"Costa Rica":     1743,
	"Corea del Sur":  1786,
	"Croacia":        1927,
	"Dinamarca":      1971,
	"Ecuador":        1833,
	"España":         2048,
	"Estados Unidos": 1798,
	"Francia":        2005,
	"Gales":          1790,
	"Ghana":          1567,
	"Inglaterra":     1920,
	"Irán":           1797,
	"Japón":          1787,
	"Marruecos":      1766,
	"México":         1809,
	"Países Bajos":   2040,
	"Polonia":        1814,
	"Portugal":       2006,
	"Senegal":        1687,
	"Serbia":         1898,
	"Suiza":          1902,
	"Túnez":          1707,
	"Uruguay":        1936,
}

var elo map[string]int

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
func partido(e1, e2 int) (g1, g2 int) {
	if e1 < e2 {
		g2, g1 = partido(e2, e1)
		return g1, g2
	}

	dif := float64(e1-e2) / 400
	exp := 1 / (math.Pow(10, -dif) + 1)

	i, _ := slices.BinarySearch(probsMax, exp)
	if i >= len(goles) {
		i = len(goles) - 1
	}
	exp1 := distuv.Poisson{Lambda: goles[i]}
	exp2 := distuv.Poisson{Lambda: promedioGoles - goles[i]}
	g1 = int(exp1.Rand())
	g2 = int(exp2.Rand())
	return g1, g2
}

// Extra es el tiempo suplementario.
func extra(e1, e2 int) (g1, g2 int) {
	if e1 < e2 {
		g2, g1 = extra(e2, e1)
		return g1, g2
	}

	dif := float64(e1-e2) / 400
	exp := 1 / (math.Pow(10, -dif) + 1)

	i, _ := slices.BinarySearch(probsMax, exp)
	if i >= len(goles) {
		i = len(goles) - 1
	}
	exp1 := distuv.Poisson{Lambda: goles[i] / 3}
	exp2 := distuv.Poisson{Lambda: (promedioGoles - goles[i]) / 3}
	g1 = int(exp1.Rand())
	g2 = int(exp2.Rand())
	return g1, g2
}

func cambioDePuntos(e1, e2 string, g1, g2 int) int {
	if elo[e1] < elo[e2] {
		return -cambioDePuntos(e2, e1, g2, g1)
	}

	dif := float64(elo[e1]-elo[e2]) / 400
	exp := 1 / (math.Pow(10, -dif) + 1)

	res := 0.5
	if g1 > g2 {
		res = 1
	} else if g2 > g1 {
		res = 0
	}

	difGol := g1 - g2
	if difGol < 0 {
		difGol = -difGol
	}
	G := float64(1)
	if difGol == 2 {
		G = 3.0 / 2
	} else if difGol >= 3 {
		G = (11 + float64(difGol)) / 8
	}

	// peso por copa del mundo
	peso := 60.0

	return int(math.Round(peso * G * (res - exp)))
}

// Contador guarda los resultados de un equipo
type contador struct {
	nombre string

	// posiciones
	p1   int
	p2   int
	p3   int
	p4   int
	oct  int
	crt  int
	sf   int
	f    int
	camp int

	// goles
	mas   int
	menos int

	// elo final
	elo int
}

var resultados map[string]*contador

type grupoPos struct {
	nombre string
	puntos int

	// goles
	mas   int
	menos int

	// suerte usado para los desempates
	suerte float64
}

func partidoDeGrupo(p1, p2 *grupoPos) {
	g1, g2 := partido(elo[p1.nombre], elo[p2.nombre])
	if g1 > g2 {
		p1.puntos += 3
	} else if g2 > g1 {
		p2.puntos += 3
	} else {
		p1.puntos += 1
		p2.puntos += 1
	}

	p1.mas += g1
	p1.menos += g2

	p2.mas += g2
	p2.menos += g1

	pts := cambioDePuntos(p1.nombre, p2.nombre, g1, g2)
	elo[p1.nombre] += pts
	elo[p2.nombre] -= pts
}

func partidoEliminatorio(e1, e2 string) string {
	c1 := resultados[e1]
	c2 := resultados[e2]

	g1, g2 := partido(elo[e1], elo[e2])
	if g1 == g2 {
		x1, x2 := extra(elo[e1], elo[e2])
		g1 += x1
		g2 += x2
	}

	c1.mas += g1
	c1.menos += g2
	c2.mas += g1
	c2.menos += g2

	pts := cambioDePuntos(e1, e2, g1, g2)
	elo[e1] += pts
	elo[e2] -= pts

	if g1 > g2 {
		return e1
	}
	if g2 > g1 {
		return e2
	}

	// Penales: una moneda al aire
	if rand.Float64() < 0.5 {
		return e1
	}
	return e2
}

func ordenarGrupo(pos []*grupoPos) {
	slices.SortFunc(pos, func(a, b *grupoPos) bool {
		// numero de puntos
		if a.puntos != b.puntos {
			return a.puntos > b.puntos
		}

		// diferencia de goles
		dA := a.mas - a.menos
		dB := b.mas - b.menos
		if dA != dB {
			return dA > dB
		}

		// goles anotados
		if a.mas != b.mas {
			return a.mas > b.mas
		}

		// suerte
		return a.suerte < b.suerte
	})
}

func resultadosGrupo(pos []*grupoPos) {
	for i, p := range pos {
		c := resultados[p.nombre]
		switch i {
		case 0:
			c.p1++
		case 1:
			c.p2++
		case 2:
			c.p3++
		case 3:
			c.p4++
		}
		if i < 2 {
			c.oct++
		}
		c.mas += p.mas
		c.menos += p.menos
	}
}

func grupoA() (a1, a2 string) {
	pos := []*grupoPos{
		{
			nombre: "Catar",
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Ecuador",
			puntos: 3,
			mas:    2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Senegal",
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Países Bajos",
			puntos: 3,
			mas:    2,
			suerte: rand.Float64(),
		},
	}
	elo["Catar"] = 1642 + 100
	elo["Ecuador"] = 1871
	elo["Países Bajos"] = 2050
	elo["Senegal"] = 1677

	// partidoDeGrupo(pos[0], pos[1]) // Catar vs Ecuador
	// partidoDeGrupo(pos[2], pos[3]) // Senegal vs Países Bajos
	partidoDeGrupo(pos[0], pos[2]) // Catar vs Senegal
	partidoDeGrupo(pos[1], pos[3]) // Ecuador vs Países Bajos
	partidoDeGrupo(pos[0], pos[3]) // Catar vs Países Bajos
	partidoDeGrupo(pos[1], pos[2]) // Ecuador vs Senegal

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoB() (b1, b2 string) {
	pos := []*grupoPos{
		{
			nombre: "Inglaterra",
			puntos: 3,
			mas:    6,
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Irán",
			mas:    2,
			menos:  6,
			suerte: rand.Float64(),
		},
		{
			nombre: "Estados Unidos",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Gales",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
	}
	elo["Inglaterra"] = 1957
	elo["Irán"] = 1760
	elo["Estados Unidos"] = 1797
	elo["Gales"] = 1791

	// partidoDeGrupo(pos[0], pos[1]) // Inglaterra vs Irán
	// partidoDeGrupo(pos[2], pos[3]) // Estados Unidos vs Gales
	partidoDeGrupo(pos[1], pos[3]) // Irán vs Gales
	partidoDeGrupo(pos[0], pos[2]) // Inglaterra vs Estados Unidos
	partidoDeGrupo(pos[0], pos[3]) // Inglaterra vs Gales
	partidoDeGrupo(pos[1], pos[2]) // Irán vs Estados Unidos

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoC() (c1, c2 string) {
	pos := []*grupoPos{
		{
			nombre: "Argentina",
			mas:    1,
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Arabia Saudita",
			puntos: 3,
			mas:    2,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "México",
			puntos: 1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Polonia",
			puntos: 1,
			suerte: rand.Float64(),
		},
	}
	elo["Argentina"] = 2086
	elo["Arabia Saudita"] = 1692
	elo["México"] = 1809
	elo["Polonia"] = 1814

	// partidoDeGrupo(pos[0], pos[1]) // Argentina vs Arabia Saudita
	// partidoDeGrupo(pos[2], pos[3]) // México vs Polonia
	partidoDeGrupo(pos[1], pos[3]) // Arabia Saudita vs Polonia
	partidoDeGrupo(pos[0], pos[2]) // Argentina vs México
	partidoDeGrupo(pos[0], pos[3]) // Argentina vs Polonia
	partidoDeGrupo(pos[1], pos[2]) // Arabia Saudita vs México

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoD() (d1, d2 string) {
	pos := []*grupoPos{
		{
			nombre: "Francia",
			puntos: 3,
			mas:    4,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Australia",
			mas:    1,
			menos:  4,
			suerte: rand.Float64(),
		},
		{
			nombre: "Dinamarca",
			puntos: 1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Túnez",
			puntos: 0,
			mas:    0,
			menos:  0,
			suerte: rand.Float64(),
		},
	}
	elo["Francia"] = 2022
	elo["Australia"] = 1702
	elo["Dinamarca"] = 1952
	elo["Túnez"] = 1726

	// partidoDeGrupo(pos[2], pos[3]) // Dinamarca vs Túnez
	// partidoDeGrupo(pos[0], pos[1]) // Francia vs Australia
	partidoDeGrupo(pos[1], pos[3]) // Australia vs Túnez
	partidoDeGrupo(pos[0], pos[2]) // Francia vs Dinamarca
	partidoDeGrupo(pos[0], pos[3]) // Túnez vs Francia
	partidoDeGrupo(pos[1], pos[2]) // Australia vs Dinamarca

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoE() (e1, e2 string) {
	pos := []*grupoPos{
		{
			nombre: "España",
			puntos: 3,
			mas:    7,
			suerte: rand.Float64(),
		},
		{
			nombre: "Costa Rica",
			menos:  7,
			suerte: rand.Float64(),
		},
		{
			nombre: "Alemania",
			mas:    1,
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Japón",
			puntos: 3,
			mas:    2,
			menos:  1,
			suerte: rand.Float64(),
		},
	}

	elo["España"] = 2068
	elo["Costa Rica"] = 1723
	elo["Alemania"] = 1919
	elo["Japón"] = 1831

	// partidoDeGrupo(pos[2], pos[3]) // Alemania vs Japón
	// partidoDeGrupo(pos[0], pos[1]) // España vs Costa Rica
	partidoDeGrupo(pos[1], pos[3]) // Costa Rica vs Japón
	partidoDeGrupo(pos[0], pos[2]) // España vs Alemania
	partidoDeGrupo(pos[0], pos[3]) // España vs Japón
	partidoDeGrupo(pos[1], pos[2]) // Costa Rica vs Alemania

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoF() (f1, f2 string) {
	pos := []*grupoPos{
		{
			nombre: "Bélgica",
			puntos: 3,
			mas:    1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Canadá",
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Marruecos",
			puntos: 1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Croacia",
			puntos: 1,
			suerte: rand.Float64(),
		},
	}

	elo["Bélgica"] = 2020
	elo["Canadá"] = 1763
	elo["Marruecos"] = 1779
	elo["Croacia"] = 1914

	// partidoDeGrupo(pos[2], pos[3]) // Marruecos vs Croacia
	// partidoDeGrupo(pos[0], pos[1]) // Bélgica vs Canadá
	partidoDeGrupo(pos[0], pos[2]) // Bélgica vs Marruecos
	partidoDeGrupo(pos[1], pos[3]) // Canadá vs Croacia
	partidoDeGrupo(pos[0], pos[3]) // Bélgica vs Croacia
	partidoDeGrupo(pos[1], pos[2]) // Canadá vs Marruecos

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoG() (g1, g2 string) {
	pos := []*grupoPos{
		{
			nombre: "Brasil",
			puntos: 3,
			mas:    2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Serbia",
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Suiza",
			puntos: 3,
			mas:    1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Camerún",
			menos:  1,
			suerte: rand.Float64(),
		},
	}
	elo["Brasil"] = 2185
	elo["Serbia"] = 1882
	elo["Suiza"] = 1911
	elo["Camerún"] = 1601

	// partidoDeGrupo(pos[2], pos[3]) // Suiza vs Camerún
	// partidoDeGrupo(pos[0], pos[1]) // Brasil vs Serbia
	partidoDeGrupo(pos[1], pos[3]) // Camerún vs Serbia
	partidoDeGrupo(pos[0], pos[2]) // Brasil vs Suiza
	partidoDeGrupo(pos[0], pos[3]) // Brasil vs Camerún
	partidoDeGrupo(pos[1], pos[2]) // Serbia vs Suiza

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func grupoH() (g1, g2 string) {
	pos := []*grupoPos{
		{
			nombre: "Portugal",
			puntos: 3,
			mas:    3,
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Ghana",
			mas:    2,
			menos:  3,
			suerte: rand.Float64(),
		},
		{
			nombre: "Uruguay",
			puntos: 1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Corea del Sur",
			puntos: 1,
			suerte: rand.Float64(),
		},
	}
	elo["Portugal"] = 2010
	elo["Ghana"] = 1463
	elo["Uruguay"] = 1924
	elo["Corea del Sur"] = 1798

	// partidoDeGrupo(pos[2], pos[3]) // Uruguay vs Corea del Sur
	// partidoDeGrupo(pos[0], pos[1]) // Portugal vs Ghana
	partidoDeGrupo(pos[1], pos[3]) // Ghana vs Corea del Sur
	partidoDeGrupo(pos[0], pos[2]) // Portugal vs Uruguay
	partidoDeGrupo(pos[0], pos[3]) // Portugal vs Corea del Sur
	partidoDeGrupo(pos[1], pos[2]) // Ghana vs Uruguay

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos[0].nombre, pos[1].nombre
}

func simulacion() {
	elo = make(map[string]int, len(eloBase))
	for n, e := range eloBase {
		elo[n] = e
	}

	// fase de grupos
	a1, a2 := grupoA()
	b1, b2 := grupoB()
	c1, c2 := grupoC()
	d1, d2 := grupoD()
	e1, e2 := grupoE()
	f1, f2 := grupoF()
	g1, g2 := grupoG()
	h1, h2 := grupoH()

	// octavos
	m49 := partidoEliminatorio(a1, b2)
	m50 := partidoEliminatorio(c1, d2)
	m53 := partidoEliminatorio(e1, f2)
	m54 := partidoEliminatorio(g1, h2)
	m51 := partidoEliminatorio(b1, a2)
	m52 := partidoEliminatorio(d1, c2)
	m55 := partidoEliminatorio(f1, e2)
	m56 := partidoEliminatorio(h1, g2)

	resultados[m49].crt++
	resultados[m50].crt++
	resultados[m53].crt++
	resultados[m54].crt++
	resultados[m51].crt++
	resultados[m52].crt++
	resultados[m55].crt++
	resultados[m56].crt++

	// cuartos
	m57 := partidoEliminatorio(m49, m50)
	m58 := partidoEliminatorio(m53, m54)
	m59 := partidoEliminatorio(m51, m52)
	m60 := partidoEliminatorio(m55, m56)

	resultados[m57].sf++
	resultados[m58].sf++
	resultados[m59].sf++
	resultados[m60].sf++

	// semi-finales
	m61 := partidoEliminatorio(m57, m58)
	m62 := partidoEliminatorio(m59, m60)

	resultados[m61].f++
	resultados[m62].f++

	camp := partidoEliminatorio(m61, m62)
	resultados[camp].camp++

	for n, e := range elo {
		resultados[n].elo += e
	}
}

var simulaciones int
var outFormat string

func main() {
	flag.IntVar(&simulaciones, "sims", 1_000_000, "número de simulaciones")
	flag.StringVar(&outFormat, "fmt", "", "formato de salida, \"md\" para markdown")
	flag.Parse()

	resultados = make(map[string]*contador, len(elo))
	for nombre := range eloBase {
		c := &contador{
			nombre: nombre,
		}
		resultados[nombre] = c
	}

	for i := 0; i < simulaciones; i++ {
		simulacion()
	}

	res := make([]*contador, 0, len(resultados))
	for _, c := range resultados {
		res = append(res, c)
	}
	slices.SortFunc(res, func(a, b *contador) bool {
		if a.camp != b.camp {
			return a.camp > b.camp
		}
		if a.f != b.f {
			return a.f > b.f
		}
		if a.sf != b.sf {
			return a.sf > b.sf
		}
		if a.crt != b.crt {
			return a.crt > b.crt
		}
		if a.oct != b.oct {
			return a.oct > b.oct
		}
		if a.p3 != b.p3 {
			return a.p3 > b.p3
		}
		return eloBase[a.nombre] > eloBase[b.nombre]
	})

	sims := float64(simulaciones) / 100

	if outFormat == "md" {
		fmt.Printf("Equipo | ELO | ELO final | P1 | P2 | P3 | P4 | Ocv | Ct | Sf | Fin | Camp | Goles\n")
		fmt.Printf("------ | --- | --------- | -- | -- | -- | -- | --- | -- | -- | --- | ---- | -----\n")
		for _, c := range res {
			fmt.Printf("%s | %d | %d | ", c.nombre, eloBase[c.nombre], int((float64(c.elo))/float64(simulaciones)))
			fmt.Printf("%d | %d | %d | %d | ", int(float64(c.p1)/sims), int(float64(c.p2)/sims), int(float64(c.p3)/sims), int(float64(c.p4)/sims))
			fmt.Printf("%d | %d | %d | %d | ", int(float64(c.oct)/sims), int(float64(c.crt)/sims), int(float64(c.sf)/sims), int(float64(c.f)/sims))
			fmt.Printf("%d | ", int(float64(c.camp)/sims))
			fmt.Printf("%.1f-%.1f\n", float64(c.mas)/float64(simulaciones), float64(c.menos)/float64(simulaciones))
		}
		return
	}
	fmt.Printf("# simulaciones %d\n", simulaciones)
	fmt.Printf("Equipo\tELO -> ELO final\tP1 P2 P3 P4\tOc Ct Sf F\tCamp\tGoles\n")
	for _, c := range res {
		fmt.Printf("%s\t%d -> %d\t", c.nombre, eloBase[c.nombre], int((float64(c.elo))/float64(simulaciones)))
		fmt.Printf("%d %d %d %d\t", int(float64(c.p1)/sims), int(float64(c.p2)/sims), int(float64(c.p3)/sims), int(float64(c.p4)/sims))
		fmt.Printf("%d %d %d %d\t", int(float64(c.oct)/sims), int(float64(c.crt)/sims), int(float64(c.sf)/sims), int(float64(c.f)/sims))
		fmt.Printf("%d\t", int(float64(c.camp)/sims))
		fmt.Printf("%.1f-%.1f\n", float64(c.mas)/float64(simulaciones), float64(c.menos)/float64(simulaciones))
	}
}
