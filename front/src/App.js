import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom';

import Home from './components/Home';
import Login from './components/Login';
import Admin from './components/Admin';

import Movies from './components/Movies';
import OneMovie from './components/OneMovie';
import EditMovie from './components/EditMovie';

import Genres from './components/Genres'
import OneGenre from './components/OneGenre';

export default function App() {

  const [token, setToken] = useState('')
  const [loaded, setLoaded] = useState(false)

  function handleJwtChange(jwt) {
    setToken(jwt)
  }

  function logout() {
    setToken('')
    window.localStorage.removeItem('token')
  }

  useEffect(() => {
    const t = window.localStorage.getItem('token')
    if (t) {
      if (token === '') {
        setToken(t)
      }
    }

    setLoaded(true)
  }, [token])

  if (!loaded) return (<p>Loading....</p>)

  return (
    <Router>

      <div className="container">
        <div className="row">
          <div className='col mt-3'>
            <h1 className="mt-3">
              Watch a Movie!
            </h1>
          </div>
          <div className='col mt-3 text-end'>
            {token === '' && (<Link to="/login">Login</Link>)}
            {token !== '' && (<Link to="/logout" onClick={() => logout()}>Logout</Link>)}
          </div>
          <hr className="mb-3"></hr>
        </div>

        <div className="row">
          <div className="col-md-2">
            <nav>
              <ul className="list-group">
                <li className="list-group-item">
                  <Link to="/">Home</Link>
                </li>
                <li className="list-group-item">
                  <Link to="/movies">Movies</Link>
                </li>
                <li className="list-group-item">
                  <Link to="/genres">Genres</Link>
                </li>
                {token !== '' &&
                  (
                    <>
                      <li className="list-group-item">
                        <Link to="/admin/movie/0">Add Movie</Link>
                      </li>
                      <li className="list-group-item">
                        <Link to="/admin">Manage Catalogue</Link>
                      </li>
                    </>
                  )}

              </ul>
            </nav>
          </div>

          <div className="col-md-10">
            <Switch>

              <Route exact path="/login" component={((props) =>
                <Login {...props} handleJwtChange={handleJwtChange} />)}
              />


              <Route path="/movies/:id" component={OneMovie} />

              <Route path="/movies">
                <Movies />
              </Route>

              <Route path="/genre/:id" component={OneGenre} />
              <Route exact path="/genres">
                <Genres />
              </Route>

              <Route path="/admin/movie/:id" component={((props) =>
                <EditMovie {...props} token={token} />)}
              />

              <Route path="/admin" component={((props) =>
                <Admin {...props} token={token} />)}
              />

              <Route path="/">
                <Home />
              </Route>
            </Switch>
          </div>
        </div>
      </div>
    </Router>
  );
}