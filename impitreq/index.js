function resolvePlatformPackage() {
  const platform = process.platform;
  const arch = process.arch;

  if (platform === 'win32' && arch === 'x64') return 'impitreq-win32-x64';
  if (platform === 'linux' && arch === 'x64') return 'impitreq-linux-x64';
  if (platform === 'darwin' && arch === 'x64') return 'impitreq-darwin-x64';

  throw new Error(`Unsupported platform: ${platform}-${arch}`);
}

function getLibPath() {
  const platformPkg = resolvePlatformPackage();
  const mod = require(platformPkg);
  const libPath = mod && (mod.libPath || (typeof mod.getLibPath === 'function' ? mod.getLibPath() : undefined));
  if (!libPath) {
    throw new Error(`Platform package "${platformPkg}" did not export "libPath"`);
  }
  return libPath;
}

module.exports = {
  getLibPath,
  resolvePlatformPackage,
};

