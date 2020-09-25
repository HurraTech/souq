import axios from 'axios';
import { JAWHAR_API } from './constants';
const http = require('http');

export default class HurraServer {
  static getState() {
    console.log(`GETTING STATE from ${JAWHAR_API}`)
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .get(`${JAWHAR_API}/apps/${auid}/state`)
      .then(res => {
        console.log("RESULT OF STATE IS", res.data)
          resolve(res.data)
      })


    })
  }

  static setState(state) {
    console.log("SETTING STATE")
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .post(`${JAWHAR_API}/apps/${auid}/state`, state)
      .then(res => {
          resolve(res.data)
      })

    })
  }

  static patchState(state) {
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .patch(`${JAWHAR_API}/apps/${auid}/state`, state)
      .then(res => {
          resolve(res.data)
      })

    })
  }

  static exec(container, command, args = [], env = {}) {
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .post(`${JAWHAR_API}/apps/${auid}/${container}/command`, {
        Cmd: command,
        Args: args,
        Env: env
      })
      .then(res => {
          resolve(res.data)
      })

    })
  }


  static start_container(container) {
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .put(`${JAWHAR_API}/apps/${auid}/${container}`)
      .then(res => {
          resolve(res.data)
      })

    })
  }

  static stop_container(container) {
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .delete(`${JAWHAR_API}/apps/${auid}/${container}`)
      .then(res => {
          resolve(res.data)
      })

    })
  }


  static exec_sync(container, command, args = [], env = {}) {
    let auid = process.env.REACT_APP_AUID
    console.log(`EXECUTING COMMAND ${JAWHAR_API}/apps/${auid}/${container}/command`,{
        Cmd: command,
        Args: args,
        Env: env
    })

    return new Promise((resolve, reject) => {
      axios
      .post(`${JAWHAR_API}/apps/${auid}/${container}/command`, {
        Cmd: command,
        Args: args,
        Env: env
      })
      .then(res => {
        HurraServer.wait_for_cmd(res.data, resolve)
      })

    })
  }

  static wait_for_cmd(command, resolver)
  {
    HurraServer.get_command(command.ID).then(command_update => {
      if (command_update.Status == "completed") {
        console.log("Command completed", command_update)
        resolver(command_update)
      } else {
        setTimeout(() => { HurraServer.wait_for_cmd(command, resolver) }, 1000)
      }
    })
  }

  static get_command(cmd_id) {
    let auid = process.env.REACT_APP_AUID
    return new Promise((resolve, reject) => {
      axios
      .get(`${JAWHAR_API}/commands/${cmd_id}`)
      .then((statusRes) => {
        console.log(`Command ${cmd_id} Status`, statusRes.data.Status);
        resolve(statusRes.data)
      })
    })
  }


  static service_http_proxy = (proxy_host, proxy_port=80) => (oreq,ores) => {
	console.log("MAKING PROXY REQ TO", proxy_host, oreq.path, oreq.method)
     const options = {
       // host to forward to
       host: proxy_host,
       // port to forward to
       port: proxy_port,
       // path to forward to
       path: `${oreq.path}.html`,
       // request method
       method: oreq.method,
       // headers to send
       headers: oreq.headers,
     };

     const creq = http
       .request(options, pres => {
         // set encoding
         // pres.setEncoding('utf8');

         // set http status code based on proxied response
         ores.writeHead(pres.statusCode);

         // wait for data
         pres.on('data', chunk => {
           ores.write(chunk);
         });

         pres.on('close', () => {
           // closed, let's end client request as well
           ores.end();
         });

         pres.on('end', () => {
           // finished, let's finish client request as well
           ores.end();
         });
       })
       .on('error', e => {
         // we got an error
         console.log(e.message);
         try {
           // attempt to set error message and http status
           ores.writeHead(500);
           ores.write(e.message);
         } catch (e) {
           // ignore
         }
         ores.end();
       });

     creq.end();
  }

}

