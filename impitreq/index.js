const path = require('path');
const fs = require('fs');

function getLibPath() {
  const platform = process.platform;
  const arch = process.arch;
  
  let folder = '';
  let file = '';

  if (platform === 'win32' && arch === 'x64') {
    folder = 'win32-x64';
    file = 'impitreq.dll';
  } else if (platform === 'linux' && arch === 'x64') {
    folder = 'linux-x64';
    file = 'impitreq.so';
  } else if (platform === 'darwin' && arch === 'x64') {
    folder = 'darwin-x64';
    file = 'impitreq.dylib';
  } else {
    throw new Error(`Unsupported platform: ${platform}-${arch}`);
  }

  const libPath = path.join(__dirname, 'prebuilds', folder, file);
  if (!fs.existsSync(libPath)) {
    throw new Error(`Library not found at: ${libPath}.`);
  }
  return libPath;
}

module.exports = {
  getLibPath
};
