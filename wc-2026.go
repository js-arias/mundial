// Copyright (c) 2022 J. Salvador Arias <jsalarias@gmail.com>.
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// wc-2022 es un simulador del mundial de fútbol
// basado en el ELO.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"slices"

	"gonum.org/v1/gonum/stat/distuv"
)

// ELO antes de comenzar el mundial
// fuente: <https://www.eloratings.net/2026_World_Cup>.
var eloBase = map[string]int{
	"Algeria":                1772,
	"Argentina":              2115,
	"Australia":              1777,
	"Austria":                1830,
	"Belgium":                1894,
	"Bosnia and Herzegovina": 1595,
	"Brazil":                 1991,
	"Canada":                 1788,
	"Cape Verde":             1578,
	"Colombia":               1982,
	"Croatia":                1912,
	"Curazao":                1434,
	"Czech Republic":         1740,
	"DR Congo":               1652,
	"Ecuador":                1938,
	"Egypt":                  1696,
	"England":                2024,
	"France":                 2063,
	"Germany":                1932,
	"Ghana":                  1510,
	"Haiti":                  1548,
	"Iran":                   1772,
	"Iraq":                   1607,
	"Ivory Coast":            1695,
	"Japan":                  1906,
	"Jordan":                 1680,
	"Mexico":                 1875,
	"Morocco":                1827,
	"Netherlands":            1948,
	"New Zealand":            1562,
	"Norway":                 1914,
	"Panama":                 1730,
	"Paraguay":               1834,
	"Portugal":               1989,
	"Qatar":                  1421,
	"Saudi Arabia":           1576,
	"Scotland":               1782,
	"Senegal":                1860,
	"South Africa":           1517,
	"South Korea":            1758,
	"Spain":                  2157,
	"Sweden":                 1712,
	"Switzerland":            1891,
	"Tunisia":                1628,
	"Turkey":                 1911,
	"United States":          1726,
	"Uruguay":                1892,
	"Uzbekistan":             1714,
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
	f16  int
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

func (c *contador) round(sims int) {
	c.p1 = int(math.Round(float64(c.p1*100) / float64(sims)))
	c.p2 = int(math.Round(float64(c.p2*100) / float64(sims)))
	c.p3 = int(math.Round(float64(c.p3*100) / float64(sims)))
	c.p4 = int(math.Round(float64(c.p4*100) / float64(sims)))
	c.f16 = int(math.Round(float64(c.f16*100) / float64(sims)))
	c.oct = int(math.Round(float64(c.oct*100) / float64(sims)))
	c.crt = int(math.Round(float64(c.crt*100) / float64(sims)))
	c.sf = int(math.Round(float64(c.sf*100) / float64(sims)))
	c.f = int(math.Round(float64(c.f*100) / float64(sims)))
	c.camp = int(math.Round(float64(c.camp*100) / float64(sims)))

	c.elo = int(math.Round(float64(c.elo) / float64(sims)))
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

func partidoDeGrupo(grupo string, p1, p2 *grupoPos) {
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

	if verbose {
		fmt.Printf("%s: %s %d - %s %d\n", grupo, p1.nombre, g1, p2.nombre, g2)
	}
}

func partidoEliminatorio(fase, e1, e2 string) string {
	c1 := resultados[e1]
	c2 := resultados[e2]

	g1, g2 := partido(elo[e1], elo[e2])

	// tiempo extra es un empate
	pts := cambioDePuntos(e1, e2, g1, g2)
	elo[e1] += pts
	elo[e2] -= pts

	ext := ""
	if g1 == g2 {
		x1, x2 := extra(elo[e1], elo[e2])
		g1 += x1
		g2 += x2
		ext = "[extra]"
	}
	c1.mas += g1
	c1.menos += g2
	c2.mas += g1
	c2.menos += g2

	if verbose {
		fmt.Printf("%s: %s %d - %s %d %s\n", fase, e1, g1, e2, g2, ext)
	}

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
	slices.SortFunc(pos, func(a, b *grupoPos) int {
		// numero de puntos
		if a.puntos != b.puntos {
			if a.puntos > b.puntos {
				return -1
			}
			return 1
		}

		// diferencia de goles
		dA := a.mas - a.menos
		dB := b.mas - b.menos
		if dA != dB {
			if dA > dB {
				return -1
			}
			return 1
		}

		// goles anotados
		if a.mas != b.mas {
			if a.mas > b.mas {
				return -1
			}
			return 1
		}

		// suerte
		if a.suerte < b.suerte {
			return -1
		}
		return 1
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
			c.f16++
		}
		c.mas += p.mas
		c.menos += p.menos
	}
}

func grupoA() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Mexico",
			puntos: 3,
			mas:    2,
			suerte: rand.Float64(),
		},
		{
			nombre: "South Africa",
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "South Korea",
			puntos: 3,
			mas:    2,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Czech Republic",
			mas:    1,
			menos:  2,
			suerte: rand.Float64(),
		},
	}
	// elo["Mexico"] = 1875
	// elo["South Africa"] = 1517
	// elo["South Korea"] = 1758
	// elo["Czech Republic"] = 1740

	elo["Mexico"] = 1881
	elo["South Africa"] = 1511
	elo["South Korea"] = 1786
	elo["Czech Republic"] = 1712

	// partidoDeGrupo("A", pos[0], pos[1]) // Mexico vs South Africa
	// partidoDeGrupo("A", pos[2], pos[3]) // South Korea vs Czech Republic
	partidoDeGrupo("A", pos[0], pos[2]) // Mexico vs South Korea
	partidoDeGrupo("A", pos[1], pos[3]) // South Africa vs Czech Republic
	partidoDeGrupo("A", pos[0], pos[3]) // Mexico vs Czech Republic
	partidoDeGrupo("A", pos[1], pos[2]) // South Africa vs South Korea

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoB() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Canada",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Bosnia and Herzegovina",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Qatar",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Switzerland",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
	}
	// elo["Canada"] = 1788
	// elo["Bosnia and Herzegovina"] = 1595
	// elo["Qatar"] = 1421
	// elo["Switzerland"] = 1891

	elo["Canada"] = 1767
	elo["Bosnia and Herzegovina"] = 1616
	elo["Qatar"] = 1447
	elo["Switzerland"] = 1865

	// partidoDeGrupo("B", pos[0], pos[1]) // Canada vs Bosnia and Herzegovina
	// partidoDeGrupo("B", pos[2], pos[3]) // Qatar vs Switzerland
	partidoDeGrupo("B", pos[0], pos[2]) // Canada vs Qatar
	partidoDeGrupo("B", pos[1], pos[3]) // Bosnia and Herzegovina vs Switzerland
	partidoDeGrupo("B", pos[0], pos[3]) // Canada vs Switzerland
	partidoDeGrupo("B", pos[1], pos[2]) // Bosnia and Herzegovina vs Qatar

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoC() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Brazil",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Morocco",
			puntos: 1,
			mas:    1,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Haiti",
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Scotland",
			puntos: 3,
			mas:    1,
			suerte: rand.Float64(),
		},
	}
	// elo["Brazil"] = 1991
	// elo["Morocco"] = 1827
	// elo["Haiti"] = 1548
	// elo["Scotland"] = 1782

	elo["Brazil"] = 1978
	elo["Morocco"] = 1840
	elo["Haiti"] = 1536
	elo["Scotland"] = 1794

	// partidoDeGrupo("C", pos[0], pos[1]) // Brazil vs Morocco
	// partidoDeGrupo("C", pos[2], pos[3]) // Haiti vs Scotland
	partidoDeGrupo("C", pos[0], pos[2]) // Brazil vs Haiti
	partidoDeGrupo("C", pos[1], pos[3]) // Morocco vs Scotland
	partidoDeGrupo("C", pos[0], pos[3]) // Brazil vs Scotland
	partidoDeGrupo("C", pos[1], pos[2]) // Morocco vs Haiti

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoD() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "United States",
			puntos: 3,
			mas:    4,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Paraguay",
			mas:    1,
			menos:  4,
			suerte: rand.Float64(),
		},
		{
			nombre: "Australia",
			puntos: 3,
			mas:    2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Turkey",
			menos:  2,
			suerte: rand.Float64(),
		},
	}
	// elo["United States"] = 1726
	// elo["Paraguay"] = 1834
	// elo["Australia"] = 1777
	// elo["Turkey"] = 1911

	elo["United States"] = 1780
	elo["Paraguay"] = 1780
	elo["Australia"] = 1839
	elo["Turkey"] = 1849

	// partidoDeGrupo("D", pos[0], pos[1]) // United States vs Paraguay
	// partidoDeGrupo("D", pos[2], pos[3]) // Australia vs Turkey
	partidoDeGrupo("D", pos[0], pos[2]) // United States vs Australia
	partidoDeGrupo("D", pos[1], pos[3]) // Paraguay vs Turkey
	partidoDeGrupo("D", pos[0], pos[3]) // United States vs Turkey
	partidoDeGrupo("D", pos[1], pos[2]) // Paraguay vs Australia

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoE() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Germany",
			puntos: 3,
			mas:    7,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Curazao",
			mas:    1,
			menos:  7,
			suerte: rand.Float64(),
		},
		{
			nombre: "Ivory Coast",
			puntos: 3,
			mas:    1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Ecuador",
			menos:  1,
			suerte: rand.Float64(),
		},
	}
	// elo["Germany"] = 1932
	// elo["Curazao"] = 1434
	// elo["Ivory Coast"] = 1695
	// elo["Ecuador"] = 1938

	elo["Germany"] = 1939
	elo["Curazao"] = 1427
	elo["Ivory Coast"] = 1743
	elo["Ecuador"] = 1890

	// partidoDeGrupo("E", pos[0], pos[1]) // Germany vs Curazao
	// partidoDeGrupo("E", pos[2], pos[3]) // Ivory Coast vs Ecuador
	partidoDeGrupo("E", pos[0], pos[2]) // Germany vs Ivory Coast
	partidoDeGrupo("E", pos[1], pos[3]) // Curazao vs Ecuador
	partidoDeGrupo("E", pos[0], pos[3]) // Germany vs Ecuador
	partidoDeGrupo("E", pos[1], pos[2]) // Curazao vs Ivory Coast

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoF() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Netherlands",
			puntos: 1,
			mas:    2,
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Japan",
			puntos: 1,
			mas:    2,
			menos:  2,
			suerte: rand.Float64(),
		},
		{
			nombre: "Sweden",
			puntos: 3,
			mas:    5,
			menos:  1,
			suerte: rand.Float64(),
		},
		{
			nombre: "Tunisia",
			mas:    1,
			menos:  5,
			suerte: rand.Float64(),
		},
	}
	// elo["Netherlands"] = 1948
	// elo["Japan"] = 1906
	// elo["Sweden"] = 1712
	// elo["Tunisia"] = 1628

	elo["Netherlands"] = 1944
	elo["Japan"] = 1910
	elo["Sweden"] = 1755
	elo["Tunisia"] = 1585

	// partidoDeGrupo("F", pos[0], pos[1]) // Netherlands vs Japan
	// partidoDeGrupo("F", pos[2], pos[3]) // Sweden vs Tunisia
	partidoDeGrupo("F", pos[0], pos[2]) // Netherlands vs Sweden
	partidoDeGrupo("F", pos[1], pos[3]) // Japan vs Tunisia
	partidoDeGrupo("F", pos[0], pos[3]) // Netherlands vs Tunisia
	partidoDeGrupo("F", pos[1], pos[2]) // Japan vs Sweden

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoG() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Belgium",
			suerte: rand.Float64(),
		},
		{
			nombre: "Egypt",
			suerte: rand.Float64(),
		},
		{
			nombre: "Iran",
			suerte: rand.Float64(),
		},
		{
			nombre: "New Zealand",
			suerte: rand.Float64(),
		},
	}
	elo["Belgium"] = 1894
	elo["Egypt"] = 1696
	elo["Iran"] = 1772
	elo["New Zealand"] = 1562

	partidoDeGrupo("G", pos[0], pos[1]) // Belgium vs Egypt
	partidoDeGrupo("G", pos[2], pos[3]) // Iran vs New Zealand
	partidoDeGrupo("G", pos[0], pos[2]) // Belgium vs Iran
	partidoDeGrupo("G", pos[1], pos[3]) // Egypt vs New Zealand
	partidoDeGrupo("G", pos[0], pos[3]) // Belgium vs New Zealand
	partidoDeGrupo("G", pos[1], pos[2]) // Egypt vs Iran

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoH() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Spain",
			suerte: rand.Float64(),
		},
		{
			nombre: "Cape Verde",
			suerte: rand.Float64(),
		},
		{
			nombre: "Saudi Arabia",
			suerte: rand.Float64(),
		},
		{
			nombre: "Uruguay",
			suerte: rand.Float64(),
		},
	}
	elo["Spain"] = 2157
	elo["Cape Verde"] = 1578
	elo["Saudi Arabia"] = 1576
	elo["Uruguay"] = 1892

	partidoDeGrupo("H", pos[0], pos[1]) // Spain vs Cape Verde
	partidoDeGrupo("H", pos[2], pos[3]) // Saudi Arabia vs Uruguay
	partidoDeGrupo("H", pos[0], pos[2]) // Spain vs Saudi Arabia
	partidoDeGrupo("H", pos[1], pos[3]) // Cape Verde vs Uruguay
	partidoDeGrupo("H", pos[0], pos[3]) // Spain vs Uruguay
	partidoDeGrupo("H", pos[1], pos[2]) // Cape Verde vs Saudi Arabia

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoI() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "France",
			suerte: rand.Float64(),
		},
		{
			nombre: "Senegal",
			suerte: rand.Float64(),
		},
		{
			nombre: "Iraq",
			suerte: rand.Float64(),
		},
		{
			nombre: "Norway",
			suerte: rand.Float64(),
		},
	}
	elo["France"] = 2063
	elo["Senegal"] = 1860
	elo["Iraq"] = 1607
	elo["Norway"] = 1914

	partidoDeGrupo("I", pos[0], pos[1]) // France vs Senegal
	partidoDeGrupo("I", pos[2], pos[3]) // Iraq vs Norway
	partidoDeGrupo("I", pos[0], pos[2]) // France vs Iraq
	partidoDeGrupo("I", pos[1], pos[3]) // Senegal vs Norway
	partidoDeGrupo("I", pos[0], pos[3]) // France vs Norway
	partidoDeGrupo("I", pos[1], pos[2]) // Senegal vs Iraq

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoJ() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Argentina",
			suerte: rand.Float64(),
		},
		{
			nombre: "Algeria",
			suerte: rand.Float64(),
		},
		{
			nombre: "Austria",
			suerte: rand.Float64(),
		},
		{
			nombre: "Jordan",
			suerte: rand.Float64(),
		},
	}
	elo["Argentina"] = 2115
	elo["Algeria"] = 1772
	elo["Austria"] = 1830
	elo["Jordan"] = 1680

	partidoDeGrupo("J", pos[0], pos[1]) // Argentina vs Algeria
	partidoDeGrupo("J", pos[2], pos[3]) // Austria vs Jordan
	partidoDeGrupo("J", pos[0], pos[2]) // Argentina vs Austria
	partidoDeGrupo("J", pos[1], pos[3]) // Algeria vs Jordan
	partidoDeGrupo("J", pos[0], pos[3]) // Argentina vs Jordan
	partidoDeGrupo("J", pos[1], pos[2]) // Algeria vs Austria

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoK() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "Portugal",
			suerte: rand.Float64(),
		},
		{
			nombre: "DR Congo",
			suerte: rand.Float64(),
		},
		{
			nombre: "Uzbekistan",
			suerte: rand.Float64(),
		},
		{
			nombre: "Colombia",
			suerte: rand.Float64(),
		},
	}
	elo["Portugal"] = 1989
	elo["DR Congo"] = 1652
	elo["Uzbekistan"] = 1714
	elo["Colombia"] = 1982

	partidoDeGrupo("K", pos[0], pos[1]) // Portugal vs DR Congo
	partidoDeGrupo("K", pos[2], pos[3]) // Uzbekistan vs Colombia
	partidoDeGrupo("K", pos[0], pos[2]) // Portugal vs Uzbekistan
	partidoDeGrupo("K", pos[1], pos[3]) // DR Congo vs Colombia
	partidoDeGrupo("K", pos[0], pos[3]) // Portugal vs Colombia
	partidoDeGrupo("K", pos[1], pos[2]) // DR Congo vs Uzbekistan

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func grupoL() []*grupoPos {
	pos := []*grupoPos{
		{
			nombre: "England",
			suerte: rand.Float64(),
		},
		{
			nombre: "Croatia",
			suerte: rand.Float64(),
		},
		{
			nombre: "Ghana",
			suerte: rand.Float64(),
		},
		{
			nombre: "Panama",
			suerte: rand.Float64(),
		},
	}
	elo["England"] = 2024
	elo["Croatia"] = 1912
	elo["Ghana"] = 1510
	elo["Panama"] = 1730

	partidoDeGrupo("L", pos[0], pos[1]) // England vs Croatia
	partidoDeGrupo("L", pos[2], pos[3]) // Ghana vs Panama
	partidoDeGrupo("L", pos[0], pos[2]) // England vs Ghana
	partidoDeGrupo("L", pos[1], pos[3]) // Croatia vs Panama
	partidoDeGrupo("L", pos[0], pos[3]) // England vs Panama
	partidoDeGrupo("L", pos[1], pos[2]) // Croatia vs Ghana

	ordenarGrupo(pos)
	resultadosGrupo(pos)

	return pos
}

