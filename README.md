# Simulador del mundial

## Funcionamiento del modelo

Para estimar los goles en un partido,
utilicé 2.6 como promedio de gol.
En base a eso,
se asumí que en todos los partidos,
la suma de la expectativa de goles
entre ambos equipos sumaría 2.6.
El número de goles se ajusto entonces
al valor esperado de la probabilidad de victoria
(calculada con el ELO)
comparado con la probabilidad de victoria
dado una combinación de la expectativa de goles
y usando una [distribución de Poisson](https://en.wikipedia.org/wiki/Poisson_distribution).
Esto se calculo con el programa de la carpeta ´probs´.

Esta es la tabla:

Equipo fuerte | Equipo débil | Expectativa de victoria
------------- | ------------ | -----------------------
1.3 | 1.3 | 0.500
1.4 | 1.2 | 0.547
1.5 | 1.1 | 0.594
1.6 | 1.0 | 0.638
1.7 | 0.9 | 0.682
1.8 | 0.8 | 0.725
1.9 | 0.7 | 0.764
2.0 | 0.6 | 0.802
2.1 | 0.5 | 0.837
2.2 | 0.4 | 0.869
2.3 | 0.3 | 0.897
2.4 | 0.2 | 0.922
2.5 | 0.1 | 0.944
