# Performance e Profiling

## Objetivo

Registrar os experimentos realizados com o `pprof` na API Go, documentando CPU, memória e goroutines.

O objetivo desta etapa é aprender a coletar e interpretar perfis de performance, e não estabelecer um benchmark de capacidade ou limite de throughput da aplicação.

## Ambiente

- Aplicação: `desafio-go`
- Linguagem: Go
- API HTTP: porta `8080`
- Pprof: porta `6060`
- Banco: PostgreSQL
- Perfilador: `net/http/pprof`
- CPU profile: 10 segundos

## Metodologia

Foram realizados:

1. CPU profile durante chamadas autenticadas de `GET /products`.
2. Heap profile com `GET /debug/pprof/heap`.
3. Goroutine profile com `GET /debug/pprof/goroutine?debug=2`.

Também foi realizado um primeiro experimento gerando carga no próprio endpoint do Pprof, utilizado apenas para validar o funcionamento do profiler.

## Resumo

| Perfil | Resultado |
|---|---:|
| CPU | 80 ms de samples em 10 s |
| Heap `inuse_space` | 11.644,76 kB |
| Heap `alloc_space` | 12.701,44 kB |
| Goroutines | 6 no snapshot |

## Limitações

Os experimentos utilizaram carga pequena e controlada. Portanto, não representam capacidade máxima, throughput, latência sob carga ou comportamento em produção.

Os perfis de memória também incluem alocações relacionadas à inicialização e aos recursos do Swagger.

## Conclusão

O Pprof foi integrado em um servidor separado na porta `6060` e os principais perfis foram coletados com sucesso. Os resultados estabelecem uma baseline inicial para futuras comparações após alterações arquiteturais ou de performance.
EOF