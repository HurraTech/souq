import React from 'react';
import axios from 'axios'
import  { Redirect } from 'react-router-dom'
class HurraApp extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      content: ""
    };
  }

  async componentDidMount() {
    let page = (await (await fetch('/index.html')).text());
    this.setState({
      content: page
    })
  }

  render() {
	return <div dangerouslySetInnerHTML={{__html: this.state.content}} />
  }
}


export default HurraApp;
