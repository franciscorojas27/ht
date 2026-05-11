const fs = require('fs');
const path = require('path');
const https = require('https');
const os = require('os');

const REPO = "franciscorojas27/HT";
const VERSION = require('../package.json').version;

const platformMap = { win32: "windows", darwin: "darwin", linux: "linux" };
const archMap = { x64: "amd64", arm64: "arm64" };

const platform = platformMap[os.platform()];
const arch = archMap[os.arch()];

if (!platform || !arch) {
  console.error(`Plataforma no soportada: ${os.platform()} ${os.arch()}`);
  process.exit(1);
}

const exe = platform === "windows" ? "ht.exe" : "ht";
const url = `https://github.com/${REPO}/releases/download/v${VERSION}/ht-${platform}-${arch}.zip`;
const dest = path.join(__dirname, '../bin', exe);

function download(url, dest) {
  https.get(url, (res) => {
    if (res.statusCode === 302 || res.statusCode === 301) {
      return download(res.headers.location, dest);
    }
    if (res.statusCode !== 200) {
      console.error(`Error al descargar: ${res.statusCode}`);
      process.exit(1);
    }
    const file = fs.createWriteStream(dest);
    res.pipe(file);
    file.on('finish', () => {
      file.close();
      if (process.platform !== "win32") fs.chmodSync(dest, 0o755);
      console.log("ht instalado correctamente.");
    });
  });
}

console.log(`Descargando ht v${VERSION} para ${platform}-${arch}...`);
download(url, dest);