# Goroutine Profile

## 1. Coleta

O perfil de goroutines foi coletado utilizando o endpoint `/debug/pprof/goroutine` com `debug=2`, permitindo obter o stack trace das goroutines:

```bash
curl -o goroutine.txt "http://localhost:6060/debug/pprof/goroutine?debug=2"

Para verificar a quantidade de goroutines capturadas no snapshot:

grep -c '^goroutine ' goroutine.txt
2. Resultado

O snapshot capturou:

6 goroutines

Foram identificadas principalmente:

servidor HTTP principal;
database/sql connection opener;
health check do pgxpool;
servidor Pprof;
requisição ao endpoint de profiling;
runtime/atividades HTTP.
3. Servidor HTTP Principal

A goroutine principal estava em estado IO wait, com stack apontando para:

main.main()
cmd/api/main.go

Isso corresponde ao servidor HTTP principal da aplicação.

4. PostgreSQL

Foi identificada uma goroutine responsável pela abertura de conexões do database/sql:

database/sql.(*DB).connectionOpener

Também foi identificada a goroutine responsável pelo health check do pool PostgreSQL:

github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck
5. Pprof

Foi identificada uma goroutine executando:

desafio-go/internal/app.StartPprofServer.func1()
internal/app/pprof.go

Isso confirma que o servidor Pprof está sendo executado separadamente do servidor HTTP principal.

6. Interpretação

O número de 6 goroutines representa uma fotografia do momento em que o snapshot foi coletado. Portanto, não representa um número fixo nem um limite de goroutines da aplicação.

No snapshot analisado, não foi observado um volume anormal de goroutines.

7. Conclusão

O profile confirmou a separação do servidor Pprof e permitiu identificar as principais goroutines relacionadas ao servidor HTTP, ao PostgreSQL e ao próprio mecanismo de profiling.