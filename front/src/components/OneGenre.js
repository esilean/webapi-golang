import React, { useState, useEffect, Fragment } from 'react'
import { Link } from 'react-router-dom'

const OneGenre = (props) => {

    const [movies, setMovies] = useState({})
    const [isLoaded, setLoaded] = useState(false)
    const [error, setError] = useState('')
    const [genreName, setGenreName] = useState('')

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/v1/movies/genre/${props.match.params.id}`)
            .then((response) => {

                if (response.status !== 200) {
                    setError("Invalid response code: " + response.status)
                    return;
                }

                return response.json()

            })
            .then((json) => {
                setGenreName(props.location.genrename)
                setMovies(json.movies)
            })
            .catch(() => {
            })
            .finally(() => {
                setLoaded(true)
            })

    }, [props.location.genrename, props.match.params.id])

    if (error) return (<p>Error: {error}</p>)
    if (!isLoaded) return (<p>Loading...</p>)

    return (
        <Fragment>
            <h2>Genre: {genreName} </h2>

            <div className='list-group'>
                {movies.map((m) => (
                    <Link key={m.id} to={`/movies/${m.id}`} className="list-group-item list-group-item-action">{m.title}</Link>
                ))}
            </div>
        </Fragment>
    );
}

export default OneGenre