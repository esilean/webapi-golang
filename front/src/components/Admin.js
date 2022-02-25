import React, { Fragment, useEffect, useState } from "react";
import { Link } from 'react-router-dom'

const Admin = (props) => {
    const [movies, setMovies] = useState([])
    const [isLoaded, setLoaded] = useState(false)
    const [error, setError] = useState('')

    useEffect(() => {

        if (props.token === '') {
            props.history.push({
                pathname: '/login'
            })
            return
        }

        fetch(`${process.env.REACT_APP_API_URL}/v1/movies`)
            .then((response) => {

                if (response.status !== 200) {
                    setError("Invalid response code: " + response.status)
                    return;
                }

                return response.json()

            })
            .then((json) => {
                setMovies(json.movies)
            })
            .catch(() => {
            })
            .finally(() => {
                setLoaded(true)
            })

    }, [props.history, props.token])

    if (error) return (<p>Error: {error}</p>)
    if (!isLoaded) return (<p>Loading...</p>)

    return (
        <Fragment>
            <h2>Choose a movie</h2>

            <div className='list-group'>
                {movies.map((m) => (
                    <Link key={m.id} to={`admin/movie/${m.id}`} className="list-group-item list-group-item-action">{m.title}</Link>
                ))}
            </div>
        </Fragment>
    );
}

export default Admin;