func simulacion() {
	elo = make(map[string]int, len(eloBase))
	for n, e := range eloBase {
		elo[n] = e
	}

	terceros := make([]*grupoPos, 0, 12)
	// fase de grupos
	gA := grupoA()
	terceros = append(terceros, gA[2])
	gB := grupoB()
	terceros = append(terceros, gB[2])
	gC := grupoC()
	terceros = append(terceros, gC[2])
	gD := grupoD()
	terceros = append(terceros, gD[2])
	gE := grupoE()
	terceros = append(terceros, gE[2])
	gF := grupoF()
	terceros = append(terceros, gF[2])
	gG := grupoG()
	terceros = append(terceros, gG[2])
	gH := grupoH()
	terceros = append(terceros, gH[2])
	gI := grupoI()
	terceros = append(terceros, gI[2])
	gJ := grupoJ()
	terceros = append(terceros, gJ[2])
	gK := grupoK()
	terceros = append(terceros, gK[2])
	gL := grupoL()
	terceros = append(terceros, gL[2])
	ordenarGrupo(terceros)
	terceros = terceros[:8]
	for _, p := range terceros {
		c := resultados[p.nombre]
		c.f16++
	}
	rand.Shuffle(len(terceros), func(i, j int) {
		terceros[i], terceros[j] = terceros[j], terceros[i]
	})

	// 16vos de final
	m74 := partidoEliminatorio("16", gE[0].nombre, terceros[0].nombre)
	m77 := partidoEliminatorio("16", gI[0].nombre, terceros[1].nombre)
	m73 := partidoEliminatorio("16", gA[1].nombre, gB[1].nombre)
	m75 := partidoEliminatorio("16", gF[0].nombre, gC[1].nombre)
	m83 := partidoEliminatorio("16", gK[1].nombre, gL[1].nombre)
	m84 := partidoEliminatorio("16", gH[0].nombre, gJ[1].nombre)
	m81 := partidoEliminatorio("16", gD[0].nombre, terceros[2].nombre)
	m82 := partidoEliminatorio("16", gG[0].nombre, terceros[3].nombre)
	m76 := partidoEliminatorio("16", gC[0].nombre, gF[1].nombre)
	m78 := partidoEliminatorio("16", gE[1].nombre, gI[1].nombre)
	m79 := partidoEliminatorio("16", gA[0].nombre, terceros[4].nombre)
	m80 := partidoEliminatorio("16", gL[0].nombre, terceros[5].nombre)
	m86 := partidoEliminatorio("16", gJ[0].nombre, gH[1].nombre)
	m88 := partidoEliminatorio("16", gD[1].nombre, gG[1].nombre)
	m85 := partidoEliminatorio("16", gB[0].nombre, terceros[6].nombre)
	m87 := partidoEliminatorio("16", gK[0].nombre, terceros[7].nombre)

	resultados[m74].oct++
	resultados[m75].oct++
	resultados[m76].oct++
	resultados[m77].oct++
	resultados[m78].oct++
	resultados[m79].oct++
	resultados[m80].oct++
	resultados[m81].oct++
	resultados[m82].oct++
	resultados[m83].oct++
	resultados[m84].oct++
	resultados[m85].oct++
	resultados[m86].oct++
	resultados[m87].oct++
	resultados[m88].oct++

	// octavos
	m89 := partidoEliminatorio("8v", m74, m77)
	m90 := partidoEliminatorio("8v", m73, m75)
	m93 := partidoEliminatorio("8v", m83, m84)
	m94 := partidoEliminatorio("8v", m81, m82)
	m91 := partidoEliminatorio("8v", m76, m78)
	m92 := partidoEliminatorio("8v", m79, m80)
	m95 := partidoEliminatorio("8v", m86, m88)
	m96 := partidoEliminatorio("8v", m85, m87)

	resultados[m89].crt++
	resultados[m90].crt++
	resultados[m91].crt++
	resultados[m92].crt++
	resultados[m93].crt++
	resultados[m94].crt++
	resultados[m95].crt++
	resultados[m96].crt++

	// cuartos
	m97 := partidoEliminatorio("4t", m89, m90)
	m98 := partidoEliminatorio("4t", m93, m94)
	m99 := partidoEliminatorio("4t", m91, m92)
	m100 := partidoEliminatorio("4t", m95, m96)

	resultados[m97].sf++
	resultados[m98].sf++
	resultados[m99].sf++
	resultados[m100].sf++

	// semi-finales
	m101 := partidoEliminatorio("SF", m97, m98)
	m102 := partidoEliminatorio("SF", m99, m100)

	resultados[m101].f++
	resultados[m102].f++

	camp := partidoEliminatorio("Final", m101, m102)
	resultados[camp].camp++

	for n, e := range elo {
		resultados[n].elo += e
	}
}

