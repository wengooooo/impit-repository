const path = require('path');

function resolvePrebuildPath() {
  const platform = process.platform;
  const arch = process.arch;

  if (platform === 'win32' && arch === 'x64') return path.join(__dirname, 'prebuilds', 'win32-x64', 'impitreq.dll');
  if (platform === 'linux' && arch === 'x64') return path.join(__dirname, 'prebuilds', 'linux-x64', 'impitreq.so');
  throw new Error(`Unsupported platform: ${platform}-${arch}`);
}

function getLibPath() {
  const libPath = resolvePrebuildPath();
  if (!libPath) {
    throw new Error('Prebuilt library path not resolved');
  }
  return libPath;
}

module.exports = {
  getLibPath,
  resolvePrebuildPath,
};
