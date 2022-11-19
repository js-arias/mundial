// Copyright (c) 2018 J. Salvador Arias <jsalarias@gmail.com>.
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type team struct {
	name   string
	points int
	goals  int
	elo    int
}

var teams = map[string]*team{
	"bra": &team{"bra", 0, 0, 2142},
	"deu": &team{"deu", 0, 0, 2077},
	"esp": &team{"esp", 0, 0, 2044},
	"fra": &team{"fra", 0, 0, 1987},
	"arg": &team{"arg", 0, 0, 1986},
	"prt": &team{"prt", 0, 0, 1970},
	"eng": &team{"eng", 0, 0, 1948},
	"bel": &team{"bel", 0, 0, 1939},
	"col": &team{"col", 0, 0, 1928},
	"per": &team{"per", 0, 0, 1915},
	"uru": &team{"uru", 0, 0, 1894},
	"che": &team{"che", 0, 0, 1890},
	"dnk": &team{"dnk", 0, 0, 1856},
	"hrv": &team{"hrv", 0, 0, 1853},
	"mex": &team{"mex", 0, 0, 1850},
	"pol": &team{"pol", 0, 0, 1831},
	"swe": &team{"swe", 0, 0, 1795},
	"irn": &team{"irn", 0, 0, 1789},
	"rus": &team{"rus", 0, 0, 1778}, // +100 host
	"srb": &team{"srb", 0, 0, 1777},
	"isl": &team{"isl", 0, 0, 1764},
	"sen": &team{"sen", 0, 0, 1750},
	"cri": &team{"cri", 0, 0, 1744},
	"aus": &team{"aus", 0, 0, 1742},
	"mar": &team{"mar", 0, 0, 1733},
	"kor": &team{"kor", 0, 0, 1714},
	"jpn": &team{"jpn", 0, 0, 1684},
	"nga": &team{"nga", 0, 0, 1681},
	"pan": &team{"pan", 0, 0, 1659},
	"tun": &team{"tun", 0, 0, 1657},
	"egy": &team{"egy", 0, 0, 1646},
	"sau": &team{"sau", 0, 0, 1591},
}

type byPoints []*team

func (a byPoints) Len() int { return len(a) }

func (a byPoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a byPoints) Less(i, j int) bool {
	if a[i].points == a[j].points {
		return a[i].goals < a[j].goals
	}
	return a[i].points < a[j].points
}

func match(t1, t2 *team) {
	if t2.elo > t1.elo {
		t1, t2 = t2, t1
	}
	dif := float64(t1.elo-t2.elo) / 400
	e := 1 / (math.Pow(10, -dif) + 1)
	g := int(rand.ExpFloat64()) + 1
	if rand.Float64() < e {
		t1.points += 3
		t1.goals += g
		t2.goals -= g
	} else {
		t2.points += 3
		t2.goals += g
		t1.goals -= g
	}
}

var allGames = false

