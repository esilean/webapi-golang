import React, { Fragment, useEffect, useState } from "react";
import { Link } from 'react-router-dom'

const Genres = () => {
    const [genres, setGenres] = useState([])
    const [isLoaded, setLoaded] = useState(false)
    const [error, setError] = useState('')

    useEffect(() => {
        fetch('http://localhost:4000/v1/genres')
            .then((response) => {

                if (response.status !== 200) {
                    setError("Invalid response code: " + response.status)
                    return;
                }

                return response.json()

            })
            .then((json) => {
                setGenres(json.genres)
            })
            .catch(() => {
            })
            .finally(() => {
                setLoaded(true)
            })

    }, [])

    if (error) return (<p>Error: {error}</p>)
    if (!isLoaded) return (<p>Loading...</p>)

    return (
        <Fragment>
            <h2>Choose a Genre</h2>

            <div className='list-group'>
                {genres.map((g) => (
                    <Link key={g.id} to={{
                        pathname: `/genre/${g.id}`,
                        genrename: g.name
                    }} className="list-group-item list-group-item-action">{g.name}</Link>
                ))}
            </div>
        </Fragment>
    );
}

export default Genres