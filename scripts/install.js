const os = require('os');
const fs = require('fs');
const https = require('https');
const path = require('path');


const version = "0.1.2";
const repo = "franciscorojas27/ht";

const platform = os.platform() === 'win32' ? 'windows' : os.platform();
const arch = os.arch() === 'x64' ? 'amd64' : (os.arch() === 'arm64' ? 'arm64' : '386');
const extension = platform === 'windows' ? 'zip' : 'tar.gz';

const filename = `ht_${version}_${platform}_${arch}.${extension}`;
const url = `https://github.com/${repo}/releases/download/v${version}/${filename}`;

const binDir = path.join(__dirname, '../bin');
if (!fs.existsSync(binDir)) fs.mkdirSync(binDir);

const dest = path.join(binDir, platform === 'windows' ? 'ht.exe' : 'ht');

console.log(`Descargando ht v${version} desde: ${url}`);

const download = (url, dest) => {
    https.get(url, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
            download(res.headers.location, dest);
            return;
        }
        if (res.statusCode !== 200) {
            console.error(`Error: Servidor respondió con ${res.statusCode}`);
            process.exit(1);
        }

        const file = fs.createWriteStream(dest);
        res.pipe(file);
        file.on('finish', () => {
            file.close();
            if (platform !== 'windows') fs.chmodSync(dest, 0o755);
            console.log('Instalación completada.');
        });
    }).on('error', (err) => {
        console.error(`Error de red: ${err.message}`);
        process.exit(1);
    });
};

download(url, dest);