func groupA() (*team, *team) {
	if !allGames {
		teams["rus"].points = 6
		teams["rus"].goals = 4
		teams["uru"].points = 9
		teams["uru"].goals = 5
		teams["egy"].points = 0
		teams["egy"].goals = -4
		teams["sau"].points = 3
		teams["sau"].goals = -5
	} else {
		match(teams["rus"], teams["sau"])
		match(teams["uru"], teams["egy"])
		match(teams["rus"], teams["egy"])
		match(teams["uru"], teams["sau"])
		match(teams["uru"], teams["rus"])
		match(teams["sau"], teams["egy"])
	}

	gv := []*team{
		teams["rus"],
		teams["uru"],
		teams["sau"],
		teams["egy"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupB() (*team, *team) {
	if !allGames {
		teams["esp"].points = 5
		teams["esp"].goals = 2 // to make spain the leader
		teams["prt"].points = 5
		teams["prt"].goals = 1
		teams["irn"].points = 4
		teams["irn"].goals = 0
		teams["mar"].points = 1
		teams["mar"].goals = -2
	} else {
		match(teams["esp"], teams["prt"])
		match(teams["irn"], teams["mar"])
		match(teams["prt"], teams["mar"])
		match(teams["esp"], teams["irn"])
		match(teams["irn"], teams["prt"])
		match(teams["esp"], teams["mar"])
	}

	gv := []*team{
		teams["esp"],
		teams["prt"],
		teams["irn"],
		teams["mar"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupC() (*team, *team) {
	if !allGames {
		teams["fra"].points = 7
		teams["fra"].goals = 2
		teams["dnk"].points = 5
		teams["dnk"].goals = 1
		teams["aus"].points = 1
		teams["aus"].goals = -3
		teams["per"].points = 3
		teams["per"].goals = 0
	} else {
		match(teams["fra"], teams["aus"])
		match(teams["dnk"], teams["per"])
		match(teams["dnk"], teams["aus"])
		match(teams["fra"], teams["per"])
		match(teams["dnk"], teams["fra"])
		match(teams["aus"], teams["per"])
	}

	gv := []*team{
		teams["fra"],
		teams["dnk"],
		teams["per"],
		teams["aus"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupD() (*team, *team) {
	if !allGames {
		teams["hrv"].points = 9
		teams["hrv"].goals = 6
		teams["isl"].points = 1
		teams["isl"].goals = -3
		teams["arg"].points = 4
		teams["arg"].goals = -2
		teams["nga"].points = 3
		teams["nga"].goals = -1
	} else {
		match(teams["arg"], teams["isl"])
		match(teams["hrv"], teams["nga"])
		match(teams["arg"], teams["hrv"])
		match(teams["nga"], teams["isl"])
		match(teams["nga"], teams["arg"])
		match(teams["isl"], teams["hrv"])
	}

	gv := []*team{
		teams["arg"],
		teams["hrv"],
		teams["isl"],
		teams["nga"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupE() (*team, *team) {
	if !allGames {
		teams["srb"].points = 3
		teams["srb"].goals = -2
		teams["bra"].points = 7
		teams["bra"].goals = 4
		teams["che"].points = 5
		teams["che"].goals = 1
		teams["cri"].points = 1
		teams["cri"].goals = -3
	} else {
		match(teams["cri"], teams["srb"])
		match(teams["bra"], teams["che"])
		match(teams["bra"], teams["cri"])
		match(teams["srb"], teams["che"])
		match(teams["srb"], teams["bra"])
		match(teams["che"], teams["cri"])
	}

	gv := []*team{
		teams["bra"],
		teams["che"],
		teams["srb"],
		teams["cri"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupF() (*team, *team) {
	if !allGames {
		teams["swe"].points = 6
		teams["swe"].goals = 3
		teams["mex"].points = 6
		teams["mex"].goals = -1
		teams["deu"].points = 3
		teams["deu"].goals = -2
		teams["kor"].points = 3
		teams["kor"].goals = 0
	} else {
		match(teams["deu"], teams["mex"])
		match(teams["swe"], teams["kor"])
		match(teams["kor"], teams["mex"])
		match(teams["deu"], teams["swe"])
		match(teams["kor"], teams["deu"])
		match(teams["mex"], teams["swe"])
	}
	gv := []*team{
		teams["deu"],
		teams["swe"],
		teams["mex"],
		teams["kor"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupG() (*team, *team) {
	if !allGames {
		teams["bel"].points = 9
		teams["bel"].goals = 7
		teams["eng"].points = 6
		teams["eng"].goals = 5
		teams["tun"].points = 3
		teams["tun"].goals = -3
		teams["pan"].points = 0
		teams["pan"].goals = -9
	} else {
		match(teams["bel"], teams["pan"])
		match(teams["eng"], teams["tun"])
		match(teams["bel"], teams["tun"])
		match(teams["eng"], teams["pan"])
		match(teams["eng"], teams["bel"])
		match(teams["pan"], teams["tun"])
	}

	gv := []*team{
		teams["eng"],
		teams["bel"],
		teams["tun"],
		teams["pan"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func groupH() (*team, *team) {
	if !allGames {
		teams["jpn"].points = 4
		teams["jpn"].goals = 1 // to classify japan
		teams["sen"].points = 4
		teams["sen"].goals = 0
		teams["pol"].points = 3
		teams["pol"].goals = -3
		teams["col"].points = 6
		teams["col"].goals = 3
	} else {
		match(teams["col"], teams["jpn"])
		match(teams["pol"], teams["sen"])
		match(teams["jpn"], teams["sen"])
		match(teams["pol"], teams["col"])
		match(teams["jpn"], teams["pol"])
		match(teams["sen"], teams["col"])
	}

	gv := []*team{
		teams["col"],
		teams["pol"],
		teams["sen"],
		teams["jpn"],
	}
	sort.Sort(byPoints(gv))
	return gv[3], gv[2]
}

func eliminator(t1, t2 *team) *team {
	t1.points = 0
	t2.points = 0
	match(t1, t2)
	if t1.points > t2.points {
		return t1
	}
	return t2
}

type counter struct {
	name  string
	group int
	r16   int
	qf    int
	sf    int
	champ int
}

type byChamp []counter

func (a byChamp) Len() int { return len(a) }

func (a byChamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a byChamp) Less(i, j int) bool {
	if a[i].champ == a[j].champ {
		if a[i].group == a[j].group {
			return teams[a[j].name].elo < teams[a[i].name].elo
		}
		return a[j].group < a[i].group
	}
	return a[j].champ < a[i].champ
}

func main() {
	champ := make(map[string]*counter)
	for _, t := range teams {
		champ[t.name] = &counter{name: t.name}
	}
	const reps = 100000
	verbose := false
	for i := 0; i < reps; i++ {
		a1, a2 := groupA()
		champ[a1.name].group++
		champ[a2.name].group++
		b1, b2 := groupB()
		champ[b1.name].group++
		champ[b2.name].group++
		c1, c2 := groupC()
		champ[c1.name].group++
		champ[c2.name].group++
		d1, d2 := groupD()
		champ[d1.name].group++
		champ[d2.name].group++
		e1, e2 := groupE()
		champ[e1.name].group++
		champ[e2.name].group++
		f1, f2 := groupF()
		champ[f1.name].group++
		champ[f2.name].group++
		g1, g2 := groupG()
		champ[g1.name].group++
		champ[g2.name].group++
		h1, h2 := groupH()
		champ[h1.name].group++
		champ[h2.name].group++

		if verbose {
			fmt.Printf("a1: %s [%d], a2: %s [%d]\n", a1.name, a1.points, a2.name, a2.points)
			fmt.Printf("b1: %s [%d], b2: %s [%d]\n", b1.name, b1.points, b2.name, b2.points)
			fmt.Printf("c1: %s [%d], c2: %s [%d]\n", c1.name, c1.points, c2.name, c2.points)
			fmt.Printf("d1: %s [%d], d2: %s [%d]\n", d1.name, d1.points, d2.name, d2.points)
			fmt.Printf("e1: %s [%d], e2: %s [%d]\n", e1.name, e1.points, e2.name, e2.points)
			fmt.Printf("f1: %s [%d], f2: %s [%d]\n", f1.name, f1.points, f2.name, f2.points)
			fmt.Printf("g1: %s [%d], g2: %s [%d]\n", g1.name, g1.points, g2.name, g2.points)
			fmt.Printf("h1: %s [%d], h2: %s [%d]\n", h1.name, h1.points, h2.name, h2.points)
			fmt.Printf("\n")
		}

		// round of 16
		var m49, m50, m51, m52 *team
		var m53, m54, m55, m56 *team
		if !allGames {
			m49 = teams["uru"]
			m50 = teams["fra"]
			m51 = teams["rus"]
			m52 = teams["hrv"]
			m53 = teams["bra"]
			m54 = teams["bel"]
			m55 = teams["swe"]
			m56 = teams["eng"]
		} else {
			m49 = eliminator(a1, b2)
			m50 = eliminator(c1, d2)
			m51 = eliminator(b1, a2)
			m52 = eliminator(d1, c2)
			m53 = eliminator(e1, f2)
			m54 = eliminator(g1, h2)
			m55 = eliminator(f1, e2)
			m56 = eliminator(h1, g2)
		}
		champ[m49.name].r16++
		champ[m50.name].r16++
		champ[m51.name].r16++
		champ[m52.name].r16++
		champ[m53.name].r16++
		champ[m54.name].r16++
		champ[m55.name].r16++
		champ[m56.name].r16++

		// qfinals
		m57 := eliminator(m49, m50)
		champ[m57.name].qf++
		m58 := eliminator(m53, m54)
		champ[m58.name].qf++
		m59 := eliminator(m51, m52)
		champ[m59.name].qf++
		m60 := eliminator(m55, m56)
		champ[m60.name].qf++

		// semis
		m61 := eliminator(m57, m58)
		champ[m61.name].sf++
		m62 := eliminator(m59, m60)
		champ[m62.name].sf++

		// final
		ch := eliminator(m61, m62)
		if verbose {
			fmt.Printf("8v: %s - %s -> %s\n", a1.name, b2.name, m49.name)
			fmt.Printf("8v: %s - %s -> %s\n", c1.name, d2.name, m50.name)
			fmt.Printf("8v: %s - %s -> %s\n", e1.name, f2.name, m53.name)
			fmt.Printf("8v: %s - %s -> %s\n", g1.name, h2.name, m54.name)
			fmt.Printf("8v: %s - %s -> %s\n", b1.name, a2.name, m51.name)
			fmt.Printf("8v: %s - %s -> %s\n", d1.name, c2.name, m52.name)
			fmt.Printf("8v: %s - %s -> %s\n", f1.name, e2.name, m55.name)
			fmt.Printf("8v: %s - %s -> %s\n", h1.name, g2.name, m56.name)
			fmt.Printf("\n")
			fmt.Printf("4t: %s - %s -> %s\n", m49.name, m50.name, m57.name)
			fmt.Printf("4t: %s - %s -> %s\n", m53.name, m54.name, m58.name)
			fmt.Printf("4t: %s - %s -> %s\n", m51.name, m52.name, m59.name)
			fmt.Printf("4t: %s - %s -> %s\n", m55.name, m56.name, m60.name)
			fmt.Printf("\n")
			fmt.Printf("sm: %s - %s -> %s\n", m57.name, m58.name, m61.name)
			fmt.Printf("sm: %s - %s -> %s\n", m59.name, m60.name, m62.name)
			fmt.Printf("\n")
			fmt.Printf("fn: %s - %s -> %s\n", m61.name, m62.name, ch.name)
			fmt.Printf("\n")
		}
		champ[ch.name].champ++

		// reset games
		for _, t := range teams {
			t.points = 0
			t.goals = 0
		}
	}
	cl := make([]counter, 0, len(champ))
	for _, c := range champ {
		cl = append(cl, *c)
	}
	sort.Sort(byChamp(cl))

	fmt.Printf("rk\tteam\tGrp\tr16\tqf\tsm\tchamp\n")
	for i, c := range cl {
		fmt.Printf("%d\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.6f\n", i+1, c.name, float64(c.group)/reps, float64(c.r16)/reps, float64(c.qf)/reps, float64(c.sf)/reps, float64(c.champ)/reps)
	}
}
