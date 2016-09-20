// Make Webpack include these in the dist folder.
require('../index.scss');

const Elm = require('../elm/Main');

var mountNode = document.getElementById('main');

// The third value on embed are the initial values for incomming ports into Elm
var app = Elm.Main.embed(mountNode);
