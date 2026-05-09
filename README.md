# HT

HT es un cliente HTTP en linea de comandos (CLI) pensado para pruebas rapidas, debugging y ejecucion de flujos HTTP reproducibles. Funciona con definiciones en YAML y tambien con un modo rapido tipo cURL.

**Caracteristicas principales**
- Ejecuta multiples requests definidas en YAML.
- Sustitucion de variables (`${VAR}`) con bloque `vars` o `env_file`.
- Modo rapido con `-url` y opciones tipo cURL.
- `dry-run` para imprimir la peticion sin enviarla.
- `verbose` con traza detallada (DNS, TCP, TLS, primer byte).
- Formatea JSON con coloreado.
- Headers globales y por-request, body en linea o `body_file`.

**Instalacion**
```bash
# desde releases (recomendado)
# descarga el binario para tu OS y ponlo en tu PATH

# desde el repo
go install ./cmd/ht

# o compilar localmente
go build -o ht ./cmd/ht
```

**Uso rapido (modo cURL)**
El modo rapido se activa si no pasas `-yml` o si usas `-url`/`-X`.

```bash
# GET simple
ht -url https://example.com

# POST con body
ht -url https://example.com -d "a=1&b=2"

# Headers multiples
ht -url https://example.com -H "Accept: application/json" -H "X-Token: 123"

# HEAD
ht -url https://example.com -I

# Basic auth
ht -url https://example.com -u user:pass

# TLS inseguro
ht -url https://example.com -k
```

**Uso con YAML**
```bash
# Ejecutar una request por nombre
ht -yml http/dev.yml -name get-user

# Listar requests definidas en el YAML
ht -yml http/dev.yml -ls

# Mostrar la peticion sin enviarla
ht -yml http/dev.yml -name create-post -dry-run

# Habilitar traza detallada
ht -yml http/dev.yml -name get-user -verbose
```

**Flags principales**
- `-yml` : Ruta al archivo YAML de configuracion.
- `-ls` : Lista las requests definidas.
- `-name` : Nombre de la request a ejecutar.
- `-dry-run` : Muestra la peticion sin enviarla.
- `-verbose` : Muestra trazas detalladas (httptrace).
- `-url` : URL directa para modo rapido.
- `-X` : Metodo HTTP para modo rapido.
- `-H` : Header extra (repetible).
- `-d` : Body para modo rapido.
- `-I` : Solo headers (HEAD).
- `-k` : TLS inseguro.
- `-u` : Basic auth `user:pass`.

**Formato del YAML**
Ejemplo basico (saneado):

```yaml
config:
  base_url: https://api.example.com
  timeout: 30
  headers:
    Content-Type: application/json; charset=UTF-8

vars:
  token: "MY_TOKEN"

env_file: ./http/.env

requests:
  - name: get-user
    method: GET
    url: /
    headers:
      User-Agent: "MyClient/1.0"

  - name: create-post
    method: POST
    url: /posts
    body:
      title: "hola"
      body: "ejemplo"
      userId: 1
```

Campos principales:
- `config.base_url` : URL base opcional que se concatena con `url` de cada request.
- `config.timeout` : Timeout en segundos (por request).
- `config.headers` : Headers aplicados a todas las requests.
- `vars` : Mapa de variables que pueden referenciarse como `${var}` en el YAML.
- `env_file` : Ruta a un `.env` con `key=value` por linea. Las variables se mezclan con `vars`.
- `requests[]` : Lista de requests con `name`, `method`, `url`, `headers`, `body` o `body_file`.

**Formato de `env_file`**
Lineas `KEY=value`. Comentarios que empiezan por `#` y lineas vacias son ignoradas.
