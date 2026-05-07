# HT

HT es un cliente HTTP en línea de comandos (CLI) que ejecuta peticiones definidas en archivos YAML. Está pensado para pruebas rápidas, debugging y ejecución de flujos HTTP reproducibles.

**Principales características**
- Ejecuta múltiples peticiones definidas en YAML.
- Soporta sustitución de variables (`${VAR}`) a partir de un bloque `vars` o un `env_file`.
- `dry-run` para imprimir la petición sin ejecutarla.
- `verbose` para traza detallada (DNS, TCP, TLS, primer byte).
- Formatea JSON en la salida con coloreado.
- Carga headers globales y por-request, admite body en línea o `body_file`.

**Instalación**
```bash
# desde el directorio del repo
go install ./cmd/ht

# o compilar localmente
go build -o ht ./cmd/ht
```

**Uso**
```bash
# Ejecutar una request por nombre
ht -yml http/dev.yml -name get-user

# Listar requests definidas en el YAML
ht -yml http/dev.yml -ls

# Mostrar la petición sin enviarla
ht -yml http/dev.yml -name create-post -dry-run

# Habilitar traza detallada
ht -yml http/dev.yml -name get-user -verbose
```

Flags:
- `-yml` : Ruta al archivo YAML de configuración (obligatorio).
- `-ls` : Lista las requests definidas.
- `-name` : Nombre de la request a ejecutar.
- `-dry-run` : Muestra la petición sin enviarla.
- `-verbose` : Muestra trazas detalladas (httptrace).

**Formato del YAML**

Ejemplo básico (saneado):

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
- `env_file` : Ruta a un `.env` con `key=value` por línea. Las variables se mezclan con `vars`.
- `requests[]` : Lista de requests con `name`, `method`, `url`, `headers`, `body` o `body_file`.

**Formato de `env_file`**
Líneas `KEY=value`. Comentarios que empiezan por `#` y líneas vacías son ignoradas.
