import React, { useState, useEffect, Fragment } from 'react'

const OneMovie = (props) => {

    const [movie, setMovie] = useState({})

    useEffect(() => {

        setMovie({
            id: props.match.params.id,
            title: "Some movie",
            runtime: 150,
        })
    }, [props.match.params.id])

    return (
        <Fragment>
            <h2>Movie: {movie.title} {movie.id}</h2>

            <table className="table table-compact table-striped">
                <thead></thead>
                <tbody>
                    <tr>
                        <td><strong>Title:</strong></td>
                        <td>{movie.title}</td>
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