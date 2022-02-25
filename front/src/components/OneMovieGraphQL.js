import React, { useState, useEffect, Fragment } from 'react'

const OneMovieGraphQL = (props) => {

    const [movie, setMovie] = useState({})
    const [isLoaded, setLoaded] = useState(false)
    const [error, setError] = useState('')

    useEffect(() => {

        const payload = `
        {
            movie(id: ${props.match.params.id}) {
                id
                title
                runtime
                year
                description
                rating
                poster
            }
        }
        `

        const headers = new Headers()
        headers.append("Content-Type", "application/json")

        const requestOptions = {
            method: 'POST',
            body: payload,
            headers: headers
        }

        fetch("http://localhost:4000/v1/graphql", requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setMovie(data.data.movie)

            })
            .catch((err) => {
                setError(err.message)
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

            {movie.poster !== "" && (
                <div>
                    <img src={`https://image.tmdb.org/t/p/w200${movie.poster}`} alt={'poster'} />
                </div>
            )}

            <div className='float-start'>
                <small>Rating: {movie.rating}</small>
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

export default OneMovieGraphQL