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

lib.func('ImpitHandleRequestJSON', 'str', ['str']);
