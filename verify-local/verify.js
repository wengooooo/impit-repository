const fs = require('fs');
const path = require('path');
const koffi = require('koffi');
const imp = require('impitreq');

const p = imp.getLibPath();
if (!fs.existsSync(p)) {
  throw new Error('lib not found: ' + p);
}
const lib = koffi.load(p);
console.log('koffi loaded:', p);
const handle = lib.func('ImpitHandleRequestJSON', 'str', ['str']);
const result = handle('hello');
console.log('ImpitHandleRequestJSON:', result);

lib.func('ImpitHandleRequestJSON', 'str', ['str']);
