import express from 'express';
import HurraApp from './HurraApp/app';

var args = process.argv.slice(2);

const server = express();
const port = process.env.PORT || 5000;

server.use(express.json())

let app = new HurraApp(server)
app.start()

if (args[0] == "init") {
  app.init();
} else {
  server.listen(port, () => console.log(`Listening on port ${port}`));
}
