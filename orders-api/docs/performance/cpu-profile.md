# CPU Profile

## 1. Procedimento

Foi utilizado o endpoint de profiling de CPU do Pprof durante um período de 10 segundos:

```bash
curl -o cpu-api.prof "http://localhost:6060/debug/pprof/profile?seconds=10"

Durante a coleta foram realizadas 200 chamadas autenticadas:

GET /products

A análise do perfil foi realizada com:

go tool pprof -text cpu-api.prof
2. Resultado

O profile apresentou:

Duration: 10s
Total samples = 80ms (0.8%)
Principais amostras
Função	CPU	Percentual
internal/runtime/syscall/linux.Syscall6	40 ms	50,00%
log/slog.Value.Kind	10 ms	12,50%
runtime.(*waitq).dequeue	10 ms	12,50%
runtime.nanotime	10 ms	12,50%
time.appendInt	10 ms	12,50%

Também apareceram componentes relacionados a:

net/http;
syscall;
slog;
runtime;
middleware Logger;
middleware Recovery.
3. Interpretação

A utilização de CPU foi baixa durante o experimento, com 80 ms de samples em uma janela de 10 segundos.

O slog aparece no caminho de execução, incluindo:

log/slog.Info
log/slog.(*JSONHandler).Handle

Essa presença é coerente com a instrumentação de logs estruturados adicionada à aplicação.

Os percentuais apresentados pelo Pprof representam a distribuição das amostras coletadas pelo profiler, e não a porcentagem da CPU total da máquina.

Com a carga utilizada no experimento, não foi identificado um hotspot relevante relacionado à lógica de negócio da aplicação.

4. Limitação

As 200 requisições utilizadas durante a coleta não constituem um benchmark de capacidade.

Para avaliar capacidade de forma adequada, seriam necessários testes controlados considerando, por exemplo:

throughput;
latência;
utilização de CPU;
utilização de memória;
comportamento sob diferentes níveis de carga.
5. Conclusão

O CPU profile confirmou o funcionamento do Pprof e permitiu observar o caminho de execução da API durante as requisições.

Com os dados coletados neste experimento, não há evidência suficiente para considerar o slog um gargalo de CPU da aplicação.