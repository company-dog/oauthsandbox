import React, { useState, useEffect } from 'react';

import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom';

import './App.css';

const App = () => {
  const [accessToken, setAccessToken] = useState('');
  const [expiresIn, setExpiresIn] = useState('');
  const [sessionState, setSessionState] = useState('');
  const [tokenType, setTokenType] = useState('');

  const setStateValue = hashMap => {
    setAccessToken(hashMap.access_token);
    setExpiresIn(hashMap.expires_in);
    setSessionState(hashMap.session_state);
    setTokenType(hashMap.token_type);
  };

  const onCheckStateClick = () => {
    console.log({
      accessToken,
      expiresIn,
      sessionState,
      tokenType
    });
  };

  return (
    <Router>
      <div>
        <div className="App">
          <h1>Implicit Grant type</h1>
        </div>
        <button onClick={onCheckStateClick}>Check state.</button>
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
            <li>
              <Link to="/service">Service</Link>
            </li>
            <li>
              <Link to="/logout">Logout</Link>
            </li>
          </ul>
        </nav>

        {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
        <Switch>
          <Route path="/login">
            <Login />
          </Route>
          <Route path="/callback">
            <Callback
              setStateValue={setStateValue}
              accessToken={accessToken}
              expiresIn={expiresIn}
              sessionState={sessionState}
              tokenType={tokenType}
            />
          </Route>
          <Route path="/service">
            <Service accessToken={accessToken} />
          </Route>
          <Route path="/">
            <Home />
          </Route>
        </Switch>
      </div>
    </Router>
  );
};

function Home() {
  return <h2>Home</h2>;
}

function Login() {
  window.location =
    'http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/auth?client_id=implicitClient&response_type=token&redirect_uri=http://localhost:3000/callback&scope=getBillingService';
  return null;
}

function Service(props) {
  const [services, setServices] = useState([]);
  const { accessToken } = props;

  useEffect(() => {
    // access protected resources
    // Post + form
    const formData = new FormData();
    formData.append('access_token', accessToken);
    fetch('http://localhost:8081/billing/v1/services', {
      method: 'POST',
      body: formData
    })
      .then(response => response.json())
      .then(data => {
        console.log(data);
        setServices(data);
      });
  }, [accessToken]);

  return (
    <div>
      <h2>Services</h2>
      <div>{JSON.stringify({ services }, null, 2)}</div>
    </div>
  );

  // // parse response
  // const services = [];
  // services.push(<div key="a">billing A</div>);
  // services.push(<div key="b">billing B</div>);
  // services.push(<div key="c">billing C</div>);
  // return services;
}

function Callback(props) {
  // Get access token
  const hashStr = window.location.hash;
  const hashMap = hashStr
    .substr(1)
    .split('&')
    .reduce((acc, item) => {
      // add item to accumulator
      const kv = item.split('=');
      acc[kv[0]] = kv[1];
      return acc;
    }, {});
  // console.log(hashMap);

  // setState...
  const {
    setStateValue,
    accessToken,
    tokenType,
    expiresIn,
    sessionState
  } = props;

  setStateValue(hashMap);

  return <h2>Callback</h2>;
}

export default App;
