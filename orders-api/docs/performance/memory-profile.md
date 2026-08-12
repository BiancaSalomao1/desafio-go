# Memory Profile

## 1. Coleta

O perfil de memória foi coletado utilizando o endpoint `/debug/pprof/heap`:

```bash
curl -o heap.prof "http://localhost:6060/debug/pprof/heap"

Análise do perfil:

go tool pprof -text heap.prof

Para analisar o espaço total alocado (alloc_space):

go tool pprof -sample_index=alloc_space -text heap.prof
2. In-use Space

O perfil inuse_space representa a memória que permanecia em uso no momento da coleta.

Total observado: 11.644,76 kB

Principais entradas
Função	Memória	Percentual
golang.org/x/net/webdav.(*memFile).Write	8.569,41 kB	73,59%
runtime.mallocgc	2.563,30 kB	22,01%
golang.org/x/net/webdav.(*memFS).OpenFile	512,04 kB	4,40%

Também foram observadas funções de inicialização de:

github.com/swaggo/files
3. Alloc Space

O perfil alloc_space representa o espaço total alocado contabilizado pelo profiler.

Total observado: 12.701,44 kB

Principais entradas
Função	Memória	Percentual
golang.org/x/net/webdav.(*memFile).Write	8.569,41 kB	67,47%
runtime.mallocgc	2.563,30 kB	20,18%
golang.org/x/net/webdav.(*memFS).OpenFile	1.024,05 kB	8,06%
net.open	544,67 kB	4,29%
4. Interpretação

Os resultados mostram uma predominância de funções relacionadas a:

golang.org/x/net/webdav;
github.com/swaggo/files;
mecanismos de alocação do runtime do Go.

A presença de webdav e github.com/swaggo/files indica influência da inicialização e dos recursos estáticos utilizados pela documentação Swagger.

Importante: não é correto concluir que os 73,59% representam consumo de memória da lógica de negócio da API. Esse percentual representa a distribuição das amostras de memória observadas pelo profiler durante o snapshot.

5. Limitações

O snapshot inclui alocações relacionadas à inicialização da aplicação.

Para investigar possíveis vazamentos de memória ou crescimento contínuo do consumo, seria necessário realizar múltiplos snapshots, após períodos de execução e sob cargas controladas, comparando os resultados ao longo do tempo.

6. Conclusão

O profile estabelece uma baseline inicial de memória da aplicação e permite identificar as principais fontes de alocação observadas durante o experimento.

Neste cenário, as maiores ocorrências estão relacionadas principalmente aos recursos de webdav, Swagger e ao runtime do Go, não havendo evidência suficiente neste experimento para atribuir o consumo à lógica de negócio da aplicação.