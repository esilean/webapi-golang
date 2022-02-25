import React, { useState, useEffect, Fragment } from 'react'

const OneMovie = (props) => {

    const [movie, setMovie] = useState({})
    const [isLoaded, setLoaded] = useState(false)
    const [error, setError] = useState('')

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/v1/movie/${props.match.params.id}`)
            .then((response) => {

                if (response.status !== 200) {
                    setError("Invalid response code: " + response.status)
                    return;
                }

                return response.json()

            })
            .then((json) => {
                setMovie(json.movie)
            })
            .catch(() => {
            })
            .finally(() => {
                setLoaded(true)
            })

    }, [props.match.params.id])

    if (error) return (<p>Error: {error}</p>)
    if (!isLoaded) return (<p>Loading...</p>)

    return (
        <Fragment>
            <h2>{movie.title} ({movie.year})</h2>

            <div className='float-start'>
                <small>Rating: {movie.mpaa_rating}</small>
            </div>

            <div className='float-end'>
                {Object.values(movie.genres).map((m, index) => (
                    <span className='badge bg-secondary me-1' key={index}>
                        {m}
                    </span>
                ))}
            </div>
            <div className='clearfix'></div>

            <hr />

            <table className="table table-compact table-striped">
                <thead></thead>
                <tbody>
                    <tr>
                        <td><strong>Title:</strong></td>
                        <td>{movie.title}</td>
                    </tr>
                    <tr>
                        <td><strong>Description:</strong></td>
                        <td>{movie.description}</td>
                    </tr>
                    <tr>
                        <td><strong>Run time:</strong></td>
                        <td>{movie.runtime} minutes</td>
                    </tr>
                </tbody>
            </table>
        </Fragment>
    );
}

export default OneMovie