var simulaciones int
var outFormat string
var verbose bool

func main() {
	flag.BoolVar(&verbose, "verbose", false, "imprime los resultados")
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

	args := flag.Args()
	if len(args) > 0 {
		c, ok := resultados[args[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "equipo no encontrado: %q\n", args[0])
			os.Exit(1)
		}
		fmt.Printf("%s\n", c.nombre)
		fmt.Printf("ELO: %d -> %.3f\n", eloBase[c.nombre], float64(c.elo)/float64(simulaciones))
		fmt.Printf("P1: %.6f\n", float64(c.p1)/float64(simulaciones))
		fmt.Printf("P2: %.6f\n", float64(c.p2)/float64(simulaciones))
		fmt.Printf("P3: %.6f\n", float64(c.p3)/float64(simulaciones))
		fmt.Printf("P4: %.6f\n", float64(c.p4)/float64(simulaciones))
		fmt.Printf("F16: %.6f\n", float64(c.f16)/float64(simulaciones))
		fmt.Printf("OCV: %.6f\n", float64(c.oct)/float64(simulaciones))
		fmt.Printf("CT: %.6f\n", float64(c.crt)/float64(simulaciones))
		fmt.Printf("SF: %.6f\n", float64(c.sf)/float64(simulaciones))
		fmt.Printf("F: %.6f\n", float64(c.f)/float64(simulaciones))
		fmt.Printf("Camp: %.6f\n", float64(c.camp)/float64(simulaciones))
		fmt.Printf("Goles: %.3f-%.3f\n", float64(c.mas)/float64(simulaciones), float64(c.menos)/float64(simulaciones))
		return
	}

	res := make([]*contador, 0, len(resultados))
	for _, c := range resultados {
		c.round(simulaciones)
		res = append(res, c)
	}
	slices.SortFunc(res, func(a, b *contador) int {
		if a.camp != b.camp {
			if a.camp > b.camp {
				return -1
			}
			return 1
		}
		if a.f != b.f {
			if a.f > b.f {
				return -1
			}
			return 1
		}
		if a.sf != b.sf {
			if a.sf > b.sf {
				return -1
			}
			return 1
		}
		if a.crt != b.crt {
			if a.crt > b.crt {
				return -1
			}
			return 1
		}
		if a.oct != b.oct {
			if a.oct > b.oct {
				return -1
			}
			return 1
		}
		if a.f16 != b.f16 {
			if a.f16 > b.f16 {
				return -1
			}
			return 1
		}
		if a.p3 != b.p3 {
			if a.p3 > b.p3 {
				return -1
			}
			return 1
		}
		if eloBase[a.nombre] > eloBase[b.nombre] {
			return -1
		}
		return 1
	})

	if outFormat == "md" {
		fmt.Printf("Equipo | ELO | ELO final | P1 | P2 | P3 | P4 | F16 | Ocv | Ct | Sf | Fin | Camp | Goles\n")
		fmt.Printf("------ | --- | --------- | -- | -- | -- | -- | --- | --- | -- | -- | --- | ---- | -----\n")
		for _, c := range res {
			fmt.Printf("%s | %d | %d | ", c.nombre, eloBase[c.nombre], c.elo)
			fmt.Printf("%d | %d | %d | %d | ", c.p1, c.p2, c.p3, c.p4)
			fmt.Printf("%d | %d | %d | %d | %d | ", c.f16, c.oct, c.crt, c.sf, c.f)
			fmt.Printf("%d | ", c.camp)
			fmt.Printf("%.1f-%.1f\n", float64(c.mas)/float64(simulaciones), float64(c.menos)/float64(simulaciones))
		}
		return
	}
	fmt.Printf("# simulaciones %d\n", simulaciones)
	fmt.Printf("Equipo\tELO -> ELO final\tP1 P2 P3 P4\t16 Oc Ct Sf F\tCamp\tGoles\n")
	for _, c := range res {
		fmt.Printf("%s\t%d -> %d\t", c.nombre, eloBase[c.nombre], c.elo)
		fmt.Printf("%d %d %d %d\t", c.p1, c.p2, c.p3, c.p4)
		fmt.Printf("%d %d %d %d %d\t", c.f16, c.oct, c.crt, c.sf, c.f)
		fmt.Printf("%d\t", c.camp)
		fmt.Printf("%.1f-%.1f\n", float64(c.mas)/float64(simulaciones), float64(c.menos)/float64(simulaciones))
	}
}
