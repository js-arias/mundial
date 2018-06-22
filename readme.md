# Simulador del mundial

Aunque llega un poco tarde,
este es mi simulador del mundial de futbol.

Utiliza el ELO de los equipos
justo antes de que el mundial arrancara
[wikipedia](https://en.wikipedia.org/wiki/World_Football_Elo_Ratings#Elo_Ratings_before_each_World_Championship).
Me gustaría una medida basada más
en la fuerza del equipo en el campo,
pero no conozco ningún indice así!

Al usar el ELO,
es posible calcular
[el resultado esperado del partido](https://en.wikipedia.org/wiki/World_Football_Elo_Ratings#Expected_result_of_match)
pero esa expectativa no incluye
la probabilidad de empates,
sin embargo,
creo que a lo largo de las simulaciones
el efecto de los empates debería disminuir.

Otra situación,
es que esta expectativa
tampoco permite estimar el número
de goles de diferencia.
Esto lo simulo de forma simple,
usando la parte entera
de una distribución exponencial
con lamda de 1
(promedio de 1),
y sumando 1
(el mínimo de diferencia en una victoria es un gol).

Otro inconveniente,
es que el calendario esta escrito
de forma directa sobre el código
(sería mucho más lindo que leyera
el calendario y los resultados previos
de un archivo independiente).
Yo lo ejecuto como si fuera un script
(usando *go run*)
y guardo los resultados en probs.tab.
Si quiero simular todo el mundial
desde el inicio,
la variable gobal allGames
se asigna como verdadera.
