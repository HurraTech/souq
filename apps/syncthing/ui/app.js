import React from 'react';
import axios from 'axios'
import  { Redirect } from 'react-router-dom'
class HurraApp extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
    };
  }

  render() {
	return ( window.location = `//${window.location.hostname}:28384`)
  }
}


export default HurraApp;
